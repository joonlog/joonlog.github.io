---
title : 자체 관리형 Kubernetes에서의 Teleport 구축
date : 2025-10-10 09:00:00 +09:00
categories : [Kubernetes, Teleport]
tags : [kubernetes, k8s, self managed k8s, teleport]  #소문자만 가능
---

> 매니지드 K8S가 아닌 자체 관리형 K8S에서의 Teleport구축 방법
> 

> UI 접근을 위한 외부 통신으로 MetalLB+HAproxy+Nginx Controller가 구성 되었다는 전제 하에 구축
> 
> 
> > CSP 환경이라면 Ingress 설정에서 AWS ALB Controller와 같은 로드밸런서를 사용해 훨씬 간편하게 구축 가능
> > 
- Teleport 공식문서:
https://goteleport.com/docs/
https://goteleport.com/docs/zero-trust-access/deploy-a-cluster/helm-deployments/kubernetes-cluster/

### Teleport 설치

1. Helm 리포지토리 추가
    
    ```bash
    helm repo add teleport https://charts.releases.teleport.dev
    helm repo update
    ```
    
2. teleport-cluster 네임스페이스 생성
    - 네임스페이스에 `baseline` 보안 수준 적용
    
    ```bash
    kubectl create ns teleport-cluster
    kubectl label namespace teleport-cluster 'pod-security.kubernetes.io/enforce=baseline'
    ```
    
3. pvc 생성
    - teleportpvc.yaml
        - kubectl apply -f teleportpvc.yaml
    
    ```bash
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: teleport-pvc
      namespace: teleport-cluster
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: longhorn
      resources:
        requests:
          storage: 10Gi
    ```
    
4. SSL 인증서 Secret 등록
    - 인증서 발급
        - `/etc/letsencrypt/live/teleport.<도메인>/fullchain.pem`
        - `/etc/letsencrypt/live/teleport.<도메인>/privkey.pem`
    
    ```bash
    certbot certonly --manual --preferred-challenges dns \
      -d teleport.<도메인> \
      --key-type ecdsa
    ```
    
    - secret 생성
    
    ```bash
    kubectl -n teleport-cluster create secret tls teleport-tls \
      --cert=/etc/letsencrypt/live/teleport.<도메인>/fullchain.pem \
      --key=/etc/letsencrypt/live/teleport.<도메인>/privkey.pem
    ```
    
5. teleport 설치
    - https://goteleport.com/docs/reference/helm-reference/
    - values.yaml 수정
        - helm으로 설치
        
        ```bash
        helm show values teleport/teleport-cluster > values.yaml
        vim values.yaml
        
        helm install teleport teleport/teleport-cluster \
          -f values.yaml \
          -n teleport-cluster
          
        helm upgrade --install teleport teleport/teleport-cluster \
          -f values.yaml \
          -n teleport-cluster
        ```
        
    - 클러스터 이름 기입
        
        ```bash
        clusterName: teleport.<도메인>
        ```
        
    - 외부 접근 주소
        
        ```bash
        publicAddr: ['teleport.<도메인>:8443']
        tunnelPublicAddr: ['teleport.<도메인>:8443']
        ```
        
    - 프록시가 443 포트 하나에서 ALPN으로 Web UI/SSH 리버스 터널/Kubernetes 프록시를 다 처리
        - 기본값이 seprate인데 이렇게 되면 3024/3025/3026 등등 각각 포트 사용해야함
        - 외부 주소는 publicAddr로 노출됨
        
        ```bash
        proxyListenerMode: "multiplex"
        ```
        
    - Proxy Service
        
        ```bash
        service:
          type: LoadBalancer
          spec:
            # MetalLB에서 고정 IP를 쓰고 싶을 때만 지정
            # loadBalancerIP: 172.27.1.103
        ```
        
    - ACME 비활성화
        
        ```bash
          acme: false
        ```
        
    - TLS 인증서 Secret 지정
        - cert-manager가 발급한 Secret 이름으로 설정
        
        ```bash
          tls:
            existingSecretName: teleport-tls
        ```
        
    - Auth 스토리지 PV 설정
        - pv는 longhorn이 자동으로 프로비저닝
        
        ```bash
        persistence:
          enabled: true
          existingClaimName: "teleport-pvc"
          volumeSize: 10Gi
        ```
        
    - HA replica 설정
        - HA 안쓰면 1
        
        ```bash
        highAvailability:
          replicaCount: 1
        ```
        

### Teleport UI 접근

- Teleport는 443포트에서 HTTPS/SSH 리버스 터널/K8S Access Proxy 등을 `ALPN`으로 동시에 처리하기 때문에, 다른 애플리케이션처럼 ingress로 프록시하면 안됨
    - ingress가 처리하면 ALPN 정보가 사라져서 Teleport Proxy는 트래픽의 소스를 알 수 없게 됨
        - web ui는 접근 되더라도 tsh login, tctl, kubectl 등은 모두 실패할 것
    - Teleport의 Proxy 서비스를 Load Balancer 타입 서비스로 생성해야 함
- 마찬가지 이유로 haproxy 설정에서도 L7 프록시 말고 L4 프록시가 필요
1. HAproxy 설정
- values로 생성한 load balancer를 조회하여 MetalLB로부터 부여받은 IP 확인
    - 현재는 172.27.1.103
    
    ```bash
    # kubectl get svc -n teleport-cluster 
    NAME                TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)                                                                     AGE
    teleport            LoadBalancer   10.233.62.211   172.27.1.103   443:30523/TCP,3023:31718/TCP,3026:31539/TCP,3024:31393/TCP,3036:30139/TCP   56s
    ```
    
- 확인한 IP를 HAproxy 설정에 추가
    - HAproxy 서버 공인 IP로 접근 시 teleport의 LB Service로 통신되도록 설정
    
    ```bash
    vim /etc/haproxy/haproxy.cfg
    haproxy -c -f /etc/haproxy/haproxy.cfg
    systemctl reload haproxy
    ```
    
    ```bash
    frontend unified_frontend_8443
        bind *:8443
        mode tcp
        option tcplog
        tcp-request inspect-delay 5s
        tcp-request content accept if { req_ssl_hello_type 1 }
        use_backend metallb_backend_teleport if { req.ssl_sni -i teleport.<도메인> }
        default_backend metallb_backend_teleport
    
    backend metallb_backend_teleport
        mode tcp
        server teleport 172.27.1.103:443
    ```
    
1. 접근 성공!

![TeleportBuild01.png](/assets/img/kubernetes/TeleportBuild01.png)

---
title: "외부에서 자체 관리형 Kubernetes 접근을 위한 MetalLB/HAproxy 설정과 통신 구조"
date: 2025-08-24T09:00:00+09:00
categories: ["Kubernetes", "Plugins"]
tags: ["kubernetes", "k8s", "local k8s", "metallb", "haproxy", "nginx controller", "self managed k8s"]
---


- 외부에서의 로컬 K8S 클러스터 접근을 위해 MetalLB와 HAproxy를 사용
- K8S 클러스터 내부의 로드밸런서로는 nginx controller를 사용 중인 상태
    - AWS라면, ALB Controller 하나만으로 MetalLB+HAproxy+nginxcontroller 기능을 커버 가능
    - 사실상의 CSP의 K8S를 쓰는 이유??

# MetalLB

> https://metallb.universe.tf/installation/
> 
- 클러스터 내부용 로드밸런서
    - 오픈소스
    - CSP의 LB를 사용하지 못하는 경우(내/외부 양쪽에서 사용가능한 LB)에, 내부용 L4 로드밸런서로 사용하기 위해 주로 MetalLB를 사용

### 사전 준비

- kube-proxy의 configmap에서 strictARP 활성화

```bash
kubectl edit configmap -n kube-system kube-proxy

apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
ipvs:
  strictARP: true
```

### MetalLB 설치

- helm으로 설치

```bash
helm repo add metallb https://metallb.github.io/metallb
helm install metallb metallb/metallb
```

### MetalLB 설정

- IP Pool 설정

```bash
cat <<EOF | kubectl apply -f -
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: default-pool
  namespace: default
spec:
  addresses:
  - 172.27.1.100-172.27.1.110
EOF
```

- L2Advertisement 설정

```bash
cat <<EOF | kubectl apply -f -
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: default
  namespace: default
spec:
  ipAddressPools:
  - default-pool
EOF
```

# HAproxy

- haproxy 설치

```bash
yum install -y haproxy
systemctl enable --now haproxy
```

- haproxy 설정
  - Jenkins를 예시로 구성
    
    ```bash
    # kubectl get ingress -n jenkins
    NAME              CLASS   HOSTS           ADDRESS        PORTS   AGE
    jenkins-ingress   nginx   jenkins.local   172.27.1.100   80      26h
    ```
    
    - 파드가 MetalLB로부터 부여받은 IP가 172.27.1.100일 경우 하기와 같이 설정
        - 백엔드 헬스체크는 MetalLB에서 하고 있으니 HAproxy에선 설정하지 않음
    - `http-request set-header Host`
        - 클라이언트에서 오는 모든 HTTP 요청의 Host 헤더를 jenkins.local로 변경
    - `http-request del-header X-Forwarded-Host`
        - 이전 프록시에서 설정한 Host 헤더 제거
    - `http-request del-header X-Forwarded-Proto`
        - 이전 프록시에서 설정한 프로토콜 정보 제거
    
    ```bash
    tee /etc/haproxy/haproxy.cfg > /dev/null <<EOF
    frontend jenkins_frontend
        bind *:8080
        default_backend metallb_backend
    
    backend metallb_backend
        server metallb 172.27.1.100:80
    EOF
    ```

# 통신 구조

- jenkins 파드가 하기와 같이 구성되어 있다 가정할 때 통신 구조
    - 현재 외부에서의 접근 시 80/443 포트를 사용하지 못하고 별도의 포트를 사용해야 하는 환경이기 때문에 8080포트를 사용

```bash
# kubectl get ingress -n jenkins
NAME              CLASS   HOSTS           ADDRESS        PORTS   AGE
jenkins-ingress   nginx   jenkins.local   172.27.1.100   80      26h
# kubectl get all -n jenkins 
NAME            READY   STATUS    RESTARTS   AGE
pod/jenkins-0   2/2     Running   0          9d

NAME                    TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
service/jenkins         ClusterIP   10.233.5.87    <none>        8080/TCP    9d
service/jenkins-agent   ClusterIP   10.233.28.70   <none>        50000/TCP   9d

NAME                       READY   AGE
statefulset.apps/jenkins   1/1     9d
```

1. 브라우저에서 <도메인>:8080으로 접속
2. DNS 해석되어 공인IP:8080으로 라우팅
    - 공인IP:8080는 HAproxy 서버
3. HAproxy frontend 설정으로 8080포트로 들어오는 요청을 HAproxy backend로 라우팅
    - HAproxy backend는 172.27.1.100:80
4. 172.27.1.100 IP는 MetalLB의 IP 풀에서부터 할당받은 Nginx Controller의 IP
    - Nginx Ingress Controller는 LoadBalancer 타입 Service
5. Nginx Ingress Controller는 HAproxy backend로부터 받은 트래픽의 Host 헤더가 Jenkins Ingress에 정의된 Host인 jenkins.local인지 확인
6. Jenkins Ingress에 정의된 규칙에 따라 Jenkins Service로 라우팅
7. Jenkins Service에서 Jenkins Pod로 라우팅
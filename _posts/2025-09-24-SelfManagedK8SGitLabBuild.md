---
title : 자체 관리형 Kubernetes에서의 GitLab 구축
date : 2025-09-24 09:00:00 +09:00
categories : [Kubernetes, GitLab]
tags : [kubernetes, k8s, self managed k8s, scm, gitlab, gitlab registry, minio]  #소문자만 가능
---

> 매니지드 K8S가 아닌 자체 관리형 K8S에서의 GitLab 구축 방법
> 

> UI 접근을 위한 외부 통신으로 MetalLB+HAproxy+Nginx Controller가 구성 되었다는 전제 하에 구축
> 
> 
> > CSP 환경이라면 Ingress 설정에서 AWS ALB Controller와 같은 로드밸런서를 사용해 훨씬 간편하게 구축 가능
> > 
- GitLab 공식문서:
****https://docs.gitlab.com/install/install_methods/#helm-chart
- GitLab Artifacthub 문서:
https://artifacthub.io/packages/helm/gitlab/gitlab

### GitLab 설치

1. Helm 리포지토리 추가
    
    ```bash
    helm repo add gitlab https://charts.gitlab.io/
    helm repo update
    ```
    
2. gitlab 네임스페이스 생성
    
    ```bash
    kubectl create ns gitlab
    ```
    
3. gitlab-pv 생성(StorageClass 없을 시에만 진행)
    - 기존에 사용하던 Ceph, Longhorn, EBS Driver, Harbor 같은 StorageClass가 없을 경우에 로컬 경로를 사용한 PV 할당 방법
        - StorageClass가 있다면 PVC 생성 시 자동으로 PV가 할당되기 때문에 이 과정은 불필요
    
    ```bash
    vim gitlab-01-volume.yaml
    kubectl apply -f gitlab-01-volume.yaml
    ```
    
    ```bash
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: gitlab-pv
    spec:
      storageClassName: gitlab-pv
      accessModes:
      - ReadWriteOnce
      capacity:
        storage: 20Gi
      persistentVolumeReclaimPolicy: Retain
      hostPath:
        path: /data/gitlab-volume/
    
    ---
    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
      name: gitlab-pv
    provisioner: kubernetes.io/no-provisioner
    volumeBindingMode: WaitForFirstConsumer
    ```
    
4. gitlab 설치
    - values.yaml 수정
        - helm으로 설치
    
    ```bash
    helm show values gitlab/gitlab > values.yaml
    vim values.yaml
    
    helm install gitlab -n gitlab -f values.yaml gitlab/gitlab
    ```
    
    - 소스코드 리포지토리 설정
    
    ```bash
    # GitLab Gitaly (Git 저장소)
    gitlab:
      gitaly:
        persistence:
          storageClass: "longhorn"
          size: 20Gi
    ```
    
    - postgresql, redis 설정
        - postgresql: gitlab 전체의 메타데이터 저장
            - redis: 사용자 세션, 백그라운드 작업, 캐시, 알림 등에 사용
    
    ```bash
    # PostgreSQL
    postgresql:
      primary:
        persistence:
          storageClass: "longhorn"
          size: 8Gi
    
    # Redis  
    redis:
      master:
        persistence:
          storageClass: "longhorn"
          size: 2Gi
    ```
    
    - minio설정
        - https://docs.gitlab.com/charts/charts/globals#configure-minio-settings
        - S3와 호환되는 오브젝트 스토리지
    
    ```bash
    global:
    
      # minio 서비스 활성화
      minio:
        enabled: true
        
      # 바이너리 파일들이 오브젝트 스토리지를 사용할 수 있도록 설정
      appConfig:
        object_store
          enabled: true
          
      # Container Registry가 minio의 registry 버킷을 사용
      registry:
        bucket: registry
        
    # minio가 Longhorn을 통한 PV 사용
    minio:
      persistence:
        enabled: true
        size: 10Gi
        storageClass: longhorn
    ```
    
    - 불필요한 기능 비활성화
    
    ```bash
    global:
      lfs:
        enabled: false
      uploads:
        enabled: false
      toolbox:
        enabled: false
        
    # GitLab KAS 비활성화 (Kubernetes Agent)
    global:
      kas:
        enabled: false
    
    # CI/CD 관련 비활성화
    global:
      appConfig:
        artifacts:
          enabled: false
        packages:
          enabled: false
    
    # 모니터링 비활성화
    prometheus:
      install: false
    
    # GitLab Runner 비활성화
    gitlab-runner:
      install: false
    ```
    
    - 도메인, 에디션(CE/EE)
    
    ```bash
    global:
      edition: ce
      hosts:
        domain: <도메인>
    ```
    
    - http 접근을 위한 설정
    
    ```bash
    global:
      hosts:
        https: false
      ingress:
        configureCertmanager: false
        tls:
          enabled: false
    installCertmanager: false
    ```
    

### GitLab UI

> GitLab Helm Chart를 통해 ingress는 생성 완료!
> 

> values.yaml 내에 ingress 생성을 위한 nginx 설정이 포함되어 있음
> 
1. HAproxy 설정
- ingress를 조회하여 MetalLB로부터 부여받은 IP 확인
    - 현재는 172.27.1.102
    
    ```bash
    # kubectl get ingress -n gitlab
    NAME                        CLASS          HOSTS                    ADDRESS        PORTS   AGE
    gitlab-registry             gitlab-nginx   registry.<도메인>   172.27.1.102   80      61m
    gitlab-webservice-default   gitlab-nginx   gitlab.<도메인>     172.27.1.102   80      61m
    ```
    
- 확인한 IP를 HAproxy 설정에 추가
    - HAproxy 서버 공인 IP로 접근 시 gitlab의 Ingress로 통신되도록 설정
    - `http-response replace-value Location ^http://gitlab\.<도메인>/(.*)$ http://gitlab.<도메인>:<포트>/\\1`
    - 두 조건을 만족할 경우 설정 필요
        - 외부 브라우저에서 GitLab 접근 시 포트를 붙여서 접근해야 될 경우
        - 애플리케이션 내부적으로 url로 응답이 나오도록 설정되어 있을 경우
            - EX) GitLab 로그인 시 GitLab에 설정된 url 응답
    - 설정하지 않으면 접근 시 자동으로 포트가 제거된 상태로 리다이렉션되어 404 에러
    
    ```bash
    vim /etc/haproxy/haproxy.cfg
    haproxy -c -f /etc/haproxy/haproxy.cfg
    systemctl reload haproxy
    ```
    
    ```bash
    frontend unified_frontend_8080
        bind *:8080
        mode http
        option forwardfor
    
        http-request set-header X-Forwarded-Host %[req.hdr(host)]
        http-request set-header X-Forwarded-Proto http
        http-request set-header X-Forwarded-Port %[dst_port]
    
        # Host 기반 라우팅
        use_backend metallb_backend_gitlab if { hdr(host) -m sub gitlab }
        # GitLab Container Registry를 위한 경로
        use_backend metallb_backend_gitlab if { hdr(host) -m sub registry }
    
    backend metallb_backend_gitlab
        http-response replace-value Location ^http://gitlab\.<도메인>/(.*)$ http://gitlab.<도메인>:<포트>/\\1
        server gitlab 172.27.1.100:80
    ```
    
1. 접근 성공!
    - gitlab 관리자(root) 비밀번호 확인
        
        ```bash
        kubectl get secret gitlab-gitlab-initial-root-password -n gitlab -o jsonpath="{.data.password}" | base64  --decode
        ```
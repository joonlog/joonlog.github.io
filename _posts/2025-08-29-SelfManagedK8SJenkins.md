---
title : 자체 관리형 Kubernetes에서의 Jenkins 구축
date : 2025-08-29 09:00:00 +09:00
categories : [Kubernetes, application]
tags : [kubernetes, k8s, self managed k8s, jenkins]  #소문자만 가능
---

> 매니지드 K8S가 아닌 자체 관리형 K8S에서의 Jenkins 구축 방법
> 

> 외부에서의 통신을 위한 MetalLB+HAproxy+Nginx Controller가 구성 되었다는 전제 하에 구축
> 

> CSP 환경이라면 Ingress 설정에서 AWS ALB Controller와 같은 로드밸런서를 사용해 훨씬 간편하게 구축 가능
> 

- **Jenkins.io 공식문서:**
https://www.jenkins.io/doc/book/installing/kubernetes/

### 1. Helm을 통한 설치

```bash
helm repo add jenkinsci https://charts.jenkins.io
helm repo update
```

### 2. jenkins-ci 네임스페이스 생성

```bash
kubectl create ns jenkins
```

### 3. jenkins-pv 생성

```bash
vim jenkins-01-volume.yaml
kubectl apply -f jenkins-01-volume.yaml
```

- https://raw.githubusercontent.com/jenkins-infra/jenkins.io/master/content/doc/tutorials/kubernetes/installing-jenkins-on-kubernetes/jenkins-01-volume.yaml

```bash
apiVersion: v1
kind: PersistentVolume
metadata:
  name: jenkins-pv
spec:
  storageClassName: jenkins-pv
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 20Gi
  persistentVolumeReclaimPolicy: Retain
  hostPath:
    path: /data/jenkins-volume/

---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: jenkins-pv
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
```

### 4. ServiceAccount 생성

```bash
vim jenkins-02-sa.yaml
kubectl apply -f jenkins-02-sa.yaml
```

- https://raw.githubusercontent.com/jenkins-infra/jenkins.io/master/content/doc/tutorials/kubernetes/installing-jenkins-on-kubernetes/jenkins-02-sa.yaml

```bash
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: jenkins
  namespace: jenkins
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
  name: jenkins
rules:
- apiGroups:
  - '*'
  resources:
  - statefulsets
  - services
  - replicationcontrollers
  - replicasets
  - podtemplates
  - podsecuritypolicies
  - pods
  - pods/log
  - pods/exec
  - podpreset
  - poddisruptionbudget
  - persistentvolumes
  - persistentvolumeclaims
  - jobs
  - endpoints
  - deployments
  - deployments/scale
  - daemonsets
  - cronjobs
  - configmaps
  - namespaces
  - events
  - secrets
  verbs:
  - create
  - get
  - watch
  - delete
  - list
  - patch
  - update
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
  - list
  - watch
  - update
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
  name: jenkins
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: jenkins
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:serviceaccounts:jenkins
```

### 5. Jenkins 설치

- https://github.com/jenkinsci/helm-charts/tree/main/charts/jenkins
    - helm으로 설치
    
    ```bash
    wget https://raw.githubusercontent.com/jenkinsci/helm-charts/main/charts/jenkins/values.yaml
    vim values.yaml
    
    helm install jenkins -n jenkins -f helm/values.yaml jenkinsci/jenkins
    ```
    
- https://raw.githubusercontent.com/jenkinsci/helm-charts/main/charts/jenkins/values.yaml
    
    ```bash
    persistence:
      storageClass: jenkins-pv
      size: "20Gi"
      
    serviceAccount:
      create: false
      name: jenkins
    ```
    
    - /data/jenkins-volume 접근 권한으로 인해 오류 생길 경우 권한 부여 필요
        
        ```bash
        sudo chown -R 1000:1000 /data/jenkins-volume/
        sudo chmod -R 755 /data/jenkins-volume/
        ```
        

### 6. ingress 생성

- 8080포트를 사용하는 jenkins service로 연결

```bash
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: jenkins-ingress
  namespace: jenkins
  annotations:
    nginx.ingress.kubernetes.io/backend-protocol: "HTTP"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
spec:
  ingressClassName: nginx
  rules:
  - host: "jenkins.local"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: jenkins
            port:
              number: 8080
EOF
```

### 7. HAproxy 설정 추가

- ingress를 조회하여 MetalLB로부터 부여받은 IP 확인
    - 현재는 172.27.1.100
    
    ```bash
    # kubectl get ingress -n jenkins
    NAME              CLASS   HOSTS           ADDRESS        PORTS   AGE
    jenkins-ingress   nginx   jenkins.local   172.27.1.100   80      26h
    ```
    
- 확인한 IP를 HAproxy 설정에 추가
    - HAproxy 서버 공인 IP로 접근 시 Jenkins의 Ingress로 통신되도록 설정
    
    ```bash
    tee /etc/haproxy/haproxy.cfg > /dev/null <<EOF
    frontend jenkins_frontend
        bind *:8080
        default_backend metallb_backend
    
    backend metallb_backend
        server metallb 172.27.1.100:80 check
    EOF
    ```
    
- 접근 성공!
    
    ![Jenkins1.png](/assets/img/kubernetes/Jenkins1.png)

    - Jenkins 관리자 비밀번호 확인
        
        ```bash
        kubectl exec -it jenkins-0 -n jenkins -- cat /var/jenkins_home/secrets/initialAdminPassword
        ```
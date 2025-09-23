---
title : Kubernetes 인증 / 권한 관리 Overview
date : 2025-09-23 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, role, rolebinding, clusterrole, clusterrolebinding]  #소문자만 가능
---

> kubernetes에서 kubectl 같이 클라이언트나 파드 내부에서 API 서버에 접근할 때는 인증과 권한이 필요
> 

### 접근 제어

User/Group → 인증

- 인증서는 클라이언트 신원을 확인하는 데 사용
    - ex) `kubectl get nodes` 실행 시 사용자 인증 필요
- 인증서(certificate)를 kubeconfig에 등록하여 사용

ServiceAccount → 권한

- 파드 내부에서 쿠버네티스 API를 호출할 때 사용
    - Pod 실행 시 자동으로 토큰이 `/var/run/secrets/kubernetes.io/serviceaccount/token` 에 마운트
- 파드 내부에서 자동으로 마운트되거나 Secret을 통해 주입

### 사용자 생성 절차

1. 개인키 생성
2. CSR 생성
3. CSR 요청 및 승인
    
    ```bash
    kubectl apply -f csr.yaml
    kubectl certificate approve myuser-csr
    ```
    
4. Role/RoleBinding 생성
    
    ```bash
    kind: Role
    apiVersion: rbac.authorization.k8s.io/v1
    metadata:
      namespace: default
      name: pod-reader
    rules:
    - apiGroups: [""]
      resources: ["pods"]
      verbs: ["get", "list"]
    ---
    kind: RoleBinding
    apiVersion: rbac.authorization.k8s.io/v1
    metadata:
      name: read-pods
      namespace: default
    subjects:
    - kind: User
      name: myuser
    roleRef:
      kind: Role
      name: pod-reader
      apiGroup: rbac.authorization.k8s.io
    ```
    
5. kubeconfig에 사용자 추가
    
    ```bash
    kubectl config set-credentials myuser --client-certificate=user.crt --client-key=user.key
    ```
    
6. context 변경
    
    ```bash
    kubectl config set-context myuser-context --cluster=mycluster --namespace=default --user=myuser
    kubectl config use-context myuser-context
    ```
    

### 서비스 어카운트 생성 절차

1. 서비스 어카운트 생성
    
    ```bash
    kubectl create serviceaccount mysa -n default
    ```
    
2. Rolebinding 생성
    
    ```bash
    kind: RoleBinding
    apiVersion: rbac.authorization.k8s.io/v1
    metadata:
      name: sa-read-pods
      namespace: default
    subjects:
    - kind: ServiceAccount
      name: mysa
      namespace: default
    roleRef:
      kind: Role
      name: pod-reader
      apiGroup: rbac.authorization.k8s.io
    ```
    
3. pod 실행시 서비스 어카운트 지정
    
    ```bash
    apiVersion: v1
    kind: Pod
    metadata:
      name: test-pod
    spec:
      serviceAccountName: mysa
      containers:
      - name: curl
        image: curlimages/curl
        command: ["sleep", "3600"]
    ```
    

### 권한 관리(RBAC)

- Role
    - 네임스페이스 단위 권한
        - ex) default namespace에서 pod 읽기
- RoleBinding
    - Role을 특정 User/Group/ServiceAccount에 바인딩
- ClusterRole
    - 클러스터 전체 권한
        - ex) 모든 namespace에서 pod 읽기
- ClusterRoleBinding
    - ClusterRole을 User/Group/ServiceAccount에 바인딩
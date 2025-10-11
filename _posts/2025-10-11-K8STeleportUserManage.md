---
title : Kubernetes환경 Teleport에서 사용자 관리하기
date : 2025-10-11 09:00:00 +09:00
categories : [Kubernetes, Teleport]
tags : [kubernetes, k8s, self managed k8s, teleport, tsh]  #소문자만 가능
---

### Teleport 사용자 생성

- member Role 생성
    - member.yaml
    
    ```bash
    kind: role
    version: v7
    metadata:
      name: member
    spec:
      allow:
        kubernetes_groups: ["system:masters"]
        kubernetes_labels:
          '*': '*'
        kubernetes_resources:
          - kind: '*'
            namespace: '*'
            name: '*'
            verbs: ['*']
    ```
    
    ```bash
    kubectl exec -i deployment/teleport-auth -n teleport-cluster -- tctl create -f < member.yaml
    ```
    
- myuser 사용자 생성
    - 출력된 링크를 브라우저에서 접속
    
    ```bash
    kubectl exec -it deployment/teleport-auth -n teleport-cluster -- tctl users add myuser --roles=member,access,editor
    ```
    
- 로그인 성공!
    
    ![TeleportUserManage01.png](/assets/img/kubernetes/TeleportUserManage01.png)
    

### tsh 로그인

- tsh login
    
    ```bash
    tsh login --proxy=teleport.<도메인>:8443 --user=myuser
    ```
    
- 사용자에 할당된 사용 가능한 k8s cluster 목록
    
    ```bash
    # tsh kube ls
    Kube Cluster Name      Labels Selected 
    ---------------------- ------ -------- 
    teleport.<domain>        *      
    ```
    
- k8s cluster login
    
    ```bash
    tsh kube login teleport.<domain>
    ```
    
- kubectl 성공!
    
    ```bash
    # kubectl get nodes
    NAME      STATUS   ROLES           AGE   VERSION
    master1   Ready    control-plane   51d   v1.31.4
    master2   Ready    control-plane   51d   v1.31.4
    master3   Ready    control-plane   51d   v1.31.4
    node1     Ready    <none>          51d   v1.31.4
    node2     Ready    <none>          51d   v1.31.4
    ```
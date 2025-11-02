---
title: "Kubernetes Pod Scheduling Overview"
date: 2025-09-19T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["kubernetes", "k8s", "node affinity", "pod affinity", "taint", "toleration", "cordon", "drain"]
---


> k8s scheduler가 어떤 기준으로 파드를 노드에 배치하는지에 대한 정리
> 

### Node Affinity

- 지정된 조건에 맞는 노드에 파드를 배포
    - `disktype=ssd` label이 있는 노드에만 스케줄링
    
    ```bash
    apiVersion: v1
    kind: Pod
    metadata:
      name: web-pod
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: disktype
                operator: In
                values:
                - ssd
      containers:
      - name: web
        image: nginx
    ```
    

### Pod Affinity

- pod를 가까이 배치(ex: 같은 노드에 배치)
    - `app=db` label이 붙은 파드와 같은 노드에 배치
        - 반대로 `podAntiAffinity`를 쓰면 다른 노드에 배치 가능
    
    ```bash
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - db
        topologyKey: "kubernetes.io/hostname"
    ```
    

### Taint & Toleration

- nodeselector: 노드레이블이 지정된 노드에 파드를 배포할 때 사용
- node taint: 노드에 파드를 배포되지 않도록 할 때 사용
    - node taint 추가 명령어
    
    ```bash
    kubectl taint nodes node1 key=value:NoSchedule
    ```
    
- pod toleration: node taint가 설정된 노드에 배포하고 싶을 때 사용
    - pod 매니페스트 내 toleration 설정
        - 이 파드는 위 node1에 배치 가능
    
    ```bash
    tolerations:
    - key: "key"
      operator: "Equal"
      value: "value"
      effect: "NoSchedule"
    ```
    

### Cordon/Uncordon & Drain

- cordon: 지정된 노드에 파드 배포 금지
    
    ```bash
    kubectl cordon node1
    ```
    
- uncordon: 지정된 노드에 파드 배포 금지 해제
    
    ```bash
    kubectl uncordon node1
    ```
    
- drain: 지정된 노드에 모든 파드 삭제
    
    ```bash
    kubectl drain node1 --ignore-daemonsets
    ```
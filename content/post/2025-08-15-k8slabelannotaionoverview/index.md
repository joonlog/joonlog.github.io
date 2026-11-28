---
title: "Kubernetes Label과 Annotation"
slug: "k8s-labels-annotations"
aliases: ["/p/kubernetes-label과-annotation/"]
date: 2025-08-17T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["Kubernetes"]
---

## Label

- 리소스를 식별하기 위한 key-value 쌍

- selector로 검색/필터링 가능

- label은 3개 이상 정의할 것을 권장

    ```bash
    environment=prod
    tier=frontend
    release=stable
    team=devops1
    ```

## Annotation

- 메타데이터를 key, value 쌍으로 정의

- k8s에 메타데이터를 전달할 용도로 사용

- 관리를 위해 필요한정보를 기록할 용도로 사용

    ```bash
    builder: BSC
    buildDate: "2024-04-03"
    imageRegistry: "https://registry.example.com/nginx"
    ```

- [참고] Labels vs Annotations

    - lables: 클러스터를 관리할 때 사용자가 원하는 값을 입력

    - annotations: 클러스터 시스템에 필요한 정보를 표시하는데 사용

- Label+Annotation 예시

```bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    environment: prod
    tier: frontend
    release: stable
    team: devops1
  annotations:
    kubernetes.io/change-cause: "nginx image upgraded to 1.14"
    builder: "BSC"
    buildDate: "2024-04-03"
    imageRegistry: "https://registry.example.com/nginx"
spec:
  replicas: 3
  selector:
    matchLabels:
      tier: frontend
  template:
    metadata:
      labels:
        environment: prod
        tier: frontend
        release: stable
        team: devops1
      annotations:
        lastChecked: "2024-04-05T10:00:00Z"
    spec:
      containers:
      - name: nginx
        image: nginx:1.14
        ports:
        - containerPort: 80
```

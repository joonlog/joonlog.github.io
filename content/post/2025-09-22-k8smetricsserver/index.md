---
title: "자체 관리형 Kubernetes에서 Metrics Server를 통한 매트릭 자원 확인"
slug: "self-managed-k8s-metrics-server"
aliases: ["/p/자체-관리형-kubernetes에서-metrics-server를-통한-매트릭-자원-확인/"]
date: 2025-09-22T09:00:00+09:00
categories: ["Kubernetes", "Plugins"]
tags: ["Kubernetes", "Metrics Server", "Self-Managed Kubernetes"]
---


> `kubectl top node node1`과 같이 즉각적인 매트릭 자원 확인이 가능하도록 설정
> 
- GitHub:
https://github.com/kubernetes-sigs/metrics-server?tab=readme-ov-file
- Artifacthub:
https://artifacthub.io/packages/helm/metrics-server/metrics-server
- helm 리포지토리 추가

```bash
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm upgrade
```

- values.yaml 수정

```bash
helm show values metrics-server/metrics-server > metrics-server-values.yaml
vim metrics-server-values.yaml
```

- metrics server가 tls 인증서를 검증하지 않도록 설정

```bash
args:
    - --kubelet-insecure-tls
```

- metrics server 설치

```bash
helm upgrade --install metrics-server metrics-server/metrics-server -f metrics-server-values.yaml
```
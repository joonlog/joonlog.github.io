---
title: "Kubernetes 아키텍쳐"
date: 2025-08-05T09:00:00+09:00
categories: ["Kubernetes", "Plugins"]
tags: ["kubernetes", "k8s", "nginx controller", "ingress controller", "helm"]
---


### Nginx Ingress Controller

- k8s 외부에서 들어오는 http/https 트래픽을 1차적으로 처리하는 Entry Point
    
    > Ingress Controller가 없다면 NodePort Service 만으로 Pod를 열어야한다.
    > 
- Ingress 리소스를 통해 내부 서비스로 라우팅
- TLS Termination, 리버스 프록시, 로드 밸린싱 등 수행 가능

https://docs.nginx.com/nginx-ingress-controller/installation/installing-nic/installation-with-helm/

- helm을 통한 설치

```bash
helm install my-release oci://ghcr.io/nginx/charts/nginx-ingress --version 2.2.1
```

- helm을 통한 설치 시 옵션 수정이 필요한 경우
    - 하기 명령어로 pull한 helm chart의 values.yaml에서 옵션 수정

```bash
helm pull oci://ghcr.io/nginx/charts/nginx-ingress --untar --version 2.2.1
```

- 수정한 옵션 적용하여 설치

```bash
helm install my-release . 
```
---
title: "Kubernetes Ingress Overview"
slug: "k8s-ingress-overview"
aliases: ["/p/kubernetes-ingress-overview/"]
date: 2025-08-13T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["Kubernetes"]
---

> Ingress API는 기능이 동결됐고, Kubernetes 프로젝트는 후속으로 Gateway API를 권장한다(Ingress를 없앨 계획은 없음). 많이 쓰던 ingress-nginx 컨트롤러는 2026년 3월에 유지보수가 끝났다(기존 배포는 계속 동작). 새로 구축한다면 Gateway API 구현체를 먼저 검토한다.

### LoadBalancer

- L4 - NLB + Service의 LoadBalancer

    웹서버가 아닌 백엔드 서비스(API, DB 프록시 등)에 대한 부하 분산에 적합

    TCP/UDP 레벨에서 패킷을 라우팅

- L7 - ALB + Ingress

    웹 애플리케이션 HTTP/HTTPS 트래픽 부하 분산에 적합

    URL 경로, 호스트명 등 애플리케이션 계층의 규칙 기반 라우팅 가능

## Ingress

- 트래픽 부하 분산 (HTTP/HTTPS)

- 경로 기반 라우팅

    예: /, /home, /pay, /login

- 도메인 기반 라우팅

    예: www.example.com, www.test.com

- SSL/TLS 인증서 처리

    HTTPS 종료(인증서는 TLS Secret으로 지정). 인증서 자동 발급·갱신은 cert-manager 같은 별도 컨트롤러가 맡음(Let's Encrypt 연계)

- 경로 기반 라우팅 예시

    ```bash
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: path-based-ingress
    spec:
      ingressClassName: nginx
      rules:
        - host: shop.example.com
          http:
            paths:
              - path: /products
                pathType: Prefix
                backend:
                  service:
                    name: product-service
                    port:
                      number: 80
              - path: /orders
                pathType: Prefix
                backend:
                  service:
                    name: order-service
                    port:
                      number: 80
    ```

- 호스트 기반 라우팅 예시

    ```bash
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: host-based-ingress
    spec:
      ingressClassName: nginx
      rules:
        - host: api.example.com
          http:
            paths:
              - path: /
                pathType: Prefix
                backend:
                  service:
                    name: api-service
                    port:
                      number: 80
        - host: admin.example.com
          http:
            paths:
              - path: /
                pathType: Prefix
                backend:
                  service:
                    name: admin-service
                    port:
                      number: 80
    ```

## Ingress Architecture

> Client(HTTP/HTTPS) → Ingress Controller(Ingress Rules 라우팅) → Service(ClusterIP) → Pod

- **Ingress Controller**

    - Ingress 리소스에 정의된 **Ingress Rules**를 읽고 동작

        - Input: `<ingress service>`

        - Output: `<service name>`:80

    - Service는 대상 Pod를 찾는 데 쓰이고, 컨트롤러는 보통 Service의 엔드포인트(Pod IP)로 바로 전달(ClusterIP를 거치지 않음)

- **Service (ClusterIP)**

    - Input: `<service name>`:80

    - Output: <pod’s label>:80

- **Deployment(ReplicaSet, Pod)**

    - 여러 Pod를 레플리카로 생성해 무중단 운영 가능

## Ingress Controller 종류

- **로컬(Kubernetes On-Prem, Minikube 등)** → **NGINX Ingress Controller** 사용(ingress-nginx는 2026-03 유지보수 종료. 새로 쓴다면 Gateway API 구현체)

- **AWS EKS** → **AWS Load Balancer Controller(ALB 기반)** 사용

    - ALB와 연동해 L7 라우팅 처리

    - ACM(AWS Certificate Manager)로 TLS 인증서 관리 가능

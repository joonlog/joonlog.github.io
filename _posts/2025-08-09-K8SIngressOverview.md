---
title : Kubernetes Ingress Overview
date : 2025-08-09 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, ingress, ingress controller]  #소문자만 가능
---

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
    
    HTTPS 통신을 위한 인증서 설치 및 자동 갱신(Let's Encrypt 등과 연계 가능)
    
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
> 
- **Ingress Controller**
    - Ingress 리소스에 정의된 **Ingress Rules**를 읽고 동작
        - Input: <ingress service>
        - Output: <service name>:80
    - 실제 트래픽은 Service(ClusterIP)를 통해 Pod로 전달
- **Service (ClusterIP)**
    - Input: <service name>:80
    - Output: <pod’s label>:80
- **Deployment(ReplicaSet, Pod)**
    - 여러 Pod를 레플리카로 생성해 무중단 운영 가능

## Ingress Controller 종류

- **로컬(Kubernetes On-Prem, Minikube 등)** → **NGINX Ingress Controller** 사용
- **AWS EKS** → **AWS Load Balancer Controller(ALB 기반)** 사용
    - ALB와 연동해 L7 라우팅 처리
    - ACM(AWS Certificate Manager)로 TLS 인증서 관리 가능
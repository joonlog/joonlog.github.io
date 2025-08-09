---
title : Kubernetes Ingress Overview
date : 2025-08-09 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, ingress, ingress controller]  #소문자만 가능
---

### LoadBalancer

- L4 - NLB - Service의 LoadBalancer: 웹서버를 제외한 나머지 부하 분산
- L7 - ALB - Ingress: 웹 로드 밸런싱

## Ingress

- 트래픽 부하 분산
- 경로 기반 라우팅(/, /home, /pay, /login)
- 도메인 기반 라우팅(ex: www.example.com, www.test.com)
- SSL/TLS 인증서 처리(ex: HTTPS)

## Ingress Architecture

- Ingress Controller
    - Ingress Rules
    Input: <ingress service>
    Output: <service name>:80
    - Service(ClusterIP)
    Input: <service name>:80
    Output: <pod’s label>:80
    - Deploument(ReplicaSet, Pod)
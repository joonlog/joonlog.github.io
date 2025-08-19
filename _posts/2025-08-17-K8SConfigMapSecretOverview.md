---
title : Kubernetes ConfigMap과 Secret
date : 2025-08-17 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, configmap, secret]  #소문자만 가능
---

- variable: 특정 파드에 변수값을 설정할 때 사용
- configmap: 네임스페이스 안에 pod에 대한 전역 변수값을 설정할 때 사용
- secret: 네임스페이스 안에 pod에 대한 전역 변수값을 설정할 때 사용. 민감한 설정/데이터(예: 아이디/비밀번호, 인증서, 토큰)

## Configmap

- key-value 쌍으로 정의
- 네임스페이스안에 실행되는 파드에서 공통적으로 사용
- 주로 설정
    - 변수
    - 설정 파일

## Secret

- key-value 쌍으로 정의
- 네임스페이스안에 실행되는 파드에서 공통적으로 사용
- 주로 설정
    - 도커 레지스트리 인증 정보
    - TLS 인증서
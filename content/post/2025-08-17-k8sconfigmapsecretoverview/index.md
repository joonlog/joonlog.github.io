---
title: "Kubernetes ConfigMap과 Secret"
date: 2025-08-17T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["kubernetes", "k8s", "configmap", "secret"]
---


- Kubernetes에서 애플리케이션을 배포할 때 환경변수나 설정 등을 관리하는 방법
    - Pod yaml에 직접 값을 적을수도 있지만 재사용/보안을 위해 ConfigMap/Secret을 사용

## Configmap

- key-value 쌍으로 정의
- DB 호스트, 포트, 로그 레벨 등의 민감하지 않은 설정 값을 저장할 때 사용
- 예시 - 환경변수
    - ConfigMap
        
        ```bash
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: app-config
          namespace: demo
        data:
          APP_MODE: "production"
          LOG_LEVEL: "debug"
        ```
        
    - Pod
        
        ```bash
        apiVersion: v1
        kind: Pod
        metadata:
          name: sample-pod
          namespace: demo
        spec:
          containers:
            - name: sample-app
              image: nginx
              envFrom:
                - configMapRef:
                    name: app-config
        ```
        
- 예시 - 파일 형태로 마운트
    - ConfigMap
        
        ```bash
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: nginx-config
        data:
          nginx.conf: |
            events {}
            http {
              server {
                listen 80;
                location / {
                  return 200 'Hello ConfigMap!';
                }
              }
            }
        ```
        
    - Pod
        
        ```bash
        volumes:
          - name: config-volume
            configMap:
              name: nginx-config
        containers:
          - name: nginx
            image: nginx
            volumeMounts:
              - mountPath: /etc/nginx/nginx.conf
                subPath: nginx.conf
                name: config-volume
        ```
        

## Secret

- key-value 쌍으로 정의
- Base64로 인코딩하여 저장
- 비밀번호, 인증 토큰, TLS 키 등 민감한 데이터를 저장할 때 사용
- AWS Secrets Manager, Vault와 같은 Secret CSI Driver을 연동해서 사용 가능
- 예시 - 환경변수
    - ConfigMap
        
        ```bash
        apiVersion: v1
        kind: Secret
        metadata:
          name: db-secret
        type: Opaque
        data:
          username: ZGJ1c2Vy   # "dbuser" → base64 인코딩
          password: c2VjcmV0MTIz # "secret123"
        ```
        
    - Pod
        
        ```bash
        containers:
          - name: db-client
            image: postgres
            env:
              - name: DB_USER
                valueFrom:
                  secretKeyRef:
                    name: db-secret
                    key: username
              - name: DB_PASS
                valueFrom:
                  secretKeyRef:
                    name: db-secret
                    key: password
        ```
        
- 예시 - 도커 레지스트리 인증 정보
    - Secret
        
        ```bash
        apiVersion: v1
        kind: Secret
        metadata:
          name: regcred
        type: kubernetes.io/dockerconfigjson
        data:
          .dockerconfigjson: <Base64-encoded-json>
        ```
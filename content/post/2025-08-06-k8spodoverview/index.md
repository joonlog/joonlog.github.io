---
title: "Kubernetes Pod Overview"
date: 2025-08-06T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["kubernetes", "k8s", "pod", "initcontainer", "static pod"]
---


## Pod

- 컨테이너를 추상화 시켜 놓은 단위
- k8s에서 배포하는 최소 단위
- 하나의 Pod에는 여러 Container가 들어 갈 수 있다
- Pod는 하나의 IP를 갖는다

Multi-Container Pod

```bash
Liveness Probe
- httpGet
    - livenessProbe:
      httpGet:
        path: /
        port: 80
- tcpSocket
    - livenessProbe:
      tcpSocket:
        port: 22
- exec:
    - livenessProbe:
      exec:
        command:
        - ls
        - /tmp/datafile

Liveness Probe Parameter
- initialDelaySeconds: 0
- periodSeconds: 5
- timeoutSeconds: 1
- successThreshold: 1
- failureThreshold: 2
```

## initcontainer

앱 컨테이너가 동작하기 위한 전초 작업을 위한 컨테이너

- (WEB-DB) 환경의 WEB 컨테이너에서 DB 컨테이너가 기동되었는지를 위한 작업
- (Network) 네트워크 인프라 점검을 위한 전초 작업

```yaml
spec:
  containers:
  - name: myapp-container
    image: busybox:1.28
    command: ['sh', '-c', 'echo The app is running! && sleep 3600']
  initContainers:
  - name: init-myservice
    image: busybox:1.28
    command: ['sh', '-c', "until nslookup myservice.$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace).svc.cluster.local; do echo waiting for myservice; sleep 2; done"]
  - name: init-mydb
    image: busybox:1.28
    command: ['sh', '-c', "until nslookup mydb.$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace).svc.cluster.local; do echo waiting for mydb; sleep 2; done"]
```

## Infra/pause container

- 앱컨테이너의 인프라를 생성 및 삭제를 담당하는 컨테이너

## Static Pod

- /var/lib/kubelet/config.yaml
    - staticPodPath: /etc/kubernetes/manifests

## Pod 자원 할당

- 자원 할당
    - cpu(단위: 1core/1vCPU = 1000m)
    - mem(단위: 500Mi, 1Gi, ..)
    
    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
      name: nginx-pod-resources
    spec:
      containers:
      - name: nginx-container
        image: nginx:1.14
        ports:
        - containerPort: 80
          protocol: TCP
        resources:
          requests:
            cpu: 200m
            memory: 250Mi
          limits:
            cpu: 1
            memory: 1Gi
    ```
    

## Pod 환경 변수

- 환경 변수 설정

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: mysql-env-pod
spec:
  containers:
  - name: mysql-container
    image: mysql
    ports:
    - containerPort: 3306
      protocol: TCP
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: "password"
```

## 디자인 패턴

- 사이드카 - 협업
- 어댑터 - 입력 표준화
- 앰버서더 - 출력 표준화
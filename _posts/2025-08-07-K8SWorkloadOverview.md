---
title : Kubernetes Workload Overview
date : 2025-08-07 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, replicaset, deployment, daemonset, statefulset, job]  #소문자만 가능
---

## ReplicaSet

- 동일한 스펙의 Pod를 지정된 개수만큼 유지
- 집합 기반 셀렉터(Set-based selector)
- rolling update 지원
- 직접 쓰기보다는 Deployment 내부에서 주로 사용됨

```bash
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: my-replicaset
spec:
  replicas: 3
  selector:
    matchExpressions:
      - key: app
        operator: In
        values: ["web", "api"]
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
        - name: nginx
          image: nginx:1.27
```

## Deployment

- ReplicaSet을 관리하며 Pod의 개수를 보장
- Pod에 대한 rolling update/rollback 기능 지원
- 무중단 파드 배포를 위해 사용

```bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
        - name: nginx
          image: nginx:1.27
```

## DaemonSet

- 하나의 노드당 1개의 Pod 실행(controlplane 제외)
- Rolling Update / Rollback 기능 지원
- 데몬셋 용도: (관리) ex) 모니터링 Pod 실행

```bash
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-monitor
spec:
  selector:
    matchLabels:
      app: node-monitor
  template:
    metadata:
      labels:
        app: node-monitor
    spec:
      containers:
        - name: node-exporter
          image: prom/node-exporter:latest
```

## StatefulSet

- State가 필요한 워크로드에 사용(Pod 이름, 스토리지 고정)
- 고유 네트워크 ID, 안정적인 스토리지, 순차 배포/종료 보장
- 주로 DB, 메시지 큐, Zookeeper, Kafka 등에 사용

```bash
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web-stateful
spec:
  serviceName: "web"
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
        - name: nginx
          image: nginx:1.27
  volumeClaimTemplates:
    - metadata:
        name: www
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 1Gi
```

## Job

- 트랜잭션의 단위
- 잡이 성공적으로 종료 될 때 까지 지속적으로 파드 실행
- 잡이 성공하면 completed
- 잡은 실패하면 종료 후 다시 실행

```bash
apiVersion: batch/v1
kind: Job
metadata:
  name: data-migration
spec:
  backoffLimit: 3
  template:
    spec:
      containers:
        - name: migrate
          image: busybox
          command: ["sh", "-c", "echo Running migration... && sleep 10"]
      restartPolicy: OnFailure
```

## CronJob

- Job 컨트롤러를 사용하여 애플리케이션 파드를 주기적으로 반복해서 실행
- Job이 종료되면 스케쥴링에 의해 주기적으로 실행
- Job이 비정상 종료되면 다시 실행

```bash
apiVersion: batch/v1
kind: CronJob
metadata:
  name: db-backup
spec:
  schedule: "0 2 * * *"  # 매일 새벽 2시
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: backup
              image: busybox
              command: ["sh", "-c", "echo Backing up DB... && sleep 5"]
          restartPolicy: OnFailure
```
---
title : Kubernetes Workload Overview
date : 2025-08-07 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, replicaset, deployment, daemonset, statefulset, job]  #소문자만 가능
---

## ReplicaSet

- 집합 기반 셀렉터(Set-based selector)
- rolling update 지원

## Deployment

- Pod의 개수를 보장
- Pod에 대한 rolling update/rollback 기능 지원

## DaemonSet

- 하나의 노드당 1개의 Pod 실행(controlplane 제외)
- Rolling Update / Rollback 기능 지원
- 데몬셋 용도: (관리) ex) 모니터링 Pod 실행

## StatefulSet

[참고] Deployment vs StatefulSet

## Job

- 트랜잭션의 단위
- 잡이 성공적으로 종료 될 때 까지 지속적으로 파드 실행
- 잡이 성공하면 completed
- 잡은 실패하면 종료 후 다시 실행

## CronJob

- Job 컨트롤러를 사용하여 애플리케이션 파드를 주기적으로 반복해서 실행
- Job이 종료되면 스케쥴링에 의해 주기적으로 실행
- Job이 비정상 종료되면 다시 실행
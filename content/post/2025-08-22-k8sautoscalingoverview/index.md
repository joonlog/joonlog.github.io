---
title: "Kubernetes Autoscaling Overview"
slug: "k8s-autoscaling-overview"
aliases: ["/p/kubernetes-autoscaling-overview/"]
date: 2025-08-22T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["Kubernetes"]
---

## 오토스케일링 종류

- 클러스터 레벨 오토스케일링

    - 노드 풀의 VM 개수를 오토스케일링

    - Cluster Autoscaler, Karpenter 등. 노드를 만들 API가 있는 환경(EKS, GKE, AKS 등)에서 주로 사용

- 파드 레벨 오토스케일링

    - 파드의 개수를 오토스케일링

### 수평적 파드 오토스케일링(HPA)

- metrics-server 애드온이 필요

    - 각 노드/파드가 사용하는 매트릭(CPU/MEM) 모니터링하고 데이터 수집하는 역할

동작원리

- HPA 컨트롤러가 metrics-server에서 Pod 매트릭을 주기적으로 수집하고, 사전에 설정한 임계치 조건을 만족하면 파드의 개수를 자동으로 조정

HPA 동작 조건

- 기본값 15초 간격으로 Pod 사용량을 점검(`--horizontal-pod-autoscaler-sync-period`)

- 임계값을 초과하면 Pod scale out

- scale out은 바로 실행(안정화 창 0초), scale in은 최근 5분(안정화 창 300초) 동안의 가장 높은 권장값을 기준으로 천천히 줄임. HPA의 `behavior`로 조정 가능

- 예시

    - `myapp` Deployment가 CPU 사용 평균 50%를 넘기면 Pod 개수가 늘어남

    - 만약 CPU 사용률이 낮아지면, 5분 이후 Pod 개수를 줄임

    - Pod 개수는 최소 2개, 최대 10개 사이에서 조정됨

    ```bash
    apiVersion: autoscaling/v2
    kind: HorizontalPodAutoscaler
    metadata:
      name: myapp-hpa
    spec:
      scaleTargetRef:  
        apiVersion: apps/v1
        kind: Deployment
        name: myapp
      minReplicas: 2 
      maxReplicas: 10  
      metrics:
      - type: Resource
        resource:
          name: cpu
          target:
            type: Utilization
            averageUtilization: 50 
    ```

---

애플리케이션 지표 기반 오토스케일링도 가능하지만, prometheus와 같은 툴과 연동 필요(Prometheus Adapter, 이벤트 기반이면 KEDA)

수직 오토스케일링(VPA)은 Pod 개수 대신 request/limit을 조정

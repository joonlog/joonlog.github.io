---
title : Kubernetes Autoscaling Overview
date : 2025-08-22 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, autoscaling, hpa, ca]  #소문자만 가능
---

## 오토스케일링 종류

- 클러스터 레벨 오토스케일링
    - 노드 풀의 VM 개수를 오토스케일링
    - 클라우드 프로바이더와 같은 환경에서만 가능(EKS, GKE, AKS 등)
- 파드 레벨 오토스케일링
    - 파드의 개수를 오토스케일링

### 수평적 파드 오토스케일링(HPA)

- metrics-server 애드온이 필요
    - 각 노드/파드가 사용하는 매트릭(CPU/MEM) 모니터링하고 데이터 수집하는 역할

동작원리

- HPA 컨트롤러가 metrics-server에서 Pod 매트릭을 주기적으로 수집하고, 사전에 설정한 임계치 조건을 만족하면 파드의 개수를 자동으로 조정

HPA 동작 조건

- 기본값 30초 간격으로 Pod 사용량을 점검
- 임계값을 초과하면 Pod scale out
- 확장된 이후 3분 대기, scale in은 확장 이후 5분 대기
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

애플리케이션 지표 기반 오토스케일링도 가능하지만, prometheus와 같은 툴과 연동 필요
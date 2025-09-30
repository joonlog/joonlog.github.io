---
title : Kubernetes QoS Overview
date : 2025-09-30 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, qos, guaranteed, burstable, bestEffort, eviction]  #소문자만 가능
---

> kubernetes는 노드에 리소스가 부족할 때, Pod의 QoS 클래스에 따라 어떤 Pod를 우선적으로 Eviction(퇴출)할지를 결정
> 

### QoS 클래스 종류

1. Guaranteed
    - Pod 내 모든 컨테이너의 CPU/Memory Request와 Limit 값이 동일한 경우
    - 가장 높은 우선순위를 가짐 (Eviction 대상에서 가장 안전함)
2. Burstable
    - Guaranteed와 BestEffort에 속하지 않는 모든 경우
    - 예: Request < Limit 인 경우
    - 우선순위는 중간
3. BestEffort
    - 리소스 Request와 Limit을 전혀 지정하지 않은 경우
    - 가장 낮은 우선순위 (Eviction 1순위)

---

### Eviction 동작 원리

노드에 리소스가 부족해지면, kubelet은 다음과 같은 순서로 Pod를 Eviction

- 우선순위: Guaranteed > Burstable > BestEffort
- BestEffort
    - 항상 Eviction 대상에 포함됨
- Burstable
    - CPU/Memory Request보다 실제 사용량이 큰 경우 Eviction 대상
- Guaranteed
    - Eviction 대상에서 가장 안전

> 주의: Request와 Limit 차이를 너무 크게 잡으면 Burstable로 분류되어 Eviction 대상
> 
> 
> 따라서 Request와 Limit은 가급적 동일하게 잡는 것이 안정적
> 

---

### CPU vs Memory

- CPU → 압축 가능한(compressible) 리소스
    - Request보다 더 쓰더라도 kubelet이 **throttle(스로틀링)** 해서 CPU 사용량을 줄일 수 있음
    - 즉, CPU 경합이 발생해도 프로세스가 죽지는 않음
- Memory, Storage → 압축 불가능한(incompressible) 리소스
    - Memory 사용이 Request를 초과하면, kubelet은 해당 Pod를 **Eviction**시켜야 함
    - Storage도 마찬가지로 강제 종료가 필요할 수 있음

### EX) Guaranteed Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: guaranteed-pod
spec:
  containers:
  - name: nginx
    image: nginx
    resources:
      requests:
        memory: "512Mi"
        cpu: "500m"
      limits:
        memory: "512Mi"
        cpu: "500m"

```

Request = Limit 이므로 QoS = Guaranteed

### EX) Burstable Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: burstable-pod
spec:
  containers:
  - name: nginx
    image: nginx
    resources:
      requests:
        memory: "256Mi"
        cpu: "250m"
      limits:
        memory: "512Mi"
        cpu: "500m"

```

- Request < Limit 이므로 QoS = Burstable

### EX) BestEffort Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: besteffort-pod
spec:
  containers:
  - name: nginx
    image: nginx
```

- Request/Limit 미지정 → QoS = BestEffort

### Pod QoS Class 확인 방법

```bash
kubectl get pod {POD_NAME} -o jsonpath='{ .status.qosClass}{"\n"}'
```

### 정리

- QoS 클래스는 Pod의 리소스 보장 수준
- Eviction 순서: Guaranteed > Burstable > BestEffort
- CPU는 throttle 가능(죽지 않음), Memory/Storage는 Eviction 발생
- 안정적인 서비스 운영을 위해 Request와 Limit을 동일하게 설정하는 것이 가장 안전

---

### 참고 자료

- [Kubernetes 공식 문서 - QoS](https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/)
- [Kubernetes 공식 문서 - QoS 설정 예제](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/)
- [No Easy Dev 블로그](https://no-easy-dev.tistory.com/entry/%EC%BF%A0%EB%B2%84%EB%84%A4%ED%8B%B0%EC%8A%A4-Pod-%EC%9E%90%EC%9B%90%EA%B4%80%EB%A6%AC-QoS)
- [ssup2 - Pod Eviction 분석](https://ssup2.github.io/theory_analysis/Kubernetes_Pod_Eviction/)
- https://jenakim47.tistory.com/96
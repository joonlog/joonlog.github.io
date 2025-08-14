---
title : Kubernetes Label과 Annotation
date : 2025-08-10 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, label, annotation]  #소문자만 가능
---

## Label

- label은 3개 이상 정의할 것을 권장
- environment=prod tier=frontend release=stable team=devops1

## Annotation

- annotation은 key, value 쌍으로 정의
- 쿠버네티스에게 특정 정보를 전달할 용도로 사용
ex) annotation:
kubernetes.io/change-cause: nginx:1.14
- 관리를 위해 필요한정보를 기록할 용도로 사용
ex) annotations:
builder: BSC
buildData: “2024-0403”
imageregistry: “https:…”
- [참고] Labels vs Annotations
    - lables: 클러스터를 관리할 때 사용자가 원하는 값을 입력
    - annotations: 클러스터 시스템에 필요한 정보를 표시하는데 사용

## 배포 정책 패턴/업데이터 정책 패턴

새로운 파드 배포 방법

- 블루-그린 배포: 한꺼번에 업데이트
- 카나리 배포: 일부만 업데이트
- 롤링 업데이트 배포: 순차적으로 업데이트
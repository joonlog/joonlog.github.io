---
title: "K8S에 구축한 GitLab Registry에 Container Image 푸시 시 blob upload unknown 에러 TroubleShooting"
slug: "k8s-gitlab-registry-blob-upload-troubleshooting"
aliases: ["/p/k8s에-구축한-gitlab-registry에-container-image-푸시-시-blob-upload-unknown-에러-troubleshooting/"]
date: 2025-09-27T09:00:00+09:00
categories: ["Kubernetes", "GitLab"]
tags: ["Kubernetes", "GitLab", "MinIO", "Longhorn", "Self-Managed Kubernetes"]
---

> GitLab Helm Chart에서 GitLab Registry 데이터 볼륨의 persistence 설정 기본값이 object storage로 되어 있기 때문에 발생한 이슈

### 문제 상황

- K8S 환경에 Helm Chart로 구축한 GitLab 서버의 Container Registry에 컨테이너 이미지를 푸시하면 `blob upload unknown` 에러 발생

    - docker push 하면 어느정도 레이어가 푸시되는 것 같다가 동일한 구간에서 계속 실패

### 에러 분석

- GitLab Registry 파드는 HPA replicas에 의해 여러 파드가 동작하도록 설정되어 있는데, 이 파드 간 데이터 공유가 안되는 것 같다고 느낌

- GitLab Registry가 실제로 데이터를 어디에 저장하는지 확인해보니 각 파드의 /tmp 경로에 임시로만 저장하고 있었다!!!

    - Helm Chart values.yaml 파일에서 Registry에 바로 Longhorn PV를 사용하도록 설정했는데, Registry 설정에 직접적으로 PV를 부여하지 못하도록 제한되어 있는 것을 확인

    - 찾아보니 minio라는 오브젝트 스토리지가 values.yaml에 있었고, minio는 GitLab 서비스들의 바이너리 파일들을 저장할 수 있도록 설계된 것을 확인

    - Registry에 데이터를 저장하려면 object stroage를 통해서 저장하는 설정이 필수!

### 해결 방법

- minio에 Longhorn PV를 할당하고, 활성화한 minio를 registry가 사용할 수 있도록 설정

- minio설정

    - https://docs.gitlab.com/charts/charts/globals#configure-minio-settings

    - S3와 호환되는 오브젝트 스토리지

```bash
global:

  # minio 서비스 활성화
  minio:
    enabled: true

  # 바이너리 파일들이 오브젝트 스토리지를 사용할 수 있도록 설정
  appConfig:
    object_store
      enabled: true

  # Container Registry가 minio의 registry 버킷을 사용
  registry:
    bucket: registry

# minio가 Longhorn을 통한 PV 사용
minio:
  persistence:
    enabled: true
    size: 10Gi
    storageClass: longhorn
```

- 설정 후 docker push 성공!

### Omnibus GitLab과 K8S GitLab의 차이점

Omnibus GitLab:

- 단일 서버에 모든 컴포넌트가 함께 설치

- Registry가 로컬 파일시스템 직접 사용

- 데이터 영속성 보장됨

- Object storage 선택사항

K8S GitLab:

- Registry가 별도 Pod로 실행

- Pod 재시작 시 임시 파일시스템 초기화됨

- PersistentVolume 사용해도 단일 Registry Pod에 종속

- 여러 Registry Pod 스케일링 불가 (파일시스템 공유 안됨)

- 실질적으로 object storage 필수

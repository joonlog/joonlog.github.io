---
title: "Kubernetes Storage Overview"
date: 2025-08-21T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["kubernetes", "k8s", "pv", "pvc", "storageclass"]
---


- Pod에서 볼륨 선언 → 마운트가 기본 패턴
    - spec.volumes에서 볼륨 선언 후 spec.containers.volumeMounts 에서 마운트

```bash
apiVersion: v1
kind: Pod
metadata:
  name: demo-pod
spec:
  containers:
    - name: app
      image: nginx:1.25
      volumeMounts:
        - name: html
          mountPath: /usr/share/nginx/html
  volumes:
    - name: html
      hostPath:
        path: /home/
        type: DirectoryOrCreate
```

## 볼륨 유형

| 타입 | 사용 용도 | 비고 |
| --- | --- | --- |
| emptyDir | Pod 수명=데이터 수명. 임시 공간 | Pod 종료 시 데이터 없어짐 |
| hostPath | 노드 로컬 경로를 Pod에 마운트 | 데이터가 특정 워커노드에 종속됨 |
| configMap / secret | 변수를 파일로 주입 | 주로 읽기 전용으로 마운트 |
| nfs | NAS 마운트 | 멀티 Pod 간 데이터 공유에 사용 |
| CSI 기반 블록 / 파일 | AWS EDB, GCE PD 등 | 대부분 CSI 드라이버로 제공
동적 프로비저닝 가능 |
- 예시: emptyDir / configMap / secret

```bash
# emptyDir
volumes:
  - name: work
    emptyDir: {}  

# configMap
volumes:
  - name: app-config
    configMap:
      name: my-config

# secret
volumes:
  - name: app-secret
    secret:
      secretName: db-cred
```

## PV & PVC

- PV(PersistentVolume): 클러스터에 영구 스토리지 리소스를 생성
- PVC(PersistentVolumeClaim): 생성한 PV에서 자원 할당 요청

접근 모드

- ReadWriteOnce(RWO): 한 노드에서 읽기/쓰기
- ReadOnlyMany(ROX): 여러 노드에서 읽기 전용
- ReadWriteMany(RWX): 여러 노드에서 읽기/쓰기
    - ex) NFS, EFS
- pv/pvc 예시
    - `persistentVolumeReclaimPolicy` 기본값은 동적 프로비저닝이면 `Delete`인 경우가 많고, 중요한 데이터면 `Retain` 고려

```bash
# 1) PV: NFS 예시
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-nfs
spec:
  capacity:
    storage: 20Gi
  accessModes: [ReadWriteMany]
  persistentVolumeReclaimPolicy: Retain
  nfs:
    path: /export/apps
    server: 10.0.0.10

---
# 2) PVC: 위 PV에 바인딩될 요청
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-nfs
spec:
  accessModes: [ReadWriteMany]
  resources:
    requests:
      storage: 20Gi
```

## StorageClass

- 디스크 타입, 복제/백업 정책, 토폴로지 등의 프로비저닝에 필요한 파라미터를 묶어서 제공
- 사용자는 PVC에서 StorageClassName에만 지정하면 볼륨이 동적 프로비저닝됨
- 사용을 위해선 NFS Provisioner, EBS Driver와 같은 툴 설치 필요
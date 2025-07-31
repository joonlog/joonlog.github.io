---
title : 로컬 Kubernetes에서 NFS를 통한 StorageClass 사용
date : 2025-07-31 09:00:00 +09:00
categories : [Kubernetes, Plugins]
tags : [kubernetes, storageclass, nfs, pvc]  #소문자만 가능
---

- AWS의 경우 EBS Driver를 통해 pvc를 위한 StorageClass를 사용 가능
    - 로컬 Kubernetes에서도 StorageClass를 사용하기 위해 NFS 서버를 구축해서 사용
- NFS 서버 구축

```bash
dnf install -y nfs-utils

systemctl enable --now nfs-server

mkdir -p /exports/k8s-pv
chown -R nobody:nobody /exports/k8s-pv
chmod 777 /exports/k8s-pv
echo "/exports/k8s-pv *(rw,sync,no_subtree_check,no_root_squash)" | sudo tee -a /etc/exports
exportfs -rav
```

- Helm을 통한 NFS Provisioner 설치

```bash
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
helm repo update

helm install nfs-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner   --set nfs.server=172.27.1.9   --set nfs.path=/exports/k8s-pv   --set storageClass.name=nfs-sc   --set storageClass.defaultClass=true
```

- pvc 생성

```bash
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
  storageClassName: nfs-sc
```

- pod 생성

```bash
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
  storageClassName: nfs-sc
```
---
title: "Linux에서 디스크 추가 시에 LVM 증설"
date: 2025-11-28T09:00:00+09:00
categories: ["Linux", "System"]
tags: ["linux", "lvm"]
---


> 20GB 추가디스크가 서버에 부탁됐을 때 디스크를 LVM으로 /data 경로에 마운트하는 방법
> 

### 0. 디스크 초기화

```bash
fdisk /dev/xvdb
n p 1 enter enter t 8e 2
```

### 1. PV 생성

```bash
pvcreate /dev/xvdb1
```

### 2. VG 생성

```bash
vgcreate vgdata /dev/xvdb1
```

### 3. LV 생성

```bash
# 용량 지정해서 사용
lvcreate -L 20G -n lvdata vgdata

# 남은 전체 공간 사용
lvcreate -l 100%FREE -n lvdata vgdata
```

### 4. 파일시스템 생성

```bash
mkfs.ext4 /dev/vgdata/lvdata
```

### 5. 마운트 및 fstab 등록

```bash
mkdir -p /data
mount /dev/vgdata/lvdata /data

echo "/dev/vgdata/lvdata /data ext4 defaults 0 2" >> /etc/fstab
```
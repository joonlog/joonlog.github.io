---
title: "Linux LVM 명령어"
date: 2025-11-28T09:00:00+09:00
categories: ["Linux", "System"]
tags: ["linux", "lvm"]
---


> 20GB 추가디스크를 LVM으로 구성하는 방법
> 

### 1. PV 생성

```bash
pvcreate /dev/xvdb
```

### 2. VG 생성

```bash
vgcreate vgdata /dev/xvdb
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
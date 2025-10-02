---
title : Linux Swap 메모리 증설 방법
date : 2025-10-02 09:00:00 +09:00
categories : [Linux, System]
tags : [linux, swap memory]  #소문자만 가능
---

```bash
# 현재 스왑 메모리 확인
free -mh

# Swap 파일 생성
fallocate -l 2G .swapfile 
# fallocate 안되면 dd
dd if=/dev/zero of=/swapfile2 bs=1M count=2048

# Swap 파일 권한 변경
chmod 600 .swapfile

# Swap 파일 포맷 설정 및 Swap 파일 등록
mkswap .swapfile
swapon .swapfile 

# Swap Memory 가 정상적으로 등록 되었는지 확인
swapon -s
free -m

# 부팅시에도 Swap 연결이 되도록 Swap 파일 영구적 마운트 하는 방법
vim /etc/fstab
  .swapfile swap swap defaults 0 0

# Swap Memory 비활성화 및 삭제
swapoff .swapfile
rm .swapfile

# Swap 파일 마운트 해제
vim /etc/fstab
  [swap파일경로] swap swap defaults 0 0 << 문구 주석 처리
```
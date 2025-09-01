---
title : Linux 환경에서 Docker 설치 (RedHat/Ubuntu 계열)
date : 2025-09-01 09:00:00 +09:00
categories : [Container, Docker]
tags : [container, docker, redhat, centos, ubuntu] #소문자만 가능
---

- https://docs.docker.com/engine/install

### RedHat 계열

```bash
# 기존 설치되어있는 컨테이너 런타임 삭제
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo dnf remove $pkg -y; done

# docker repo 추가
yum install -y yum-utils
yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
# docker 설치
yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 현재 사용자를 docker 그룹에 추가해서 root 권한 없어도 docker 사용 가능하게 설정
usermod -a -G docker <user>

systemctl enable --now docker
```

### Ubuntu 계열

```bash
# 기존 설치되어있는 컨테이너 런타임 삭제
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done

# docker repo 설정
apt-get update
apt-get install ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
# docker repo 추가
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

# docker 설치
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 현재 사용자를 docker 그룹에 추가해서 root 권한 없어도 docker 사용 가능하게 설정
usermod -a -G docker <user>

systemctl enable --now docker
```
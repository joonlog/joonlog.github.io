---
title: "Rocky8 Yum Proxy"
date: 2024-12-20T09:00:00+09:00
categories: ["Linux", "Package"]
tags: ["rocky8", "yum", "proxy"]
---


### 목표

- pub 1ea, priv 1ea 서버 생성 후 yum proxy 통하여 priv 서버에서 yum 설치 가능하게 설정
- 외에도 pub 서버에서 yumdownloader 명령어 통해서 rpm 빼와서 priv 서버로 옮긴 후 rpm 통한 설치
    - mariadb 최신버전 설치로 테스트

### 환경 설정

- KT Cloud
    - public 1대
        - 공인 IP 포트포워딩: 2222→22
    - private 1대
        - 공인 IP 포트포워딩: X
- OS: Rocky Linux 8.1

## 작업 과정

## 1. Yum Proxy

### 1) Public 서버

- 사용자 생성

```bash
useradd -m -G wheel yumpub
echo "yumpub" | passwd --stdin yumpub
```

- 방화벽 설정

```bash
sudo systemctl enable --now firewalld 
sudo firewall-cmd --permanent --add-port=3128/tcp
sudo firewall-cmd --reload
```

- squid 설치

```bash
sudo dnf update --exclude=kernel* -y
sudo dnf install -y squid
sudo systemctl enable --now squid
```

- squid.conf 수정

```bash
sudo vi /etc/squid/squid.conf

# private server 허용
acl priv_network src <pri-server-ip>/24
http_access allow priv_network
http_access deny all

sudo systemctl restart squid
```

### 2) Private 서버

- 사용자 생성

```bash
ssh root@<pri-server-ip>

useradd -m -G wheel yumpri
echo "yumpri" | passwd --stdin yumpri
exit

ssh yumpri@<pri-server-ip>
```

- 방화벽 설정

```bash
sudo systemctl enable --now firewalld 
sudo firewall-cmd --reload
```

- yum.conf 수정
    - pub 서버를 프록시로 사용

```bash
sudo vi /etc/yum.conf

proxy=http://<public-server-ip>:3128
```

- 기존에 하던 대로 dnf update 혹은 repo 생성해서 패키지 설치
    - repo 생성

```bash
sudo vi /etc/yum.repos.d/Rocky-Extras.repo

[extras]
name=Rocky Linux $releasever - Extras
mirrorlist=https://mirrors.rockylinux.org/mirrorlist?arch=$basearch&repo=extras-$releasever
gpgcheck=1
enabled=1
countme=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-rockyofficial
```

- extras의 패키지 epel-release 검색

```bash
sudo dnf clean all
sudo dnf repolist
sudo dnf search epel-release
```

## 2. yumdownloader

### 1) Public 서버

- yumdownloader로 rpm 다운로드

```bash
sudo dnf update --exclude=kernel* -y
sudo dnf install -y yum-utils perl
# --resolve는 의존성 포함
yumdownloader --resolve mariadb
yumdownloader --resolve mariadb-server
yumdownloader --resolve perl-DBD-MySQL perl-DBI
```

- rpm 파일 private 서버로 복사

```bash
scp mariadb*.rpm yumpri@<priv-server-ip>:/tmp
scp perl*.rpm yumpri@<priv-server-ip>:/tmp
```

### 2) Private 서버

- rpm 파일로 mariadb 설치

```bash
sudo rpm -ivh /tmp/perl*.rpm
sudo rpm -ivh /tmp/mariadb*.rpm

sudo rpm -qa | grep mariadb
```

- MariaDB 테스트

```bash
mysql --version

sudo systemctl enable --now mariadb

mysql_secure_installation 

mysql -u root
```

---

## 로그 확인

- squid 로그

```bash
sudo tail -f /var/log/squid/access.log
```

## 트러블슈팅

> yum proxy 실패
> 
- /var/log/squid/access.log 확인 시 priv 서버 ip 가 없고 모르는 사설 ip가 있음
    - tail 걸고 priv에서 yum search했을 때 로그가 발생한걸로 봐서 모르는 사설 ip가 priv 서버의 트래픽인 것으로 확인
    - /etc/yum.conf에 공인 IP로 설정해서 공인 IP에서 pub 서버로 트래픽이 이동하는 과정에서 자동으로 NAT되어 새로운 사설 ip가 할당 된 것으로 보임
    - 애초에 yum.conf에 공인 IP 말고 사설 IP를 설정해야 했지만 proxy가 안되는 원인은 아님
- yum은 자신이 활성화한 repo만 squid 프록시를 통해 접근
    - priv에 repo가 없어도 pub에 있는 repo를 통해 priv에 설치가 되는 줄 앎
    - 그게 아니고, squid는 단순히 proxy만 해주고 pub는 squid 서버의 역할(NAT 서버)
    - 기존에 하던 대로 dnf update 하거나 repo 만들어서 해결

## 참고

yum proxy:
https://hgko1207.github.io/2020/09/28/linux-3/

yumdownloader:
https://wldnjd2.tistory.com/53
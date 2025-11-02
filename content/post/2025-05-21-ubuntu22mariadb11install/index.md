---
title: "Ubuntu22 MariaDB 11.4.5 버전 설치"
date: 2025-05-21T09:00:00+09:00
categories: ["Linux", "DB"]
tags: ["linux", "db", "mysql", "mariadb", "ubuntu22", "mariadb 11.4.5"]
---

### 방법 1: 수동으로 저장소 구성 후 설치

- MariaDB의 서명 키와 저장소를 수동으로 등록한 뒤 패키지를 설치하는 방식

```bash
sudo apt-get install apt-transport-https curl
sudo mkdir -p /etc/apt/keyrings
sudo curl -o /etc/apt/keyrings/mariadb-keyring.pgp 'https://mariadb.org/mariadb_release_signing_key.pgp'
```

- APT 저장소 설정 파일을 `/etc/apt/sources.list.d/mariadb.sources`에 생성

```
X-Repolib-Name: MariaDB
Types: deb
URIs: https://tw1.mirror.blendbyte.net/mariadb/repo/11.4/ubuntu
Suites: jammy
Components: main main/debug
Signed-By: /etc/apt/keyrings/mariadb-keyring.pgp
```

- 리눅스 커널 관련 패키지 자동 업데이트를 방지하고 설치

```bash
apt-mark hold linux-image-generic linux-headers-generic
apt-get install mariadb-server
```

---

### 방법 2: MariaDB 공식 설치 스크립트 사용

- MariaDB에서 제공하는 설치 스크립트를 통해 저장소 등록과 설치를 한 번에 처리하는 방식

```bash
curl -LsS https://downloads.mariadb.com/MariaDB/mariadb_repo_setup | bash -s -- --mariadb-server-version=11.

apt-mark hold linux-image-generic linux-headers-generic
apt install mariadb-server mariadb-client -y
```

---

### 초기 공통 설정

- mysql_secure_installation

```bash
mysql_secure_installation

Enter current password for root (enter for none):
# 현재 root 계정의 비밀번호
# 비밀번호가 설정되어 있지 않은 경우 엔터

Switch to unix_socket authentication [Y/n]:
# root 로그인을 비밀번호가 아니라 시스템 계정(sudo)으로만 허용할지 결정
# Y를 선택하면 root로 로그인할 때 비밀번호 대신 sudo mysql로 접속
# 일반적으로는 Y(보안성 향상). 비밀번호로 로그인할 수 있게 하려면 n 입력.

Change the root password? [Y/n]:
# root 계정의 비밀번호를 변경할지 여부
# 새로 설정하거나, 이전 비밀번호가 없는 경우 새로 설정하는 데 사용

Remove anonymous users? [Y/n]:
# 익명 사용자 계정을 삭제할지 여부
# 익명 계정은 인증 없이 접속할 수 있으므로 Y 입력 권장

Disallow root login remotely? [Y/n]:
# root 계정으로 원격 접속을 막을지 여부
# 보안상 Y 권장. 필요한 경우 n 선택 후 방화벽 등으로 제한하는 방식 고려

Remove test database and access to it? [Y/n]:
# 테스트용 test 데이터베이스를 삭제할지 여부
# 기본으로 생성되는 샘플 데이터베이스로, 불필요하므로 Y 입력 권장

Reload privilege tables now? [Y/n]:
# 앞서 설정한 변경 사항들을 즉시 적용할지 여부
# Y 선택하면 권한 테이블이 즉시 갱신
```

- 계정 생성 및 권한 부여

```bash
mysql -u root
CREATE DATABASE test;
CREATE USER 'admin'@'localhost' IDENTIFIED BY '비밀번호';
GRANT ALL PRIVILEGES ON *.* TO 'admin'@'localhost' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

- binary log 활성화

```bash
vim /etc/mysql/my.cnf

[mysqld]
log-bin = mariadb-bin
binlog-format = MIXED
server-id = 1

systemctl restart mariadb
SHOW VARIABLES LIKE 'log_bin';
```

---

- MariaDB를 완전 제거

```bash
apt purge mariadb-server mariadb-client
rm -rf /etc/mysql /var/lib/mysql
```
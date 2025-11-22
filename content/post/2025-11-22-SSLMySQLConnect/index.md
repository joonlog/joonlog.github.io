---
title: "SSL 인증서를 사용한 MySQL/MariaDB 암호화 통신"
date: 2025-11-22T09:00:00+09:00
categories: ["Linux", "DB"]
tags: ["linux", "db", "mysql", "mariadb", "ca connect", "ssl connect"]
---


> App-DB 연동 간 사설 IP로 연동 불가 +  VPN 사용 불가한 상황에서 공인 IP를 통한 암호화 통신 설정하는 방법
> 
- App 서버에서 DB 서버로 공인 IP로 통신할 경우 평문으로 통신하기에 데이터 탈취 위험이 존재
- SSL 인증서를 사용한 암호화 통신으로 데이터 보호

### 환경

App 서버: PHP-FPM 8.3

DB 서버: MariaDB 10.3

연결 방식: App 서버 → DB 서버 공인 IP

### 작업 과정

1. DB 서버에서 SSL 인증서로 사용할 CA 인증서 생성
2. DB에서 인증서를 사용하도록 설정
3. 생성한 인증서를 앱 서버로 복사
4. 앱 서버에서 인증서를 사용하도록 설정
    - PHP DB 연동 설정 + DB 연동 소스 코드 수정 필요

## 전체 과정

### 1. DB 서버 - SSL 인증서 생성 및 설정

1. 디렉토리 생성
    
    ```bash
    mkdir -p /etc/mysql/ssl
    cd /etc/mysql/ssl
    ```
    
2. CA Private Key / 인증서 생성
    
    ```bash
    openssl genrsa 2048 > ca-key.pem
    openssl req -new -x509 -days 3650 -key ca-key.pem -out ca-cert.pem -subj "/CN=MariaDB-CA"
    ```
    
3. 서버용 인증서 생성
    
    ```bash
    openssl genrsa 2048 > server-key.pem
    openssl req -new -key server-key.pem -out server-req.pem -subj "/CN=MariaDB-Server"
    openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem \
      -CAcreateserial -out server-cert.pem -days 3650
    ```
    
- 권한 설정
    
    ```bash
    chmod 600 server-key.pem
    chown mysql:mysql server-key.pem server-cert.pem ca-cert.pem
    ```
    
1. MariaDB 설정 추가
    - `/etc/my.cnf.d/server.cnf`
    
    ```bash
    [mysqld]
    ssl-ca=/etc/mysql/ssl/ca-cert.pem
    ssl-cert=/etc/mysql/ssl/server-cert.pem
    ssl-key=/etc/mysql/ssl/server-key.pem
    ```
    
2. MariaDB 재시작
    
    ```bash
    systemctl restart mariadb
    ```
    
3. SSL 활성화 확인
    
    ```bash
    SHOW VARIABLES LIKE 'have_ssl';
    SHOW SESSION STATUS LIKE 'Ssl_cipher';
    SHOW SESSION STATUS LIKE 'Ssl_version';
    ```
    

### 2. 앱 서버 - CA 인증서 적용 및 PHP SSL 사용 설정

1. 디렉토리 생성
    
    ```bash
    mkdir -p /etc/mysql/ssl
    ```
    
2. DB 서버에서 CA 인증서 복사
    - `/etc/mysql/ssl/ca-cert.pem`
3. 환경파일 설정
    - .env 파일일 경우
    
    ```bash
    DB_SSL=true
    DB_SSL_CA=/etc/mysql/ssl/ca-cert.pem
    DB_SSL_VERIFY=false
    ```
    
4. PHP SSL 옵션 적용
    - PDO로 DB와 연동할 경우
        - DB.php
        
        ```bash
        if (getenv('DB_SSL') === 'true') {
            if ($ca = getenv('DB_SSL_CA')) {
                $options[\PDO::MYSQL_ATTR_SSL_CA] = $ca;
            }
            if (getenv('DB_SSL_VERIFY') === 'false') {
                $options[\PDO::MYSQL_ATTR_SSL_VERIFY_SERVER_CERT] = false;
            }
        }
        ```
        
5. php-fpm 재시작

### 3. SSL 적용 검증

1. DB 계정에 SSL 강제 적용하도록 설정
    - 애플리케이션이 SSL 적용이 정상적이지 않다면 에러 발생
        - EX) `SQLSTATE[HY000] [1045] Access denied for user`
    
    ```bash
    ALTER USER 'gill'@'%' REQUIRE SSL;
    FLUSH PRIVILEGES;
    ```
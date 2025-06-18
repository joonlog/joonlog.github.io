---
title : Rocky8 Apache SSL HTTPS Redirection
date : 2024-12-11 18:00:00 +09:00
categories : [Linux, Middleware]
tags : [rocky8, apache, ssl, https] #소문자만 가능
---

### 목표

- **Let’s encrypt**에서 ssl 인증서 발급해 전송 보안 강화
- **apache virtual host**로 http -> https 리디렉션 설정
    - 웹사이트의 모든 트래픽이 안전한 프로토콜을 사용하도록 강제
- 인증서의 **자동 갱신**을 설정하여 만료되지 않도록 유지

### 환경 설정

- KT Cloud
    - SSL 1대
    - 공인 IP 포트포워딩
- OS: Rocky Linux 8.1
- WEB: Apache 2.4.6
- SSL: Let’s encrypt certbot

## 개념 및 설정

### Let’s encrypt

- HTTPS Everywhere 를 추구하는 비영리 프로젝트
- IdenTrust cross-sign 됨
- SSL 인증서 100% 무료화
- 인증기관 연장 및 인증서 재발급 무료
- 콘솔상에서 인증서 발급/갱신/설치/자동화
- 멀티도메인 지원, SAN 기능(여러 도메인을 한 인증서로 묶어주는 기능) 지원

### SSL 인증서 인증 방식

**Webroot**

- `certbot`이 지정된 웹 서버의 루트 디렉토리에 인증 파일을 생성
- Let's Encrypt가 도메인에 HTTP 요청을 보내서 인증 파일을 확인
- 기존 웹 서버가 실행 중이어야 함
- 무중단 설정 가능

**Webserver**

- Certbot이 Apache/Nginx 서버 플러그인을 사용해 직접 인증 파일 관리, HTTP 요청 처리
- 인증 요청을 처리하는 동안 `certbot`이 웹 서버를 재구성하거나 중단할 수 있음
- 설정/리다이렉션 자동화

**Standalone**

- `certbot`이 자체적으로 HTTP 서버를 실행하여 인증 요청을 처리.
- 주로 웹서버 없이 테스트 용도로 사용
- 웹서버가 실행중이면 중단해야 함

**DNS**

- DNS에 TXT 레코드 추가로 소유권 확인
- 와일드카드 인증서 사용 가능
- Private 네트워크에서도 사용 가능

해당 과제에서는 `webserver` 방식 사용

### Bind

> `/etc/httpd/conf.d/<Domain>.conf`
> 
- 도메인용의 VirtualHost 파일
- VirtualHost 태그에서 Apache가 수신할 트래픽의 IP/Domain, Port 를 설정
- HTTP 설정만 있는 파일에 발급받은 SSL 인증서를 사용해 HTTPS 리다이렉션 설정

> `/etc/cron.d/certbot-renew`
> 
- cron을 사용해 인증서 발급 자동화

## 작업 과정

## 1. HTTP 실행

- 사용자 생성

```bash
useradd -m -G wheel sslhttps
echo "sslhttps" | passwd --stdin sslhttps
```

- apache 설치

```bash
sudo dnf update --exclude=kernel* -y
sudo dnf install httpd -y
sudo systemctl enable --now httpd
```

- 방화벽 설정

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

- httpd-vhosts.conf 설정

```bash
sudo vi /etc/httpd/conf.d/<Domain>.conf
```

```bash
<VirtualHost *:80>
    ServerName <Domain>
    DocumentRoot /var/www/html
</VirtualHost>
```

- http 동작 테스트

```bash
sudo systemctl restart httpd
```

## 2. SSL 실행

- certbot 설치

```bash
sudo dnf install epel-release -y
sudo dnf install certbot python3-certbot-apache -y
```

- SSL 인증서 생성
    - webserver 방식

```bash
sudo certbot --apache -d <Domain>
```

## 3. HTTPS 실행

- httpd-vhosts.conf 설정

```bash
sudo vi /etc/httpd/conf.d/<Domain>.conf
```

```bash
<VirtualHost *:80>
    ServerName <Domain>
    DocumentRoot /var/www/html

    # http를 https로 리다이렉션
    Redirect permanent / https://<Domain>/
</VirtualHost>

<VirtualHost *:443>
    ServerName <Domain>
    DocumentRoot /var/www/html

    # SSL 인증서 적용
    SSLEngine on
    SSLCertificateFile /etc/letsencrypt/live/<Domain>/fullchain.pem
    SSLCertificateKeyFile /etc/letsencrypt/live/<Domain>/privkey.pem

    <Directory /var/www/html>
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>
```

- Apache 서버 재시작

```bash
sudo apachectl configtest
sudo systemctl restart httpd
```

## 4. SSL 인증서 갱신 자동화

- 만료 기간 3개월
- crontab으로 자동 갱신 설정
- 만료일 30일 이내일 때 인증서 갱신 가능
- 갱신 테스트

```bash
sudo certbot renew --dry-run
```

- 갱신 스크립트 작성

```bash
sudo vi /etc/cron.d/certbot-renew

# 매일 새벽 3시 갱신 시도
0 3 * * * root /usr/bin/certbot renew --quiet --deploy-hook "systemctl reload httpd"
```

---

## 트러블슈팅

> Let’s Encrypt를 사용한 HTTP1.1 방식의 인증서 발급은 도메인에 80포트로 접근 가능해야 한다.
> 
- 공인 ip에 8080→80으로 포워딩해서 진행했었는데, ssl 인증서 발급에 오류가 발생
- 알아보니 Let’s Encrypt 인증 기관에서 HTTP1.1(webserver) 방식을 사용한 인증서 발급에는 80포트로 고정됨
    - 80 외에 다른 포트를 사용하려면 아래 세 방법 뿐
        - DNS TXT 레코드를 사용한 인증 사용
        - 이미 HTTPS가 적용되어 있다면 HTTPS를 사용하는 ALPN-01 인증 방식을 사용
        - 다른 인증 기관을 사용
- 이후에 포워딩을 80→80으로 수정해서 해결

## 참고

https://www.owl-dev.me/blog/42#
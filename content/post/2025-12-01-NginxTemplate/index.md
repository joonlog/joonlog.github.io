---
title: "Nginx 리버스 프록시 템플릿"
date: 2025-12-01T09:00:00+09:00
categories: ["Linux", "Middleware"]
tags: ["linux", "nginx", "ssl termination", "https redirect", "java nginx", "php nginx", "fastcgi"]
---


- Nginx 리버스 프록시 템플릿
    
    > HTTPS 리다이렉트 / SSL Termination / PHP-FPM / Java WAS 공통 패턴 정리
    > 
    - 아래 상황에 대한 Nginx 설정 정리
        - HTTP → HTTPS 리다이렉트
        - SSL termination + WAS 프록시
        - PHP-FPM 앞단의 Nginx (fastcgi)
        - Java Tomcat 앞단의 Nginx (proxy_pass)
        - LB를 통해 거쳐오는 실제 클라이언트 IP 식별
        - 로그 포맷 / gzip / timeout 등 공통 설정
    
    ## 1. Nginx 전체 구성 개요
    
    - HTTPS 리다이렉트나 SSL Termination 설정은 별도 `server`블럭으로 분리
        
        ```bash
        user nginx;
        worker_processes auto;
        error_log /var/log/nginx/error.log;
        pid /run/nginx.pid;
        
        include /usr/share/nginx/modules/*.conf;
        
        events {
            worker_connections 4096;
        }
        
        http {
            # 공통 설정 (X-Forwarded-For, 로그, gzip, timeout 등)
            # HTTPS 리다이렉트 + SSL 서버 블록
            # PHP-FPM 서버
            # Java Tomcat 서버 
        }
        ```
        
    
    ## 2. HTTPS 리다이렉트
    
    - 80 포트 HTTP로 들어온 요청을 443 포트 HTTPS로 리다이렉션
        
        ```bash
        server {
            listen 80;
            server_name <도메인>;
        
            # 모든 HTTP 요청을 HTTPS로 강제 이동
            return 301 https://$host$request_uri;
        }
        ```
        
    
    ## 3. HTTPS 서버 + SSL Termination + WAS 프록시
    
    - SSL Termination 수행
        - 클라이언트가 HTTPS로 접속 → 인증서 로드 후 TLS 핸드셰이크 → Nginx가 복호화된 HTTP를 WAS로 전달
        
        ```bash
        server {
            listen 443 ssl http2;
            server_name <도메인>;
        
            ssl_certificate     /etc/nginx/ssl/<crt 파일>;
            ssl_certificate_key /etc/nginx/ssl/<key 파일>;
        
            ssl_protocols TLSv1.2 TLSv1.3;
            ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256';
            ssl_prefer_server_ciphers on;
        
            location / {
                # WAS로는 HTTP로 프록시 (Termination 이후)
                proxy_pass http://<WAS 서버 IP>:<포트>;
        
                # LB를 통해 트래픽이 들어올 경우 LB 정보가 아닌 실제 클라이언트 정보를 수신하도록 설정
                proxy_set_header Host              $host;
                proxy_set_header X-Real-IP         $remote_addr;
                proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
                proxy_set_header X-Forwarded-Port  $server_port;
        
                proxy_http_version 1.1;
                proxy_set_header Connection "";
        
                # 버퍼/타임아웃 (앱에 맞게 조정)
                proxy_buffers         8 16k;
                proxy_buffer_size     32k;
                proxy_connect_timeout 300;
                proxy_send_timeout    600;
                proxy_read_timeout    600;
                send_timeout          600;
            }
        }
        ```
        
    
    ## 4. HTTP 블록 공통 설정
    
    - Real IP / 로그 / gzip / timeout / buffer 등
        - 모든 `server` 블록이 공유할 설정
        - `set_real_ip_from` / `real_ip_header`
            - LB 뒤에서 실제 클라이언트 IP 로그 복원, 안 하면 전부 LB IP로 찍힘
        - timeout / buffer는 서비스에 맞춰서 조정
        
        ```bash
        http {
            include /etc/nginx/mime.types;
            default_type application/octet-stream;
        
            # LB를 통해 트래픽이 들어올 경우 LB 정보가 아닌 실제 클라이언트 정보를 수신하도록 설정
            set_real_ip_from <LB IP 대역>/24;
            real_ip_header X-Forwarded-For;
            real_ip_recursive on;
        
            # 로그 포맷
            log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                            '$status $body_bytes_sent "$http_referer" '
                            '"$http_user_agent" "$http_x_forwarded_for" "$request_time"';
        
            access_log /var/log/nginx/access.log main;
        
            # gzip 압축
            gzip on;
            gzip_vary on;
            gzip_proxied any;
            gzip_comp_level 6;
            gzip_types text/plain text/css text/xml text/javascript
                       application/json application/javascript application/xml+rss
                       application/rss+xml font/truetype font/opentype
                       application/vnd.ms-fontobject image/svg+xml;
            gzip_min_length 1000;
        
            # 버퍼/타임아웃 (앱에 맞게 조정)
            keepalive_timeout 15;
            keepalive_requests 200;
        
            send_timeout        30s;
            client_body_timeout 30s;
            proxy_read_timeout  60s;
            proxy_connect_timeout 5s;
            proxy_send_timeout  60s;
        
            proxy_buffering  on;
            proxy_buffer_size 4k;
            proxy_buffers      8 4k;
        
        }
        ```
        
    
    ## 5. PHP-FPM 앞단 Nginx 패턴
    
    - `SCRIPT_FILENAME` 잘못 설정하면 빈 화면 404 에러 발생
    - `fastcgi_param`으로 IP/스킴 안 넘기면 로그인/도메인 등 의존 코드가 꼬임
    - 정적 리소스 캐싱 안걸어두면 모든 js/css 요청이 PHP로 들어가서 부하 증가
        
        ```bash
        server {
            listen 80;
            server_name <도메인>;
        
            root /home/<프로젝트>/public;
            index index.php index.html;
        
            # 보안 헤더
            add_header X-UA-Compatible "IE=Edge,chrome=1";
            add_header X-XSS-Protection "1; mode=block";
            add_header X-Frame-Options "SAMEORIGIN";
            add_header X-Content-Type-Options "nosniff";
        
            # 정적 리소스 캐싱
            location ~* \.(jpg|jpeg|png|gif|css|pdf|js|ico|svg|woff|woff2|ttf|eot|otf|map)$ {
                access_log off;
                log_not_found off;
                try_files $uri =404;
                expires 1M;
                add_header Cache-Control "public";
            }
        
            charset utf-8;
            error_page 404 /error/404;
        
            # PHP-FPM FastCGI 연동
            location / {
                include fastcgi_params;
        
                fastcgi_pass <WAS 서버 IP>:9000;
                fastcgi_index index.php;
        
                # PHP에서 실제 실행할 스크립트 경로
                fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        
                # 버퍼
                fastcgi_buffers 16 16k;
                fastcgi_buffer_size 32k;
        
                # X-Forwarded 설정
                fastcgi_param X-Real-IP         $remote_addr;
                fastcgi_param X-Forwarded-For   $proxy_add_x_forwarded_for;
                fastcgi_param X-Forwarded-Proto $scheme;
                fastcgi_param X-Forwarded-Host  $host;
                fastcgi_param X-Forwarded-Port  $server_port;
            }
        }
        ```
        
    
    ## 6. Java Tomcat 앞단 Nginx 패턴
    
    - Java WAS는 `proxy_pass`로 전달하는 구조라 PHP보다 단순하지만 X-Forwarded 헤더가 중요
        - Spring Security, OAuth redirect 등
    
    ### 6-1. Upstream 정의
    
    ```bash
    upstream was_backend {
        server <WAS 서버 IP>:8080;
        keepalive 200;
    }
    ```
    
    ### 6-2. server 블록
    
    - `X-Forwarded-Proto` 로 넘기지 않으면 Spring이 redirect URL을 http로 만들어서 https 접속인데 http로 튕기는 현상 발생
    - `Host`를 내부 IP로 넘기면 외부 링크/Redirect URI에 내부 IP가 노출됨
    - `/healthz` 같은 헬스체크 URL이 없으면 LB가 임의 URL(예: `/`) 헬스체크하면서 불필요 부하 발생
        
        ```bash
        server {
            listen 80;
            server_name <도메인>;
        
            # --- 헬스체크 엔드포인트 ---
            location = /healthz {
                add_header Content-Type text/plain;
                return 200 'ok';
            }
        
            # --- 기본 프록시 ---
            location / {
                proxy_pass http://was_backend;
        
                # 원본 정보 전달
                proxy_set_header Host              $host;
                proxy_set_header X-Real-IP         $remote_addr;
                proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
                proxy_set_header X-Forwarded-Port  $server_port;
        
                proxy_http_version 1.1;
                proxy_set_header Connection "";
        
                # 파일 업로드/다운로드가 크다면 필요에 따라 조정
                # client_max_body_size 50m;
            }
        }
        ```
---
title : Apache AJP Localhost 연결 에러
date : 2025-06-07 09:00:00 +09:00
categories : [Linux, Apache]
tags : [linux, apache, ajp, localhost, 127.0.0.1, ipv4/ipv6] #소문자만 가능
---

### 증상

- 웹 접속 시 503 오류 발생
- `mod_jk.log`에 AJP 연결 실패 로그 다수 출력
    
    ```bash
    [error] ajp_connect_to_endpoint::jk_ajp_common.c (127.0.0.1:8009) (errno=111)
    [error] jk_open_socket::jk_connect.c Failed socket to (127.0.0.1:8009)
    [error] ajp_send_request::jk_ajp_common.c (worker1) Connecting to backend failed
    ```
    
- `netstat` 기준으로 Tomcat은 AJP 포트(8009)를 정상적으로 Listen 중
    
    ```bash
    tcp  0  0 127.0.0.1:8009  ::1  LISTEN  [java 프로세스]
    ```
    

### workers.properties 설정

- 기존 설정
    
    ```bash
    worker.worker1.type=ajp13
    worker.worker1.host=localhost
    ```
    
    - host가 정상적으로 localhost로 되어 있는데, `mod_jk.log`에서는 ajp 연동 오류가 나오고, netstat은 IPv6에서만 LISTEN중
    - `/etc/hosts`에도 정상적으로 localhost와 127.0.0.1이 명시된 것으로 확인
        - localhost를 127.0.0.1로 변경
- 이후 설정
    
    ```bash
    worker.worker1.type=ajp13
    worker.worker1.host=127.0.0.1
    ```
    
    - host 값 변경 후 httpd 재기동 시 503 에러 해결

### 결론

- 로그를 종합했을 때, Tomcat은 `127.0.0.1:8009`에서 AJP를 수신하고 있었지만, `workers.properties`에 설정된 `localhost`가 내부적으로 `::1`(IPv6)로 해석되어 연결이 실패한 것으로 추정
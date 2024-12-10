---
title : Rocky8 Apache Tomcat 연동
date : 2024-12-09 18:08:00 +09:00
categories : [Linux, "Web Server"]
tags : [rocky8, apache, tomcat, web, was] #소문자만 가능
---

### 목표

-   Apache-Tomcat을 사용한 **WEB**-**WAS** 연동
-   **Proxy** / **AJP** 방식 둘 모두 수행
-   AJP는 **mod_jk** 사용

### 환경 설정

-   KT Cloud VM
    -   WEB 1대, WAS 1대
    -   공인 IP 포트포워딩
-   **OS**: Rocky Linux 8.1
-   WEB: Apache 2.4.6
-   WAS: Tomcat 9.0.97

## 개념 및 설정

### Proxy / AJP 비교

| **연동 방식** | Proxy | AJP (Apache JServ Protocol) |
| --- | --- | --- |
| **프로토콜 유형** | HTTP (텍스트 기반), HTTPS 가능 | AJP (바이너리 프로토콜) |
| **성능** | HTTP 오버헤드로 인해 상대적으로 느림 | 바이너리 프로토콜로 더 빠르고 효율적임 |
| **보안 설정** | HTTPS로 SSL/TLS 암호화 가능 | AJP 자체 암호화 미지원, `secret` 옵션으로 보안 강화 |
| **주 사용 목적** | 간단한 설정, 외부 요청 처리 | 내부 네트워크에서의 고성능 요청 처리 |
| **호환성** | Tomcat 외 서버와도 연동 가능 | Tomcat 전용 프로토콜 |


### Apache

> `/etc/httpd/conf/httpd.conf`

-   Apache 메인 설정 파일
-   `proxy-vhosts.conf`, `ajp-vhosts.conf`, `httpd-modjk.conf` 등 하위 파일 Include
-   `mod_jk` 와 같은 모듈 로드

> `/etc/httpd/conf/extra/proxy-vhosts.conf` `/etc/httpd/conf/extra/ajp-vhosts.conf`

-   Proxy / AJP 용으로 생성한 VirtualHost 파일
-   **VirtualHost** 태그에서 Apache가 수신할 트래픽의 IP/Domain, Port 를 설정
    -   여기서 도메인 / 경로 / 포트 기반 라우팅을 설정 이 실습에선 포트 기반 라우팅 사용
-   **ProxyPass**에서 트래픽을 전달할 WAS 경로, IP/Domain 등을 설정
-   AJP는 **JkMount**에서 WAS 경로, 워커를 설정

Apache-Tomcat 연동에 필요한 **proxy_module**, **proxy_http_module** 모듈은 기본 설치

AJP는 **mod_jk**, **mod_proxy_ajp** 두 모듈 중 하나만 있으면 연동 가능

| **모듈** | **mod_jk** | **mod_proxy_ajp** |
| --- | --- | --- |
| **설치** | 별도 설치 필요 | Apache 기본 모듈로 제공 |
| **설정** | 복잡 (별도 설정 파일 필요) | 간단 (Apache 설정 파일에서 구성) |
| **로드 밸런싱** | 강력 (다양한 알고리즘 지원) | 제한적 (별도 모듈 필요) |
| **클러스터링** | 세션 스티키니스 및 고급 기능 지원 | 제한적 |
| **성능** | 고성능 (대규모 환경에 적합) | 단순 연동 시 적합 |
| **사용 사례** | 대규모 클러스터, 고급 로드 밸런싱이 필요한 경우 | 간단한 AJP 연동 및 소규모 환경 |


> `/etc/httpd/conf/extra/httpd-modjk.conf\\`

-   mod_jk 모듈 설정 파일
-   <IfModule jk_module> 태그로 Apache에 mod_jk가 로드된 경우에만 적용
-   workers.properties 파일 경로 지정
-   로그 경로, 로그 레벨 지정

> `/etc/httpd/conf/extra/workers.properties`

-   mod_jk 워커(worker1 포함)의 IP, Port 등을 설정하는 파일
-   여러 워커 설정 가능

### Tomcat

> `/etc/systemd/system/tomcat.service`

-   rocky 8.1에서는 Tomcat이 systemd에 등록된 패키지가 없음
-   Tomcat을 systemd 서비스로 관리하기 위한 설정

> `/usr/local/apache-tomcat-9.0.97/conf/server.xml`

-   Proxy Connector은 기본 설정 사용
-   AJP Connector 포트, 허용 IP, Secret 설정
    -   이 실습에서 Secret은 false로 설정
        -   Secret을 사용 한다면 Apache에도 Secret 설정 필수

## 작업 과정

## 1. Proxy 방식

### 1) Apache 서버 설정

-   Apache 설치/실행

```bash
sudo dnf update --exclude=kernel*
sudo dnf install httpd -y
sudo systemctl enable --now httpd

```

-   모듈 활성화 확인 (proxy_module, proxy_http_module)

```bash
httpd -M | grep proxy

```

-   Apache 서버에서 HTTP 트래픽을 허용하도록 방화벽 설정 추가

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload

```

![Rocky8ApacheTomcat1.png](/assets/img/linux/Rocky8ApacheTomcat1.png)


-   proxy-vhosts.conf 설정 파일 추가

```bash
sudo vi /etc/httpd/conf/httpd.conf

```

```bash
# 별도의 파일에서 Virtual Host 설정 관리
Include conf/extra/proxy-vhosts.conf

```

-   Apache가 Tomcat 서버로 프록시 역할을 수행하도록 설정

```bash
sudo mkdir -p /etc/httpd/conf/extra
sudo vi /etc/httpd/conf/extra/proxy-vhosts.conf

```

```bash
# Apache가 8080 포트에서 수신된 요청 처리
<VirtualHost *:8080>
    # Apache 서버의 IP 또는 도메인 지정
    ServerName <WEB IP>
    
    # Apache가 수신한 요청을 WAS 8080 포트로 프록시 처리
    ProxyRequests Off
    ProxyPass / http://<WAS IP>:8080/
    ProxyPassReverse / http://<WAS IP>:8080/
</VirtualHost>

```

-   Apache 서버 재시작

```bash
sudo systemctl restart httpd

```

### 2) Tomcat 서버 설정

-   Java 설치

```bash
sudo dnf update --exclude=kernel*
sudo dnf install java-11-openjdk -y

java -version

```

-   ~/.bashrc 환경 변수 추가

```bash
export JAVA_HOME=$(dirname $(dirname $(readlink -f $(which java))))
export PATH=$PATH:$JAVA_HOME/bin

```

-   Tomcat 설치

```bash
wget <https://dlcdn.apache.org/tomcat/tomcat-9/v9.0.97/bin/apache-tomcat-9.0.97.tar.gz>
tar xvzf apache-tomcat-9.0.97.tar.gz apache-tomcat-9.0.97
sudo mv apache-tomcat-9.0.97 /usr/local

```

-   tomcat 사용자 설정

```bash
sudo groupadd tomcat
sudo useradd -g tomcat -d /usr/local/apache-tomcat-9.0.97 -s /bin/false tomcat
sudo chown -R tomcat:tomcat /usr/local/apache-tomcat-9.0.97

```

-   Tomcat 서비스 등록

```bash
vi /etc/systemd/system/tomcat.service

```

```bash
[Unit]
Description=tomcat 9
After=network.target syslog.target

[Service]
Type=forking
# JDK와 Tomcat 설치 경로 지정
Environment="JAVA_HOME=/usr/lib/jvm/java-11-openjdk-11.0.25.0.9-2.el8.x86_64"
Environment="CATALINA_HOME=/usr/local/apache-tomcat-9.0.97"
User=tomcat
Group=tomcat
ExecStart=/usr/local/apache-tomcat-9.0.97/bin/startup.sh
ExecStop=/usr/local/apache-tomcat-9.0.97/bin/shutdown.sh

[Install]
WantedBy=multi-user.target

```

-   Tomcat 실행

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now tomcat
sudo systemctl status tomcat

```

-   방화벽 HTTP 8080 포트 허용

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload

```

### 3) 설정 확인

-   브라우저에서 **Apache 서버의 IP 주소**로 접근하여 Tomcat의 기본 페이지가 나타나는지 확인

http://<공인 ip>:8080

![Rocky8ApacheTomcat2.png](/assets/img/linux/Rocky8ApacheTomcat2.png)


## 2. AJP 방식

### 1) Apache 서버 설정

-   Apache 설치/실행

```bash
sudo dnf update --exclude=kernel*
sudo dnf install httpd -y
sudo systemctl enable --now httpd

```

-   tomcat-connectors 설치(mod_jk)

```bash
sudo yum install autoconf libtool httpd-devel -y

wget <https://dlcdn.apache.org/tomcat/tomcat-connectors/jk/tomcat-connectors-1.2.50-src.tar.gz>
tar -zxvf tomcat-connectors-1.2.50-src.tar.gz
cd tomcat-connectors-1.2.50-src/native
./configure --with-apxs=/usr/bin/apxs

sudo dnf install make -y
sudo dnf install redhat-rpm-config -y

make 
sudo make install

```

-   모듈 활성화 확인

```bash
sudo systemctl restart httpd
httpd -M | grep jk

```

-   Apache 서버에서 HTTP 트래픽을 허용하도록 방화벽 설정 추가

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-port=8090/tcp
sudo firewall-cmd --reload

```

-   httpd.conf에 mod_jk 모듈 추가
-   httpd.conf에 httpd-modjk.conf 설정 파일 추가

```bash
sudo vi /etc/httpd/conf/httpd.conf

```

```bash
# 설치한 mod_jk 모듈 로드
LoadModule jk_module modules/mod_jk.so

# ajp 설정 파일 포함
Include conf/extra/httpd-modjk.conf

```

-   vhosts.conf 에 jkmount 추가

```bash
sudo mkdir -p /etc/httpd/conf/extra
sudo vi /etc/httpd/conf/extra/ajp-vhosts.conf

```

```bash
# Apache가 8090 포트에서 수신된 요청 처리
<VirtualHost *:8090>
    # AJP 요청을 수신할 Apache 서버 IP 지정
    ServerName <WEB IP>
    
    # 모든 요청(/*)을 worker1으로 라우팅
    JkMount /* worker1
</VirtualHost>

```

-   httpd-modjk.conf 파일 설정
    -   로그 레벨 debug로 변경 후 mod_jk.log에서 로그 확인 가능

```bash
sudo vi /etc/httpd/conf/extra/httpd-modjk.conf

```

```bash
<IfModule jk_module>
    # ajp worker 설정 파일 지정
    JkWorkersFile conf/extra/workers.properties
    # 상태 정보 파일 경로 지정
    JkShmFile logs/mod_jk.shm
    # 로그 파일 경로 지정
    JkLogFile logs/mod_jk.log
    # 로그 레벨 지정
    JkLogLevel debug
    JkLogStampFormat "[%a %b %d %H:%M:%S %Y]"
</IfModule>

```

-   workers.properties 파일 설정

```bash
sudo vi /etc/httpd/conf/extra/workers.properties

```

```bash
# 워커 목록
worker.list=worker1

# AJP 포트 
worker.worker1.port=8009
# Tomcat 서버 IP
worker.worker1.host=<WAS IP>
# 프로토콜 버전 지정
worker.worker1.type=ajp13

```

-   workers.properties, httpd-modjk.conf 권한 설정

```bash
sudo chmod 755 /etc/httpd/conf/extra/workers.properties
sudo chmod 755 /etc/httpd/conf/extra/httpd-modjk.conf

```

-   Apache 서버 재시작

```bash
sudo systemctl restart httpd

```

### 2) Tomcat 서버 설정

-   **Java 설치**

```bash
sudo dnf update --exclude=kernel*
sudo dnf install java-11-openjdk -y

java -version

```

-   **~/.bashrc 환경 변수 추가**

```bash
export JAVA_HOME=$(dirname $(dirname $(readlink -f $(which java))))
export PATH=$PATH:$JAVA_HOME/bin

```

-   **Tomcat 설치**

```bash
wget <https://dlcdn.apache.org/tomcat/tomcat-9/v9.0.97/bin/apache-tomcat-9.0.97.tar.gz>
tar xvzf apache-tomcat-9.0.97.tar.gz apache-tomcat-9.0.97
sudo mv apache-tomcat-9.0.97 /usr/local

```

-   tomcat 사용자 설정

```bash
sudo groupadd tomcat
sudo useradd -g tomcat -d /usr/local/apache-tomcat-9.0.97 -s /bin/false tomcat
sudo chown -R tomcat:tomcat /usr/local/apache-tomcat-9.0.97

```

-   Tomcat 서비스 등록

```bash
vi /etc/systemd/system/tomcat.service

```

```bash
[Unit]
Description=tomcat 9
After=network.target syslog.target

[Service]
Type=forking
# JDK와 Tomcat 설치 경로 지정
Environment="JAVA_HOME=/usr/lib/jvm/java-11-openjdk-11.0.25.0.9-2.el8.x86_64"
Environment="CATALINA_HOME=/usr/local/apache-tomcat-9.0.97"
User=tomcat
Group=tomcat
ExecStart=/usr/local/apache-tomcat-9.0.97/bin/startup.sh
ExecStop=/usr/local/apache-tomcat-9.0.97/bin/shutdown.sh

[Install]
WantedBy=multi-user.target

```

-   Tomcat AJP 커넥터 설정
-   Tomcat의 server.xml 파일을 열어 AJP 커넥터 활성화

```bash
sudo vi /usr/local/apache-tomcat-9.0.97/conf/server.xml

```

```bash
# AJP를 수신을 허용할 Source IP 지정
# AJP 키 인증 비활성화
<Connector port="8009" protocol="AJP/1.3" redirectPort="8443" address="<WAS IP>" secretRequired="false"/>

```

-   Tomcat 실행

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now tomcat
sudo systemctl status tomcat

```

-   방화벽 AJP 8080 포트 허용

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-port=8009/tcp
sudo firewall-cmd --reload

```

### 3) 설정 확인

-   브라우저에서 **Apache 서버의 IP 주소**로 접근하여 Tomcat의 기본 페이지가 나타나는지 확인

http://<공인ip>:8090

![Rocky8ApacheTomcat2.png](/assets/img/linux/Rocky8ApacheTomcat2.png)

----------

## 로그 확인

연동 시의 오류도 tomcat에 찍히기 때문에 catalina 로그를 가장 많이 확인

-   **Apache 로그 확인**

```bash
# 에러 로그
tail -f /var/log/httpd/error_log
# 액세스 로그
tail -f /var/log/httpd/access_log
# mod_jk 로그
tail -f /var/log/httpd/mod_jk.log

```

-   **Tomcat 로그 확인**

```bash
# tomcat 로그
tail -f /usr/local/apache-tomcat-9.0.97/logs/catalina.out
# localhost 로그
tail -f /usr/local/apache-tomcat-9.0.97/logs/localhost.log

```

## 트러블슈팅

> VirtualHost에서 / 경로는 webapps/ROOT를 가리킨다.

-   /ajp/* 와 같이 경로를 지정하면 apache-tomcat-9.0.97/webapps/ajp 폴더에 페이지가 존재해야 함

> Apache의 httpd.conf는 모듈이 앞에 있지 않으면 해당 모듈을 사용하는 설정은 무시된다.

-   httpd.conf에서 LoadModule이 Include보다 앞에 위치해야 함

> JkMount는 VirtualHost 내부에 위치해야 한다.
> 
> [*](https://ineastdev.tistory.com/entry/apache-tomcat-%EC%9D%84-modjk-%EB%A1%9C-%EC%97%B0%EB%8F%99%EC%8B%9C-%EC%A3%BC%EC%9D%98%EC%82%AC%ED%95%AD)[https://ineastdev.tistory.com/entry/apache-tomcat-을-modjk-로-연동시-주의사항*](https://ineastdev.tistory.com/entry/apache-tomcat-%EC%9D%84-modjk-%EB%A1%9C-%EC%97%B0%EB%8F%99%EC%8B%9C-%EC%A3%BC%EC%9D%98%EC%82%AC%ED%95%AD*)

-   VirtualHost를 설정한 적이 없다면 httpd.conf에 Include된 파일 중에 VirtualHost가 있는지 확인 필수
-   이 실습에서는 proxy 연동을 동시에 진행했기 때문에 proxy_vhosts.conf파일에 VirtualHost가 존재했고, AJP 연동에서 VirtualHost 없이 JkMount만 작성했기 때문에 오류가 발생

> Apache는 /ajp 를 입력받으면 내부적으로 /ajp/ 경로로 리다이렉션 한다.

-   경로 끝에 / 가 없는 URL을 입력할 경우, Apache에서는 `DirectorySlash Off` 설정으로 인해, 이 경로를 디렉토리로 간주하여 / 을 추가하여 리다이렉션 한 후에 내부에서 처리
-   만약 실제 디렉토리가 존재한다면 / 을 추가한 상태로 리다이렉션을 응답하고, 존재하지 않는다면 다른 리소스로 처리(기존 경로 응답)
-   `DirectorySlash Off` 설정으로 수정할 수 있지만, 대부분의 WEB의 기본값이기 때문에 수정을 권장하지 않음
-   따라서 JkMount의 경로에 /ajp 대신 /ajp/* 로 기입해 해결
-   만약 실제 디렉토리가 존재한다면 / 을 추가한 상태로 리다이렉션을 응답하고, 존재하지 않는다면 다른 리소스로 처리(기존 경로 응답)
-   `DirectorySlash Off` 설정으로 수정할 수 있지만, 대부분의 WEB의 기본값이기 때문에 수정을 권장하지 않음
-   따라서 JkMount의 경로에 /ajp 대신 /ajp/* 로 기입해 해결

----------

## 참고

**proxy 방식**:

[](https://velog.io/@ela__gin/%EB%A6%AC%EB%88%85%EC%8A%A4-%EC%84%9C%EB%B2%84%EC%97%90%EC%84%9C-ApacheTomcat-mod-proxy-%EB%B0%A9%EC%8B%9D-%EC%97%B0%EB%8F%99-%EC%9B%B9%EC%86%8C%EC%BC%93-%ED%91%B8%EC%8B%9C-%EC%95%8C%EB%A6%BC-%EC%97%B0%EB%8F%99-%ED%8F%AC%ED%95%A8)[https://velog.io/@ela__gin/리눅스-서버에서-ApacheTomcat-mod-proxy-방식-연동-웹소켓-푸시-알림-연동-포함](https://velog.io/@ela__gin/%EB%A6%AC%EB%88%85%EC%8A%A4-%EC%84%9C%EB%B2%84%EC%97%90%EC%84%9C-ApacheTomcat-mod-proxy-%EB%B0%A9%EC%8B%9D-%EC%97%B0%EB%8F%99-%EC%9B%B9%EC%86%8C%EC%BC%93-%ED%91%B8%EC%8B%9C-%EC%95%8C%EB%A6%BC-%EC%97%B0%EB%8F%99-%ED%8F%AC%ED%95%A8)

**ajp 방식**:

[https://nahosung.tistory.com/121](https://nahosung.tistory.com/121)

[https://lilo.tistory.com/83](https://lilo.tistory.com/83)

[](https://ineastdev.tistory.com/entry/apache-tomcat-%EC%9D%84-modjk-%EB%A1%9C-%EC%97%B0%EB%8F%99%EC%8B%9C-%EC%A3%BC%EC%9D%98%EC%82%AC%ED%95%AD)[https://ineastdev.tistory.com/entry/apache-tomcat-을-modjk-로-연동시-주의사항](https://ineastdev.tistory.com/entry/apache-tomcat-%EC%9D%84-modjk-%EB%A1%9C-%EC%97%B0%EB%8F%99%EC%8B%9C-%EC%A3%BC%EC%9D%98%EC%82%AC%ED%95%AD)

**tomcat 서비스 등록**: [https://getchu.tistory.com/14](https://getchu.tistory.com/14)
---
title: "Apache 버전 업그레이드"
date: 2025-11-27T09:00:00+09:00
categories: ["Linux", "Middleware"]
tags: ["linux", "apache", "apache version upgrade"]
---


> 기존 Apache가 설치된 환경에서 버전 업그레이드
> 

> Apache 2.4.63 신규 버전을 소스 컴파일로 설치하고 심볼릭 링크 기반 전환으로 서비스 다운타임 최소화
> 

```bash
httpd -t: 설정 문법을 테스트
httpd -S: 가상호스트 등 서버 설정을 확인
httpd -V: 컴파일 옵션과 버전을 조회
```

기존 apache 경로: /etc/httpd

신규 apache 경로: /usr/local/apache

### 0. 고려사항

- 기존 Apache 패키지 경로: `/etc/httpd`
    - 새 버전은 `/usr/local/apache` 아래에 설치해 충돌 방지
- 심볼릭 링크 기반 버전 관리
    
    `/usr/sbin/httpd` → `/usr/local/apache/bin/httpd` 로 링크하여 서비스 파일 수정 최소화
    
- conf, modules, workers.properties 등 핵심 설정만 이관
- logs 경로는 기존 경로를 그대로 사용
    
    `/usr/local/apache/logs` → `/etc/httpd/logs` 로 링크
    
- 업그레이드 중 필수 검증:
    - `httpd -t` 설정 문법 체크
    - `httpd -S` 가상호스트 및 include 파일 존재 여부 확인
    - 없는 파일은 `cp -arup` 로 보완

### **1. 신규 버전 다운로드 및 배포 준비**

- /usr/local에 신규 파일 배치
    
    ```bash
    cd /usr/local
    tar -zxvf apache.2.4.63.tar.gz
    mv apache apache-2.4.63
    ln -sfn /usr/local/apache-2.4.63 /usr/local/apache
    ```
    

### **2. 기존 설정 및 모듈 이관**

- `conf` / `conf.d` / `conf.modules.d` 파일 복사
    - 신규 버전에 기본 제공되지 않는 파일(proxy-html.conf 등)은 `httpd -S`에서 오류가 나는지 확인하고 불필요하면 삭제
    
    ```bash
    # httpd.conf
    cp -arup /etc/httpd/conf/httpd.conf /usr/local/apache/conf/
    
    # conf.d → extra/ 로 이관
    cp -arup /etc/httpd/conf.d/* /usr/local/apache/conf/extra/
    
    # conf.modules.d
    cp -arup /etc/httpd/conf.modules.d/ /usr/local/apache/conf.modules.d
    ```
    
- 모듈 복사
    
    ```bash
    cp -arup /etc/httpd/modules/mod_jk.so /usr/local/apache/modules/
    cp -arup /etc/httpd/modules/mod_ssl.so /usr/local/apache/modules/
    cp -arup /etc/httpd/modules/mod_rpaf.so /usr/local/apache/modules/
    cp -arup /etc/httpd/modules/mod_proxy_http2.so /usr/local/apache/modules/
    cp -arup /etc/httpd/modules/mod_brotli.so /usr/local/apache/modules/
    cp -arup /etc/httpd/modules/mod_suexec.so /usr/local/apache/modules/
    cp -arup /etc/httpd/modules/mod_systemd.so /usr/local/apache/modules/
    ```
    
- workers.properties
    
    ```bash
    cp -arup /etc/httpd/conf/workers.properties /usr/local/apache/conf/
    ```
    

### **3. 로그 경로 유지**

- 기존 경로 재사용
    
    ```bash
    rm -rf /usr/local/apache/logs
    ln -s /etc/httpd/logs /usr/local/apache/logs
    ```
    

### 4. httpd.conf 설정 변경

- `ServerRoot` 경로 수정 및 모듈 적용
    
    ```bash
    vim /usr/local/apache/conf/httpd.conf
    
    ServerRoot "/usr/local/apache"
    
    ServerSignature Off
    ServerTokens Prod
    
    LoadModule jk_module    modules/mod_jk.so
    LoadModule rpaf_module  modules/mod_rpaf.so
    
    Include conf.modules.d/*.conf
    ```
    

### 5. 설정 검증

- Missing file 오류가 나오면 해당 파일을 `/etc/httpd/` 기준으로 찾아 복사
- syntax 오류가 나면 conf 파일 수정 후 재검증
    
    ```bash
    /usr/local/apache/bin/httpd -t     # 설정 문법 검증
    /usr/local/apache/bin/httpd -S     # 가상호스트 및 include 파일 확인
    ```
    

### **6. systemd 서비스 연결**

- 기존 서비스 파일(`/usr/lib/systemd/system/httpd.service`)은 수정하지 않고, 바이너리만 신규 httpd 로 연결
    
    ```bash
    service httpd stop
    
    cd /usr/sbin
    mv httpd httpd.bak
    ln -s /usr/local/apache/bin/httpd httpd
    
    service httpd start
    ```
    

### **7. 확인 필요 사항**

- `httpd -t` 정상
- `httpd -S` 에 missing 파일 없음
- mod_jk, mod_ssl, mod_rpaf 등 필수 모듈 로딩 정상
- 서비스 기동 후
    - `/etc/httpd/logs/error_log`
    - 애플리케이션 로그
    - SSL 인증서 로드 정상 여부
    - 가상호스트 매핑 정상 작동 여부
- 원복 필요 시 `/usr/sbin/httpd.bak` 로 즉시 롤백
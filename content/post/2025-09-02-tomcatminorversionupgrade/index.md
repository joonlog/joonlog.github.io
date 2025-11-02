---
title: "Tomcat 마이너 버전 업그레이드 (Tomcat 9)"
date: 2025-09-02T09:00:00+09:00
categories: ["Linux", "Middleware"]
tags: ["linux", "tomcat", "minor version upgrade", "tomcat9"]
---


- Tomcat 9.0 버전대에서 9.0.108 버전으로의 마이너 버전 업그레이드
    - 마이너 버전 업그레이드는 전/후 버전의 설정이 대부분 호환됨

### 0. 고려사항

- 심볼릭 링크 방식으로 버전을 관리해서 롤백을 용이하게 설정
- `conf`, `lib`, `bin` 일부 스크립트만 선별적으로 이관
- 로그/웹앱 경로는 기존 디렉토리를 링크해 데이터 손실을 방지
- 업그레이드 후에는 반드시 `catalina.out`과 애플리케이션 로그를 확인해 정상 구동 여부를 검증

### 1. 새 버전 다운로드 및 설치

```bash
cd /usr/local
curl -O https://dlcdn.apache.org/tomcat/tomcat-9/v9.0.108/bin/apache-tomcat-9.0.108.tar.gz
tar xzf apache-tomcat-9.0.108.tar.gz
mv apache-tomcat-9.0.108 tomcat-9.0.108
```

### 2. 설정 및 배포물 이관

- 업그레이드 시 필요한 설정/라이브러리만 이관
    - conf
        - 웬만하면 전부 복사
    - bin
        - 이전 버전에서 catalina.sh 파일이 수정된 흔적이 있을 경우 반드시 파악해서 setenv.sh에 이관
            - setenv.sh로 관리하면 추후 작업시에도 이 파일만 옮기면 됨
        - setenv.sh 파일 존재할 시 복사
    - lib
        - 이전 버전에서 기본 라이브러리 외에 추가된 라이브러리가 있다면 복사
    
    ```bash
    rsync -a /usr/local/tomcat/conf/  /usr/local/tomcat-9.0.108/conf/
    rsync -a /usr/local/tomcat/bin/setenv.sh /usr/local/tomcat-9.0.108/bin/
    cp /usr/local/tomcat/lib/tibero6-jdbc.jar /usr/local/tomcat-9.0.108/lib/
    ```
    
- 기존 버전에서 log, webapps 경로가 심볼릭 링크로 되어 있다면 신규 버전에서도 동일하게 설정
    
    ```bash
    rmdir /usr/local/tomcat-9.0.108/logs
    ln -sfn /data/tomcat/logs /usr/local/tomcat-9.0.108/
    ```
    

### 3. 서비스 전환

- 구버전 종료 후 새 버전으로 심볼릭 링크 전환
    
    ```bash
    su - tomcat
    /usr/local/tomcat/bin/shutdown.sh
    
    mv /usr/local/tomcat /usr/local/tomcat-9.0.102
    ln -sfn /usr/local/tomcat-9.0.108 /usr/local/tomcat
    
    /usr/local/tomcat/bin/startup.sh
    ```
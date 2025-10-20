---
title : Docker Tomcat에서 Context Path 설정 무시 TroubleShooting
date : 2025-10-04 09:00:00 +09:00
categories : [Container, Java]
tags : [docker, container java 1.8, tomcat 9, context path, servlet path, appbase, docbase] #소문자만 가능
---

> IDE 환경에서는 잘 작동하던 Context Path 설정이 Docker Tomcat에서는 무시되어 발생한 문제
Tomcat의 WAR 파일 배포 메커니즘과 appBase/docBase 개념을 이해해야 해결 가능한 이슈
> 

### 환경

- Java 1.8
- Tomcat 9.0.41
- Spring 기반 CMS
- Docker로 Tomcat WAS 구성

### 문제 상황

1. 빌드한 app.war를 Docker Tomcat으로 배포 후 `/app/index.do` 접근 시 무한 리다이렉트 발생
    - `/app/index.do` → `/index.html` → `/app/index.do` 반복
2. context path가 `/`로 설정됐어야 하는데 `app`으로 설정되어 있음
    - server.xml 설정이 무시됨
- 초기 설정
    - Dockerfile
        
        ```docker
        COPY app.war ${CATALINA_HOME}/webapps/app.war
        COPY server.xml ${CATALINA_HOME}/conf/server.xml
        ```
        
    - server.xml
        
        ```xml
        <Host name="localhost" appBase="webapps"
              unpackWARs="true" autoDeploy="true">
            <Context path="" docBase="app" reloadable="false" />
        </Host>
        ```
        

### 에러 분석

1. 무한 리다이렉트 발생 이유
    - 애플리케이션에 index.do를 받을 수 있는 2개의 컨트롤러 존재
        1. `@RequestMapping("/index.do")`
            - index.html로 리다이렉트
        2. `@RequestMapping("/{siteId}/index.do")`
            - 실제 사이트 로직 처리
    - 브라우저에서 `/app/index.do`로 접근 시 2번으로 가야하는데 계속 1번으로만 처리됨
        - context path가 `/`가 아닌 `app`으로 되어 있다는 증거
            1. `/app/index.do`로 요청
            2. Tomcat이 context path 제거 → Servlet Path = `/index.do`
            3. Spring이 `/index.do` 매핑 찾음 → 1번 컨트롤러 매칭
            4. `redirect:/index.html` 반환
            5. `/index.html` 내부에 `location.href = "/app/index.do"` 스크립트
            6. 1번으로 돌아가서 무한 루프
2. server.xml 설정이 무시된 이유
    - Tomcat의 WAR 배포 우선순위
        1. `autoDeploy="true"` 설정 시 webapps 디렉토리 자동 스캔
        2. `app.war` 발견 → 파일명 기준으로 자동 배포
            - Context Path = `/app`
        3. server.xml의 `<Context>` 설정은 자동 배포에 의해 무시됨
3. IDE 환경과의 차이가 발생한 이유
    - IDE
        1. UI에서 Context Path를 `/`로 직접 지정 가능
        2. WAR 파일명과 무관하게 IDE 내부적으로 Context 설정 파일 생성
    - Docker Tomcat
        1. WAR 파일명이 Context Path를 결정
        2. Tomcat 자체 규칙만 적용

### 해결 방법 1

- war 파일을 ROOT.war로 배포
    - 가장 간단하고 확실한 방법
    - Context Path가 ROOT.war 파일명에 따라 `/`로 설정됨
    - Dockerfile
        
        ```docker
        COPY app.war ${CATALINA_HOME}/webapps/ROOT.war
        ```
        
    1. /app/index.do로 요청
    2. Tomcat이 context path 제거 → Servlet Path = `/app/index.do`
    3. Spring이 `/{siteId}/index.do` 패턴 매칭 → 2번 컨트롤러 매칭 (siteId=app)
    4. 정상 동작

### 해결 방법 2

- Context xml 파일 설정
    - app.war 파일명을 유지하면서 Context Path를 변경하는 방법
    - appBase: Tomcat이 애플리케이션을 찾는 작업 공간
        - 기본: `webapps`
    - docBase: 특정 애플리케이션의 실제 파일 위치
    - Tomcat 9 보안 정책: appBase 내부의 WAR를 docBase로 지정하면 무시됨
        - war 파일을 appBase와 다른곳으로 이동
        - Dockerfile
            
            ```docker
            # WAR를 appBase 밖에 배치
            RUN mkdir -p ${CATALINA_HOME}/apps
            COPY app.war ${CATALINA_HOME}/apps/app.war
            
            # Context 설정 파일 생성
            RUN mkdir -p ${CATALINA_HOME}/conf/Catalina/localhost
            COPY ROOT.xml ${CATALINA_HOME}/conf/Catalina/localhost/ROOT.xml
            
            # server.xml 수정 (autoDeploy 비활성화)
            COPY server.xml ${CATALINA_HOME}/conf/server.xml
            ```
            
        - ROOT.xml
            
            ```xml
            <?xml version="1.0" encoding="UTF-8"?>
            <Context docBase="${catalina.home}/apps/app.war" />
            ```
            
        - server.xml
            
            ```xml
            <?xml version="1.0" encoding="UTF-8"?>
            <Server port="8005" shutdown="SHUTDOWN">
              <Service name="Catalina">
                <Connector port="8080" protocol="HTTP/1.1"
                           connectionTimeout="20000"
                           redirectPort="8443"
                           URIEncoding="UTF-8" />
            
                <Engine name="Catalina" defaultHost="localhost">
                  <Host name="localhost" appBase="webapps"
                        unpackWARs="true" autoDeploy="false" deployOnStartup="true">
                  </Host>
                </Engine>
              </Service>
            </Server>
            ```
            
    - 동작 순서
        1. `server.xml`의 `autoDeploy="false"` 설정으로 webapps의 WAR자동 배포 비활성화
        2. `server.xml`의 `deployOnStartup="true"` 설정으로 Context XML 파일 읽기 활성화
        3. `ROOT.xml` → Context Path `/` 지정
            - 파일명이 ROOT이므로 `/`로 설정됨
        4. `ROOT.xml`의 `docBase="/opt/tomcat/apps/app.war"` 설정으로 appBase 밖의 경로에서 war참조
        5. Tomcat이 `apps/app.war`를 `webapps/ROOT/`로 압축 해제 후 실행
        6. 정상 동작

### 정리

1. Tomcat의 WAR 배포 규칙
    - WAR 파일명: Context Path
        
        ```bash
        ROOT.war -> /
        admin.war -> /admin
        my-app.war -> /my-app
        ```
        
2. server.xml의 Context 설정은 비권장
    - `autoDeploy="true"` 환경에서는 대부분 무시됨
    - 설정 변경 시 Tomcat 전체 재시작 필요
3. appBase / docBase
    
    > appBase 내부의 WAR를 docBase로 지정하면 Tomcat 9에서 무시됨
    > 
    
    | 항목 | appBase | docBase |
    | --- | --- | --- |
    | 정의 | Tomcat의 작업 공간 | 앱의 소스 위치 |
    | 기본값 | `webapps` | - |
    | 내용 | WAR + 압축 해제된 코드 | WAR 또는 디렉토리 경로 |
    | 자동 배포 | O | X |
4. Context Path / Servlet Path
    
    > 애플리케이션 라우팅 로직이 Servlet Path를 기준으로 설계되었다면 Context Path가 라우팅에 영향을 줄 수 있음
    > 
    - 요청 URL: /app/index.do
        
        ```bash
        Context Path: /
        Servlet Path: /app/index.do
        ```
        
    - 요청 URL: /app/index.do
        
        ```bash
        Context Path: /app
        Servlet Path: /index.do
        ```
        

### 참고

- https://tomcat.apache.org/tomcat-9.0-doc/config/context.html
- https://tomcat.apache.org/tomcat-9.0-doc/config/host.html#Automatic_Application_Deployment
- https://docs.spring.io/spring-framework/reference/web/webmvc/mvc-controller/ann-requestmapping.html
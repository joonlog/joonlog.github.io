---
title: "죽은 Connection 사용 중인 Tomcat Troubleshooting"
date: 2025-06-19T09:00:00+09:00
categories: ["Linux", "Middleware"]
tags: ["linux", "tomcat", "testwhileidle", "timeBetweenEvictionRunsMillis", "troubleshooting"]
---


### 문제 상황

- 웹 서비스 접근 시 500 에러 출력

### 원인 분석

- 웹 서비스에 DB 서버 2대를 이중화하여 구성 중, 2번 서버로 DB Failover된 후 장애 발생
- application log 확인
    
    ```bash
    org.springframework.jdbc.support.SQLErrorCodesFactory] Error while extracting database product name - falling back to empty error codes
    org.springframework.jdbc.support.SQLStateSQLExceptionTranslator] Extracted SQL state class 'JD' from value 'JDBC-90405:ERRJD'
    org.springframework.transaction.interceptor.TransactionInterceptor] Application exception overridden by rollback exception
    org.apache.catalina.core.StandardWrapperValve.invoke Servlet.service() for servlet [...] threw exception [...] with root cause
    ```
    
- Tomcat에서 Failover 이전 DB 서버로의 Connection을 계속 사용해서 `SQLException` 오류 발생하여 500 에러 출력

### 조치 내용

> 즉각적인 웹 서비스 정상화를 위한 임시 조치
> 
1. Tomcat 재기동 후 웹 서비스 정상 접근 확인

### 추가 조치 방안

- Tomcat의 server.xml 혹은 `context.xml`에서 `testWhileIdle` 설정
    - `testWhileIdl`: Tomcat이 Connection을 반환하기 전, 주기적으로 유휴 Connection 생존 여부를 검사
    - `timeBetweenEvictionRunsMillis`: 30초마다 검사

```bash
testWhileIdle="true"
timeBetweenEvictionRunsMillis="30000"
```

### 결론

Tomcat의 JDBC Connection Pool은 `testWhileIdle` 옵션이 `false`일 경우, 장애 이후에도 죽은 Connection을 계속 사용할 수 있어 장애가 지속됨. 이 설정을 `true`로 변경하고 `timeBetweenEvictionRunsMillis`를 설정해서 DB Failover 시 자동적으로 Connection을 복구하도록 설정

### 참고

- https://blog.naver.com/wasgosu-2010/90189771881
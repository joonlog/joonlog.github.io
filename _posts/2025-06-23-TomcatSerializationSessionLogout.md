---
title : 직렬화 기능을 사용중인 Tomcat의 세션 만료로 인한 강제 로그아웃 Troubleshooting
date : 2025-06-23 09:00:00 +09:00
categories : [Linux, Middleware]
tags : [linux, tomcat, serialization, session] #소문자만 가능
---

### 문제 상황

- 웹 서비스 사용 중 세션이 갑자기 만료되며 강제 로그아웃되는 현상 발생

---

### 원인 분석

- Catalina 로그에서 `ClassNotFoundException` 확인
    
    ```
    SEVERE [Catalina-utility-1] org.apache.catalina.session.StoreBase.processExpires
    Error processing session expiration for key [BA8BFEFE...JVM1]
    java.lang.ClassNotFoundException: egovframework.com.utl.slm.EgovHttpSessionBindingListener
    	at org.apache.catalina.session.StandardSession.doReadObject(StandardSession.java:1268)
    	at org.apache.catalina.session.StandardSession.readObjectData(StandardSession.java:846)
    	at org.apache.catalina.session.FileStore.load(FileStore.java:203)
    	at org.apache.catalina.session.StoreBase.processExpires(StoreBase.java:138)
    ```
    
- Tomcat의 세션 직렬화 기능으로 인해, 서버가 재기동되거나 일정 시간이 경과한 후 세션 정보가 디스크에 저장된 파일(`SESSIONS.ser`)을 통해 복원 중
- `SESSIONS.ser`에는 `EgovHttpSessionBindingListener`라는 클래스가 포함되어 있었으나, 현재 애플리케이션에는 이 클래스가 존재하지 않아 역직렬화 과정에서 오류 발생한 것으로 추정
- Tomcat이 해당 세션을 복원하지 못해 세션을 강제 만료 처리 → 사용자 입장에서는 강제 로그아웃으로 인식됨

---

### 조치 내용

> 세션 복원 실패에 따른 반복적인 오류를 차단하기 위해 다음과 같은 조치 수행
> 
1. Tomcat 세션 파일 SESSIONS.ser 삭제
2. Tomcat 재기동

→ 이후 세션 관련 오류 재발 없음, 사용자 세션 정상 유지 확인

---

### 결론

- Tomcat은 세션 정보를 디스크에 저장 후 복원하는 과정에서, 클래스 누락 등으로 인해 세션 역직렬화에 실패할 수 있음
- 이 경우 해당 세션이 강제 만료되는 문제 발생
- 이후에도 세션 직렬화 기능을 사용한다면, 애플리케이션 배포 전 ser 파일 삭제 권장
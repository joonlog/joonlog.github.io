---
title: "HikariCP connection is closed 에러 TroubleShooting"
date: 2025-09-10T09:00:00+09:00
categories: ["Linux", "DB"]
tags: ["linux", "spring", "hikaricp", "mysql", "wait_timeout", "maxlifetime", "troubleshooting"]
---


- MySQL은 `wait_timeout` 설정값에 따라 커넥션 중 일정 시간이상 사용하지 않은 커넥션을 **종료**
- HikariCP는 `maxLifetime` 설정값에 따라 스스로 미사용된 커넥션을 제거하고 새로 생성하는 방식으로 동작
    - 커넥션 누수 발생 시 connection closed 되지 않으므로 Hikari는 커넥션을 `activate` 상태로 파악하여 `maxLifetime` 에 따른 커넥션 갱신하지 않음
        - 따라서 `wait_timeout` 으로 인해 닫힌 커넥션을 재사용하려다 `connection is closed` 오류가 발생
        - HikariCP는 `maxLifetime`값을 DB의 `wait_timeout`보다 **몇 초 정도** 짧게 설정하라고 권장

### connection is closed 에러 히스토리

- connection pool 수치 점검
    - default 10개, 필요시 20~30개로 설정
    
    ```bash
    spring.datasource.hikari.maximum-pool-size=30
    ```
    
    - 현재 실행 중인 쿼리의 Que 값이 connection Pool 최대 값까지 도달하여 다음 실행할 쿼리를 처리하지 못해 발생
    - APM 솔루션으로 Connection Pool 자원 모니터링을 통해 증설이 필요할 시 즉시 증설, 쿼리 최적화 필요 언급
        - 이건 임시 조치만 반복할 뿐, 커넥션 누수는 코드를 수정해야 함
- db 파라미터 `wait_timeout = 500` 설정
- 이후 재발
    
    ```bash
    spring.datasource.hikari.minimum-idle=10 // 최소 커넥션 수
    spring.datasource.hikari.idle-timeout=300000 // 유휴 커넥션 타임아웃 (5분)
    spring.datasource.hikari.max-lifetime=1800000 // 최대 커넥션 생명 주기 (30분)
    spring.datasource.hikari.connection-timeout=30000 // 커넥션 대기 시간 (30초)
    
    # db 파라미터 max_connection = 2000으로 확인됨
    spring.datasource.hikari.maximum-pool-size=60 // 30에서 60으로 변경
    ```
    
- 이후 재발
    - HikariCP는 `maxLifetime`값을 DB의 `wait_timeout`보다 **몇 초 정도** 짧게 설정하라고 권장
        
        ```bash
        spring.datasource.hikari.max-lifetime=497000 // 최대 커넥션 생명 주기 (497초)
        ```
        
- 이후 재발
    - was 서버 로그 확인 시 maxLifetime 값을 더 짧게 설정하는 것을 권장하는 로그 확인
        
        ```bash
        [WARN ] [PoolBase.javalisConnectionAlive(184) : HikariPool-1 - Failed to connection org.mariadb.jdbc.MariaDbConnection@14441ec9 ((conn=8360) Connection.setNetworkTimeout cannot be called on a closed connection). Possibly consider using a shorter maxLifetime value.
        ```
        
    - maxLifetime 값 변경
        
        ```bash
        spring.datasource.hikari.max-lifetime=470000 // 최대 커넥션 생명 주기 (470초)
        ```
        
- 동일오류 재발
    - spring.datasource.hikari.connection-init-sql 프로퍼티를 추가하여 세션 wait_timeout을 설정
    - wait_timeout 값을 설정된 hikari MaxLifetime 값 + 5초 형식으로 설정
        
        ```bash
        spring.datasource.hikari.connection-init-sql=SET wait_timeout = 500
        ```
        

### 커넥션 누수 사례

https://jaehoney.tistory.com/337

- wait_timeout 이 15초로 설정되어 있는데 레거시 코드 때문에 높이지 못함
- 따라서 세션의 wait_timeout을 설정
    
    ```bash
    spring.datasource.hikari.connection-init-sql=SET wait_timeout = 500
    ```
    

https://do-study.tistory.com/97

1. EntityManager를 Bean으로 설정함
    - EntityManager는 스레드간 절대 공유하면 안됨
2. QueryDSL 구현 Resitory에서 특정 datasource의 EntityManager를 받기 위해 setEntityManager를 재정의하지 않고 별도 setting 메소드 작성
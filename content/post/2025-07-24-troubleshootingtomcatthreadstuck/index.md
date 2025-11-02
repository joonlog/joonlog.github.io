---
title: "Tomcat Thread Stuck으로 인한 사이트 접근 불가 TroubleShooting"
date: 2025-07-24T09:00:00+09:00
categories: ["Linux", "Middleware"]
tags: ["linux", "tomcat", "tibero", "tibero backup", "thread stuck", "troubleshooting"]
---


### 문제 개요

- 최근 특정 새벽 시간에 웹 서비스 접근 불가 현상이 반복 발생
- Tomcat 로그 확인 시 Thread Stuck으로 요청 처리가 지연되었으며, 일정 시간 지난 후 자동 해소

```bash
WARNING argo.server.valves.LenaStuckThreadDetectionValve.notifyStuckThreadDetected ...
at com.tmax.tibero.jdbc.common.TbStream.readMsgDs(Unknown Source)
```

### 원인 추적

> 왜 Thread Stuck이 발생했는가?
> 
- Tibero DB를 사용 중인 DB 서버의 문제 발생 시간대에 sys.log에서 대량의 임시 테이블이 생성된 것을 확인
- 매일 01시에 `tbexport full=Y` 명령어를 수행하는 cronjob 설정된 상태
- Full 백업 시 전체 용량이 몇백 GB가 넘는데, 백업 간 DB 리소스 고갈로 인해 Tomcat의 스레드가 응답 받지 못해 Thread Stuck이 발생한 것으로 확인

### 조치

- 백업 방식을 아래와 같이 논리적(tbexport) → 물리적 백업으로 전환
    - 대용량의 DB 백업 시에는 물리적 백업이 모든 면에서 유리
        - 무중단 Hot Backup 가능
            - 단, Hot Backup은 `archive log` 모드에서만 가능
        - 내부 자원 사용 최소화
        - 복구 및 운영 안정성 향상

### 물리적 백업 방안

1. tbrmgr
- 사용해보니 매일 풀백업을 사용하는 현재의 환경에선 맞지 않고, 백업한 파일을 압축하지 않고 그대로 보관하면서 증분 백업을 통해 보관하는 환경에서 맞는 방안인 것으로 확인
    - 현재는 매일 풀백업 + 풀백업한 파일 압축하여 보관하는 환경에서 사용 가능한 방법 필요
2. begin/end 백업
- 백업 모드 시작 - 데이터파일 복사 - 백업 모드 종료의 구조로 백업하는 방식
- 이 방식으로 해결
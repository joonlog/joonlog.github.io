---
title: "MySQL / MariaDB 덤프 및 복원"
date: 2025-05-16T09:00:00+09:00
categories: ["Linux", "DB"]
tags: ["linux", "db", "mysql", "mariadb", "mysql dump", "mysql restore"]
---


- 서비스 운영 중에 DB 백업/복원 작업할 때 주의할 점 정리
- 옵션에 따라 백업/복원 기능이 다르니 주의 필요

---

## 1. DB 접근

```bash
mysql -u 사용자 -p -h 도메인명
```

---

## 2. DB 덤프 (백업)

```bash
mysqldump -h 도메인명 -u 사용자 -p \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  BASE_DB > DB명_날짜.sql
```

### 옵션 설명

- `-single-transaction`: InnoDB 전용. 덤프 중 트랜잭션 시작해서 일관된 스냅샷 유지. 서비스 중에도 락 없이 가능.
- `-routines`: 저장 프로시저, 함수 포함.
- `-triggers`: 각 테이블의 트리거 포함.
- `-events`: 이벤트 스케줄러 이벤트 포함.

---

## 3. 덤프 파일 확인 (USE 문 유무)

- `mysqldump` 시 `--databases`  옵션 사용하면 `.sql` 파일에 `CREATE DATABASE` + `USE DB명;` 구문 포함됨
- 복원 시 다른 DB에 덮어씌어짐

> **반드시 확인 필요**
> 

```bash
head -n 30 DB명_날짜.sql | grep -i "^USE"
```

- 만약 `USE DB명;` 이 있을 경우 → 주석 처리

```sql
-- Before: USE BASE_DB;
-- After : -- USE BASE_DB;
```

---

## 4. 덤프 파일을 통한 복원

```bash
mysql -h 도메인명 -u 사용자 -p 생성한DB명 < DB명_날짜.sql
```

- 신규 DB는 복원 전 DB 생성 필요

---

## 5. 기존/신규 DB 확인

- 기존/신규 DB 내 테이블의 rows 수가 동일한지 확인

```bash
mysql> SELECT 기존DB명 AS source, COUNT(*) FROM 기존DB명.PTNT_RSRV_HIST_BACK
    -> UNION ALL
    -> SELECT 신규DB명, COUNT(*) FROM 신규DB명.PTNT_RSRV_HIST_BACK;
```

## 주의사항 요약

- `mysqldump`에 `-databases` 옵션 쓰면 `USE` 문 포함됨 → 다른 DB 덮어씌우는 원인.
- 클라우드 DB는 `SUPER` 권한 없을 수 있음 → `-set-gtid-purged=OFF` 넣어야 에러 안남.
- 복원할 DB명과 `.sql` 파일 내 `USE` 구문 다를 경우 조심해야 함.
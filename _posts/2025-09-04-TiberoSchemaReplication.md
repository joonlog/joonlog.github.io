---
title : 같은 Tibero DB 안에서 스키마 복제하는 방법
date : 2025-09-04 09:00:00 +09:00
categories : [Linux, DB]
tags : [linux, db, tibero6, tibero, schema replication, tbexport] #소문자만 가능
---

> Tibero DB에서 스키마를 복제할 경우 기존 DB와의 테이블스페이스 및 데이터파일 중복을 고려해서 작업 필요
> 
1. 복제 원본 스키마 정보 확인
    - 대상 스키마에서 사용하는 테이블스페이스, 데이터파일, 크기 확인
    
    ```bash
    -- 원본 스키마가 쓰는 테이블스페이스 목록
    SELECT DISTINCT tablespace_name
    FROM dba_segments
    WHERE owner = <대상 스키마명>
    ORDER BY 1;
    
    -- 각 TS에 대해 현재 데이터파일 경로와 크기
    SELECT tablespace_name,
    file_id,
    file_name,
    ROUND(bytes/1024/1024) AS size_mb,
    autoextensible,
    ROUND(maxbytes/1024/1024) AS max_mb
    FROM dba_data_files
    WHERE tablespace_name IN (
    SELECT DISTINCT tablespace_name FROM dba_segments WHERE owner=<대상 스키마명>
    )
    ORDER BY tablespace_name, file_id;
    
    -- 각 TS에서 대상 스키마가 실제 사용 중인 용량(세그먼트 기준)
    SELECT tablespace_name,
    ROUND(SUM(bytes)/1024/1024) AS used_mb
    FROM dba_segments
    WHERE owner=<대상 스키마명>
    GROUP BY tablespace_name
    ORDER BY tablespace_name;
    ```
    
2. 대상(복제) 스키마 구성
    - 대상 스키마 생성, 권한 부여, 쿼터 부여
    
    ```bash
    -- 대상 테이블 스페이스 생성
    	CREATE TABLESPACE <대상TS1> DATAFILE '<경로>/<대상TS1>.dbf' SIZE <사이즈>M <여부 예: AUTOEXTEND ON>;
    	CREATE TABLESPACE <대상TS2> DATAFILE '<경로>/<대상TS2>.dbf' SIZE <사이즈>M <여부 예: AUTOEXTEND ON>;
    -- 복제 적용할 대상 전부 동일하게 구성(원본 스키마와 테이블 스페이스를 공유할 경우 적용X)
    
    -- 대상 스키마 생성
    CREATE USER <대상 스키마명> IDENTIFIED BY "<패스워드>"
      DEFAULT TABLESPACE <대상TS, 예: <복제 스키마>>
      TEMPORARY TABLESPACE TEMP;
    
    -- 권한 부여(필요시)
    GRANT CONNECT, RESOURCE TO <대상 스키마명>;
    
    -- 대상 스키마에 각 대상 TS 쿼터 부여
    ALTER USER <대상 스키마명> QUOTA UNLIMITED ON <대상TS>;
    -- 복제 적용할 대상 전부 구성
    ```
    
3. tbexport 백업 파일 생성
    - 테이블스페이스 리맵 적용
    
    ```bash
    
    tbexport \
      username=<계정명> password=<비밀번호> sid=<DB명> \
      file=<덤프 파일 경로 예: /data/db_dump/src_full_remap.dmp> \
      user=<원본 스키마명> \
      remap_tablespace=<원본TS1>:<대상TS1>,<원본TS2>:<대상TS2>
      -- 원본 스키마와 테이블 스페이스를 분리할 테이블 스페이스만 적용
    ```
    
4. tbimport 데이터 적재
    
    ```bash
    tbimport \
      username=<계정명> password=<비밀번호> sid=<DB명> \
      file=<덤프 파일 경로 예: /data/db_dump/src_full_remap.dmp> \
      fromuser=<원본 스키마명> touser=<대상 스키마명> \
      rows=y constraint=n index=n trigger=n ignore=Y
    ```
    
5. tbimport 메타데이터 적재
    
    ```bash
    tbimport \
      username=<계정명> password=<비밀번호> sid=<DB명> \
      file=<DUMP_FILE> \
      log=<LOG_META> \
      fromuser=<원본 스키마명> touser=<대상 스키마명> \
      rows=n constraint=y index=y trigger=y ignore=Y
    ```
    
6. 검증
    
    ```bash
    -- 테이블 수 비교
    SELECT COUNT(*) FROM dba_tables WHERE owner=<대상 스키마>;
    SELECT COUNT(*) FROM dba_tables WHERE owner=<복제 스키마>;
    
    -- 인덱스 수 비교
    SELECT COUNT(*) FROM dba_indexes WHERE owner=<대상 스키마>;
    SELECT COUNT(*) FROM dba_indexes WHERE owner=<복제 스키마>;
    ```
    
7. 외부 스키마의 권한 부여(선택)
    - 복제 스키마 외에 다른 스키마가 존재하고, 해당 스키마가 복제한 스키마를 참조할 경우 권한을 별도로 부여해야 함
    
    ```bash
    1) 오브젝트 권한(테이블/뷰/시퀀스/프로시저/펑션/패키지 등)
    -- <대상 스키마>에는 있고 <복제 스키마>에는 없는 권한
    SELECT privilege, owner, table_name
    FROM dba_tab_privs
    WHERE grantee=<대상 스키마>
    MINUS
    SELECT privilege, owner, table_name
    FROM dba_tab_privs
    WHERE grantee=<복제 스키마>;
    
    2) 컬럼 단위 권한(있을 수 있음: SELECT/UPDATE 컬럼 권한)
    -- <대상 스키마>에는 있고 <복제 스키마>에는 없는 컬럼 권한
    SELECT privilege, owner, table_name, column_name
    FROM dba_col_privs
    WHERE grantee=<대상 스키마>
    MINUS
    SELECT privilege, owner, table_name, column_name
    FROM dba_col_privs
    WHERE grantee=<복제 스키마>;
    
    3) 시스템 권한(필요시)
    -- <대상 스키마>에는 있고 <복제 스키마>에는 없는 시스템 권한
    SELECT privilege
    FROM dba_sys_privs
    WHERE grantee=<대상 스키마>
    MINUS
    SELECT privilege
    FROM dba_sys_privs
    WHERE grantee=<복제 스키마>;
    
    4) 롤(역할) 부여 비교(필요시)
    -- <대상 스키마>에는 있고 <복제 스키마>에는 없는 롤
    SELECT granted_role
    FROM dba_role_privs
    WHERE grantee=<대상 스키마>
    MINUS
    SELECT granted_role
    FROM dba_role_privs
    WHERE grantee=<복제 스키마>;
    ```
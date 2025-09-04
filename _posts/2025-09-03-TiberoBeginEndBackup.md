---
title : Tibero DB에서 begin/end모드를 사용한 물리적 백업 및 복구 방법
date : 2025-09-03 09:00:00 +09:00
categories : [Linux, DB]
tags : [Linux, db, tibero6, tibero, begin-end backup, pigz] #소문자만 가능
---

- 백업 모드 시작 - 데이터파일 복사 - 백업 모드 종료의 구조로 백업하는 방식
- 백업 스크립트
    - DB내 존재하는 테이블스페이스 리스트를 읽어 반복문으로 풀 백업
    - pigz 명령어를 사용한 멀티스레드 압축
        - 다수의 대용량 데이터 파일을 압축해야하는 현재 상황에서는 pigz를 사용하지 않고 tar로만 압축하면 24시간 이상 소요 ⇒ pigz 사용 시 6시간 소요
    - 백업 파일과 DB 메타데이터 비교를 통한 누락 파일 감지
    - 이전 백업 파일 및 로그 정리

```bash
#!/bin/bash

set -euo pipefail

########## 환경설정 ##########
export TB_HOME=/tibero/tibero6
export TB_SID=<sid명>
export PATH="$TB_HOME/client/bin:/usr/bin:/bin"
export LD_LIBRARY_PATH="$TB_HOME/client/lib:${LD_LIBRARY_PATH:-}"
TBSQL="$TB_HOME/client/bin/tbsql"

USER=<계정명>
PASS=<비밀번호>

DATE=$(date +%Y%m%d)
BACKUP_ROOT="/data/db_dump"
BACKUP_DIR="$BACKUP_ROOT/$DATE"
DATA_DIR="/data/<db명>"
ARCHIVE_DIR="/data/archive"
LOGFILE="$BACKUP_ROOT/<db명>_backup_$DATE.log"
LOG_DIR="$BACKUP_ROOT"
LOG_FILE="$BACKUP_ROOT/<db명>_backup_check_$DATE.log"
PIGZ="/usr/bin/pigz"
export LC_ALL=C

mkdir -p "$BACKUP_DIR"

########## 유틸 ##########
log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOGFILE" ; }
tbq() { echo "$1" | "$TBSQL" "$USER/$PASS" >> "$LOGFILE" 2>&1 ; }

########## 1) DB에서 TS,파일 경로 Spool ##########
MAP_FILE="$BACKUP_DIR/_filemap_${DATE}.csv"        # 원본 Spool
MAP_CLEAN="$BACKUP_DIR/_filemap_${DATE}.clean.csv" # 정제본
log "===== DBA_DATA_FILES Spool 시작 ====="
"$TBSQL" "$USER/$PASS" <<EOF >> "$LOGFILE" 2>&1
SET ECHO OFF
SET PAGESIZE 0
SET FEEDBACK OFF
SET HEADING OFF
SET TRIMSPOOL ON
SPOOL $MAP_FILE
SELECT TRIM(tablespace_name) || ',' || TRIM(file_name)
FROM dba_data_files
ORDER BY tablespace_name, file_name;
SPOOL OFF
EXIT;
EOF

# 정제: CR 제거 → 앞의 "SQL> [번호]" 제거 → 공백 Trim → 패턴 필터링
sed -E 's/\r$//' "$MAP_FILE" \
| sed -E 's/^SQL>[[:space:]]*([0-9[:space:]]+)?//' \
| awk -F, 'NF==2 {
    gsub(/^[[:space:]]+|[[:space:]]+$/,"",$1);
    gsub(/^[[:space:]]+|[[:space:]]+$/,"",$2);
    if ($1 ~ /^[A-Z0-9_]+$/ && $2 ~ /^\//) print $1","$2
}' > "$MAP_CLEAN"

if [[ ! -s "$MAP_CLEAN" ]]; then
  log "ERROR: 파일 맵 정제 결과가 비었습니다. ($MAP_CLEAN)"
  exit 1
fi
log "Spool 원본 라인수: $(wc -l < "$MAP_FILE"), 정제 후 라인수: $(wc -l < "$MAP_CLEAN")"
log "===== Spool 완료(정제 적용): $(wc -l < "$MAP_CLEAN") rows ====="

########## 2) TS별 파일 목록 구성 ##########
declare -A TS_FILES
declare -a TS_ORDER
while IFS=, read -r TS F; do
  [[ -z "$TS" || -z "$F" ]] && continue
  [[ "$TS" =~ ^[A-Z0-9_]+$ ]] || continue
  [[ "$F"  =~ ^/.*$        ]] || continue
  if [[ -z "${TS_FILES[$TS]+_}" ]]; then
    TS_FILES[$TS]="$F"
    TS_ORDER+=("$TS")
  else
    TS_FILES[$TS]="${TS_FILES[$TS]} $F"
  fi
done < "$MAP_CLEAN"

log "감지된 테이블스페이스 수: ${#TS_ORDER[@]}"

########## 3) TS 하나 백업(복사/압축) ##########
backup_one_ts() {
  local TS="$1"; shift
  local FILES=("$@")
  local BASENAMES=()
  local MANIFEST="${TS}${DATE}.manifest"
  local TGZ="${TS}_${DATE}.tar.gz"

  log "===== ${TS} BEGIN BACKUP ====="
  tbq "ALTER TABLESPACE ${TS} BEGIN BACKUP;"

  for f in "${FILES[@]}"; do
    if [[ -r "$f" ]]; then
      cp "$f" "$BACKUP_DIR/" >> "$LOGFILE" 2>&1
      BASENAMES+=("$(basename "$f")")
    else
      log "WARN: ${TS} 데이터파일 읽기 불가: $f"
    fi
  done

  tbq "ALTER TABLESPACE ${TS} END BACKUP;"

  if (( ${#BASENAMES[@]} == 0 )); then
    log "WARN: ${TS} 복사된 파일이 없어 압축을 생략합니다."
    return 0
  fi

  (
    cd "$BACKUP_DIR" || exit 1
    tar --index-file="$MANIFEST" -cvf - "${BASENAMES[@]}" 2>>"$LOGFILE" \
      | "$PIGZ" -p 8 > "$TGZ"
    rm -f "${BASENAMES[@]}" 2>>"$LOGFILE" || true
  ) >> "$LOGFILE" 2>&1 &
}

########## 4) 병렬 폭 조절(PARALLEL=2: 두 개씩) ##########
PARALLEL=2
running=0
for TS in "${TS_ORDER[@]}"; do
  read -r -a arr <<< "${TS_FILES[$TS]}"
  backup_one_ts "$TS" "${arr[@]}"
  ((++running))
  if (( running >= PARALLEL )); then
    wait
    running=0
    log "===== 직전 ${PARALLEL}개 TS 백업/압축 완료(wait) ====="
  fi
done
wait
log "===== 모든 TS 백업/압축 완료 ====="

########## 5) MISC 묶음 ##########
MISC_DIR="$BACKUP_DIR/.misc_tmp"
mkdir -p "$MISC_DIR"

log "===== 컨트롤파일 SQL 백업 ====="
tbq "ALTER DATABASE BACKUP CONTROLFILE TO TRACE AS '${MISC_DIR}/crectl.sql' REUSE NORESETLOGS;"

log "===== 컨트롤파일 물리 백업 ====="
cp /data/<db명>/c1.ctl "$MISC_DIR/" 2>> "$LOGFILE" || true
cp /tibero/tibero6/config/c2.ctl "$MISC_DIR/" 2>> "$LOGFILE" || true

log "===== 아카이브 로그 스위치 및 복사 ====="
tbq "ALTER SYSTEM SWITCH LOGFILE;"
cp "$ARCHIVE_DIR"/*.arc "$MISC_DIR/" 2>> "$LOGFILE" || true

log "===== 리두로그/.passwd 백업 ====="
cp /data/<db명>/log*.log "$MISC_DIR/" 2>> "$LOGFILE" || true
cp /data/<db명>/.passwd  "$MISC_DIR/" 2>> "$LOGFILE" || true

tar -czf "${BACKUP_DIR}/MISC_${DATE}.tar.gz" -C "$MISC_DIR" . >> "$LOGFILE" 2>&1
rm -rf "$MISC_DIR"

log "===== 전체 백업 완료(MISC 포함) ====="

########## 6) LABEL 생성 및 검증(오탐 방지) ##########
: > "$BACKUP_DIR/LABEL${DATE}.manifest"
cat "$BACKUP_DIR"/*"${DATE}.manifest" 2>/dev/null \
  | awk -F/ '{print $NF}' \
  | sort -u > "$BACKUP_DIR/LABEL${DATE}.manifest"

if [[ ! -s "$BACKUP_DIR/LABEL${DATE}.manifest" ]]; then
  log "ERROR: LABEL${DATE}.manifest is empty (no fallback by policy)"
  exit 1
fi

: > "$LOG_DIR/backup_filelist.txt"
awk -F/ '{print $NF}' "$BACKUP_DIR/LABEL${DATE}.manifest" | sort -u > "$LOG_DIR/backup_filelist.txt"

# expect 목록 재생성 (tbsql 출력은 메인 로그만)
$TBSQL "$USER/$PASS" <<EOF >> "$LOGFILE" 2>&1
SET ECHO OFF
SET PAGESIZE 0
SET FEEDBACK OFF
SET HEADING OFF
SET TRIMSPOOL ON
SPOOL $LOG_DIR/expected_filelist_full.txt
SELECT file_name FROM dba_data_files ORDER BY 1;
SPOOL OFF
EXIT;
EOF

awk -F/ '{print $NF}' "$LOG_DIR/expected_filelist_full.txt" | sort > "$LOG_DIR/expected_filelist.txt"

# 체크 로그 초기화 후 누락만 기록
: > "$LOG_FILE"
comm -23 "$LOG_DIR/expected_filelist.txt" "$LOG_DIR/backup_filelist.txt" > "$LOG_FILE"

if [[ -s "$LOG_FILE" ]]; then
  log "누락된 데이터파일이 있습니다. ($LOG_FILE 참조)"
else
  log "DB expect 목록 대비 누락 없음"
fi

rm -f "$LOG_DIR/backup_filelist.txt" \
      "$LOG_DIR/expected_filelist.txt" \
      "$LOG_DIR/expected_filelist_full.txt"

########## 7) 하우스키핑 ##########
log "Deleting .arc archive logs older than 7 days..."
find "$ARCHIVE_DIR" -type f -name "*.arc" -mtime +7 -print -exec rm -f {} \; >> "$LOGFILE" 2>&1

log "Deleting .log files older than 1 day..."
find "$LOG_DIR" -maxdepth 1 -type f -name "*.log" -mtime +1 -print -exec rm -f {} \; >> "$LOGFILE" 2>&1

YESTERDAY=$(date -d "1 day ago" +%Y%m%d)
if [[ -d "$BACKUP_ROOT/$YESTERDAY" ]]; then
  log "Deleting backup directory: $BACKUP_ROOT/$YESTERDAY"
  rm -rf "$BACKUP_ROOT/$YESTERDAY" >> "$LOGFILE" 2>&1
else
  log "No backup directory to delete for $YESTERDAY"
fi

```

- 백업 복구 방법

```bash
### 파일 이름은 백업본 생성 일자에 따라 변경 필요

# DB 중지
tbdown immediate

# 압축 해제
tar -xzf /data/db_dump/20250806.tar.gz -C /data/db_dump

# 파일 복구(운영 경로로 이동)
cp /data/db_dump/<날짜>/*.arc /data/archive/
cp /data/db_dump/<날짜>/*.dtf /data/<SID명>/
cp /data/db_dump/<날짜>/*.dbf /data/<SID명>/
cp /data/db_dump/<날짜>/.passwd /data/<SID명>/
cp /data/db_dump/<날짜>/*.log /data/<SID명>/

#인스턴스 NOMOUNT
tbboot nomount

#컨트롤파일 재생성
@/data/db_dump/<날짜>/crectl.sql

#MOUNT로 변경
tbdown
tbboot mount

#복구
ALTER DATABASE RECOVER AUTOMATIC DATABASE; 

#오픈
ALTER DATABASE OPEN RESETLOGS;
```
---
title : Tibero DB에서 tbrmgr 명령어를 사용한 물리적 백업 및 복구 방법
date : 2025-09-03 09:00:00 +09:00
categories : [Linux, DB]
tags : [linux, db, tibero6, tibero, tbrmgr] #소문자만 가능
---

- 백업한 데이터파일을 압축 및 변형하지 않고, 지속적으로 증분하여 백업할 경우에 유용한 방식
- 백업 스크립트
    
    ```bash
    #!/bin/bash
    
    # 오늘 날짜
    today=$(date +%Y%m%d)
    
    # 백업 루트 디렉토리
    backup_dir="/data/db_dump"
    
    # 오늘 백업 디렉토리 및 압축 파일 경로
    today_dir="$backup_dir/$today"
    today_tar="$backup_dir/<db명>_backup_${today}.tar.gz"
    
    # 아카이브 로그 디렉토리
    archive_log_dir="/data/archive"
    
    # 로그 기록용 파일
    log_file="$backup_dir/<sid명>_backup_${today}.log"
    
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] === Tibero Hot Backup Script Start ===" | tee -a "$log_file"
    
    # 1. 오늘 날짜 디렉토리 생성
    mkdir -p "$today_dir"
    
    # 2. 백업 실행
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting backup to $today_dir" | tee -a "$log_file"
    tbrmgr backup --compress=HIGH -u -o "$today_dir"
    
    du -sh "$today_dir" | tee -a "$log_file"
    
    if [ $? -eq 0 ]; then
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup succeeded. Compressing..." | tee -a "$log_file"
        tar -czf "$today_tar" -C "$backup_dir" "$today"
        if [ $? -eq 0 ]; then
            echo "[$(date '+%Y-%m-%d %H:%M:%S')] Compression successful: $today_tar" | tee -a "$log_file"
            rm -rf "$today_dir"
        else
            echo "[$(date '+%Y-%m-%d %H:%M:%S')] Compression failed!" | tee -a "$log_file"
        fi
    else
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup failed!" | tee -a "$log_file"
        exit 1
    fi
    
    # 3. 1일 초과된 백업 파일 삭제 (.tar.gz)
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Deleting .tar.gz backups older than 1 days..." | tee -a "$log_file"
    find "$backup_dir" -maxdepth 1 -type f -name "*.tar.gz" -mtime +1 -print -exec rm -f {} \; >> "$log_file" 2>&1
    
    du -sh "$today_tar" | tee -a "$log_file"
    
    # 4. 7일 초과된 아카이브 로그(.arc) 삭제
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Deleting .arc archive logs older than 7 days from $archive_log_dir..." | tee -a "$log_file"
    find "$archive_log_dir" -type f -name "*.arc" -mtime +7 -print -exec rm -f {} \; >> "$log_file" 2>&1
    
    # 5. 1일 초과된 백업 로그 파일 삭제 (.log)
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Deleting .log files older than 1 day..." | tee -a "$log_file"
    find "$backup_dir" -maxdepth 1 -type f -name "*.log" -mtime +1 -print -exec rm -f {} \; >> "$log_file" 2>&1
    
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] === Tibero Hot Backup Script Complete ===" | tee -a "$log_file"
    ### 백업스크립트 완료
    ```
    
- 복구 스크립트
    
    ```bash
    압축 해제
    tar -xvzf /data/db_dump/<db명>_backup_<날짜>.tar.gz -C /data/db_dump/
    
    복구 실행
    tbrmgr recover -o /data/db_dump/<날짜>/
    ```
---
title : 리눅스 CMD 로그, 시스템 로그 logrotate, openfile, sar 로그 설정
date : 2025-10-30 09:00:00 +09:00
categories : [Linux, System]
tags : [linux, system log, cmd log, openfile, sar, logrotate]  #소문자만 가능
---

> Redhat 계열 리눅스에서 시스템 로그 및 openfile, sar 등 기본 설정하는 방법
> 

### 환경

- OS: Rocky 8.10

### 개요

- 명령어 실행 기록은 history로도 볼수있지만 각 계정마다 저장 위치가 다르고 없어질 수도 있음
- 명령어 추적을 위해 CMD 로그를 활성화
- 로그 파일 용량 문제를 방지하기 위해 logrotate 설정
- 운영 중 파일이 동시에 열릴 수 있는 파일 디스크럽터 값 증설
    - open file limit

## 설정

### 01. CMD 로그 설정

- /etc/profile.d/cmd_logging.sh
    - cmd 로그 생성 스크립트
    
    ```bash
    function logging
    {
        stat="$?"
        cmd=$(history|tail -1)
        remoteaddr=$(who am i)
        if [ "$cmd" != "$cmd_old" ]; then
            logger -p local1.notice "[2] STAT=$stat"
            logger -p local1.notice "[1] $USER=$remoteaddr, PID=$$, PWD=$PWD, CMD=$cmd"
        fi
        cmd_old=$cmd
    }
    trap logging DEBUG
    ```
    
- 권한 설정
    
    ```bash
    chmod 700 /etc/profile.d/cmd_logging.sh
    source /etc/profile.d/cmd_logging.sh
    ```
    
- rsyslog 설정 변경
    
    ```bash
    cp -aurp /etc/rsyslog.conf /home/ncloud24/rsyslog.conf.$(date +%Y%m%d)
    ```
    
- /etc/rsyslog.conf
    - 파일 하단에 추가
    
    ```bash
    local1.notice                                               /var/log/cmd.log
    ```
    
- rsyslog 재기동
    
    ```bash
    systemctl restart rsyslog
    ```
    

### 02. Logrotate 설정

```bash
cp /etc/logrotate.conf /etc/logrotate.conf.bak
```

- vim /etc/logrotate.conf
    
    ```bash
    daily                # 매일 회전
    rotate 365           # 365일(1년) 보관
    create               # 회전 후 새 파일 생성
    dateext              # 파일명에 날짜 확장자 사용
    compress             # 압축 활성화
    compresscmd /bin/gzip
    compressext .gz
    include /etc/logrotate.d
    ```
    
- 로그 회전 대상에 cmd.log 추가
- vim /etc/logrotate.d/syslog
    
    ```bash
    /var/log/cron
    /var/log/maillog
    /var/log/messages
    /var/log/secure
    /var/log/spooler
    /var/log/cmd.log
    {
        missingok
        sharedscripts
        postrotate
            /usr/bin/systemctl -s HUP kill rsyslog.service >/dev/null 2>&1 || true
        endscript
    }
    ```
    

### 03. 로그 관리 설정

- 로그 백업 디렉토리 추가
    
    ```bash
    mkdir -p /var/log/backup
    ```
    
- /home/scripts/logMoveDelete.sh
    - 로그 관리 스크립트
    
    ```bash
    #!/bin/bash
    PATH=/usr/sbin:/usr/bin:/sbin:/bin
    
    source_dir="/var/log"
    backup_dir="/var/log/backup"
    
    mv "$source_dir"/*.tar.gz "$backup_dir" 2>/dev/null || true
    find "$backup_dir" -name "*.tar.gz" -type f -mtime +365 -exec rm {} \;
    ```
    
    - 권한 설정
    
    ```bash
    chmod 700 /home/scripts/logMoveDelete.sh
    ```
    
- cron 등록
    
    ```bash
    crontab -e
    
    0 1 * * * /home/scripts/logMoveDelete.sh
    ```
    

### 04. open file limit 설정

- 로그 및 서비스에서 파일 디스크립터 제한이 발생하지 않도록 상한값을 조정
    
    ```bash
    echo '* soft nofile 65535' | sudo tee -a /etc/security/limits.conf
    echo '* hard nofile 65535' | sudo tee -a /etc/security/limits.conf
    ```
    
    - /etc/systemd/user.conf
        
        ```bash
        DefaultLimitNOFILE=65535
        ```
        

### 05. sar log 설정

- sysstat 설정 변경
- /etc/sysconfig/sysstat
    
    ```bash
    HISTORY=365         # 데이터 보존 기간 (일 단위)
    COMPRESSAFTER=90    # 90일 경과 시 gzip 압축
    ZIP="gzip"          # 압축 방식
    ```
    
- /etc/cron.d/sysstat
    
    ```bash
    */1 * * * * root /usr/lib64/sa/sa1 1 1     # 1분마다 sa1 실행
    53 23 * * * root /usr/lib64/sa/sa2 -A      # 매일 23:53에 sa2 실행 (일간 요약)
    ```
    
- sysstat 서비스 활성화    
    
    ```bash
    systemctl enable --now sysstat
    ```
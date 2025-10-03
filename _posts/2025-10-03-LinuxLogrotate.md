---
title : Linux logrotate를 통한 로그 관리
date : 2025-10-03 09:00:00 +09:00
categories : [Linux, System]
tags : [linux, logrotate, mariadb]  #소문자만 가능
---

- /etc/cron.daily/logrotate에 설정된 스크립트를 통해 자동으로 동작
- 예시
    - mariadb `general.log` 대상

```bash
vim /etc/logrotate.d/mariadb

/data/mariadb_log/general.log {

  # Depends on a mysql@localhost unix_socket authenticated user with RELOAD privilege
  su mysql mysql

  # If any of the files listed above is missing, skip them silently without
  # emitting any errors
  missingok

  # If file exists but is empty, don't rotate it
  notifempty

  # 달에 1번 보관
  monthly

  # 6개월 보관
  rotate 6

  # 500M 초과할 경우 달에 상관없이 즉시 보관
  maxsize 500M

  # 50M 미만일 경우 달에 상관없이 보관 X
  minsize 50M

  # 압축해서 보관
  compress

  # 가장 최근 보관한 로그는 압축하지 않음
  delaycompress

  # Don't run the postrotate script for each file configured in this file, but
  # run it only once if one or more files were rotated
  sharedscripts

  # redhat 계열 flush
    postrotate
        if command -v mariadb-admin >/dev/null 2>&1; then
            mariadb-admin --user=root --password='gC6mlwuspkkw!@' --host=localhost \
              --local flush-error-log flush-engine-log flush-general-log flush-slow-log
        fi
    endscript
}
```
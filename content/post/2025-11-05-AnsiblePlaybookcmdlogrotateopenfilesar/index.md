---
title: "Ansible Playbook으로 CMD 로그 / logrotate / openfile / SAR 설정하기"
date: 2025-11-05T09:00:00+09:00
categories: ["Ansilbe", "Playbook"]
tags: ["ansible", "playbook", "cmd log", "logrotate", "openfile", "sar"]
---


> Redhat 계열 리눅스에서 Ansible로 시스템 로그 및 openfile, sar 설정 자동화
> 

### 환경

- Rocky 8.10

## Ansible 설정

### 01. CMD 로그 설정

- cmdlog_rogrotate_openfile.yaml
    
    ```bash
    ---
    - name: Configure cmd logging, logrotate, and open files on Rocky 8.10
      hosts: prod
      become: true
    
      handlers:
        - name: restart rsyslog
          ansible.builtin.systemd:
            name: rsyslog
            state: restarted
            enabled: true
    
      tasks:
        - name: 01. Ensure directories exist
          ansible.builtin.file:
            path: "{{ item }}"
            state: directory
            mode: '0755'
          loop:
            - /home/ncloud24/scripts
            - /var/log/backup
    
        - name: 02. Deploy logMoveDelete.sh
          ansible.builtin.copy:
            src: files/logMoveDelete.sh
            dest: /home/ncloud24/scripts/logMoveDelete.sh
            owner: root
            group: root
            mode: '0700'
    
        - name: 03. Register daily cron for logMoveDelete.sh (01:00)
          ansible.builtin.cron:
            name: "Move and cleanup .tar.gz logs"
            user: root
            minute: "0"
            hour: "1"
            job: "/home/ncloud24/scripts/logMoveDelete.sh"
    
        - name: 04. Deploy /etc/logrotate.conf (backup original)
          ansible.builtin.copy:
            src: files/logrotate.conf
            dest: /etc/logrotate.conf
            owner: root
            group: root
            mode: '0644'
            backup: true
    
        # 기존 /etc/logrotate.d/syslog 목록에 /var/log/cmd.log 추가
        - name: 05. Ensure /var/log/cmd.log included in /etc/logrotate.d/syslog list
          ansible.builtin.lineinfile:
            path: /etc/logrotate.d/syslog
            regexp: '^/var/log/cmd\.log$'
            line: '/var/log/cmd.log'
            insertafter: '^/var/log/secure'
            backup: true
    
        - name: 06. Deploy /etc/profile.d/cmd_logging.sh
          ansible.builtin.copy:
            src: files/cmd_logging.sh
            dest: /etc/profile.d/cmd_logging.sh
            owner: root
            group: root
            mode: '0700'
    
        - name: 07. Ensure rsyslog rule for local1.notice -> /var/log/cmd.log
          ansible.builtin.lineinfile:
            path: /etc/rsyslog.conf
            regexp: '^local1\.notice'
            line: 'local1.notice                                               /var/log/cmd.log'
            insertafter: EOF
            backup: true
          notify: restart rsyslog
    
        # limits.conf 끝에 설정 (중복 방지)
        - name: 08. Append soft nofile to limits.conf
          ansible.builtin.lineinfile:
            path: /etc/security/limits.conf
            line: '* soft nofile 65535'
            regexp: '^\* +soft +nofile +\d+'
            insertafter: EOF
            create: true
            state: present
    
        - name: 09. Append hard nofile to limits.conf
          ansible.builtin.lineinfile:
            path: /etc/security/limits.conf
            line: '* hard nofile 65535'
            regexp: '^\* +hard +nofile +\d+'
            insertafter: EOF
            create: true
            state: present
    
        # user.conf 끝에 설정 (중복 방지)
        - name: 10. Ensure DefaultLimitNOFILE=65535 at end of /etc/systemd/user.conf
          ansible.builtin.lineinfile:
            path: /etc/systemd/user.conf
            regexp: '^#?DefaultLimitNOFILE='
            line: 'DefaultLimitNOFILE=65535'
            insertafter: EOF
            create: true
            state: present
    
        - name: 11. Reload systemd daemon
          ansible.builtin.systemd:
            daemon_reload: true
    ```
    
    - task01: 로그/스크립트 디렉터리 생성
    - task02: `logMoveDelete.sh` 배포 및 권한 설정
    - task03: 로그 관리 스크립트 cron 등록
    - task04: `/etc/logrotate.conf` 정책 설정
    - task05: `/etc/logrotate.d/syslog` 목록 설정
    - task06: `/etc/profile.d/cmd_logging.sh` 배포
    - task07: `rsyslog.conf` 설정 추가 + rsyslog 재시작
    - task08/09: `limits.conf`에 `nofile 65535` soft/hard 설정(중복 방지)
    - task10: `DefaultLimitNOFILE` 설정
    - task11: systemd daemon-reload
- files/cmd_logging.sh
    
    ```bash
    #!/bin/bash
    # bash 대화형 세션에서 명령을 로깅
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
    
- files/logMoveDelete.sh
    
    ```bash
    #!/bin/bash
    PATH=/usr/sbin:/usr/bin:/sbin:/bin
    
    source_dir="/var/log"
    backup_dir="/var/log/backup"
    
    # .tar.gz 파일을 백업 디렉토리로 이동
    mv "$source_dir"/*.tar.gz "$backup_dir" 2>/dev/null || true
    
    # 1년이 지난 파일 삭제
    find "$backup_dir" -name "*.tar.gz" -type f -mtime +365 -exec rm {} \;
    ```
    
- files/logrotate.conf
    
    ```bash
    daily
    rotate 365
    create
    dateext
    compress
    compresscmd /bin/gzip
    compressext .gz
    include /etc/logrotate.d
    
    /var/log/wtmp {
        monthly
        create 0664 root utmp
            minsize 1M
        rotate 1
    }
    
    /var/log/btmp {
        missingok
        monthly
        create 0600 root utmp
        rotate 1
    }
    ```
    

### 02. SAR 로그 설정

- sysstat 활성화 및 관리 주기 설정

```bash
---
- name: Enable and configure sysstat (sar)
  hosts: prod
  become: true

  handlers:
    - name: restart sysstat
      ansible.builtin.systemd:
        name: sysstat
        state: restarted
        enabled: true

  tasks:
    - name: 01. Ensure sysstat package is installed
      ansible.builtin.package:
        name: sysstat
        state: present

    - name: 02. Set HISTORY=365 in /etc/sysconfig/sysstat
      ansible.builtin.lineinfile:
        path: /etc/sysconfig/sysstat
        regexp: '^HISTORY='
        line: 'HISTORY=365'
        insertafter: EOF
        backup: true
      notify: restart sysstat

    - name: 03. Set COMPRESSAFTER=90 in /etc/sysconfig/sysstat
      ansible.builtin.lineinfile:
        path: /etc/sysconfig/sysstat
        regexp: '^COMPRESSAFTER='
        line: 'COMPRESSAFTER=90'
        insertafter: EOF
      notify: restart sysstat

    - name: 04. Set ZIP="gzip" in /etc/sysconfig/sysstat
      ansible.builtin.lineinfile:
        path: /etc/sysconfig/sysstat
        regexp: '^ZIP='
        line: 'ZIP="gzip"'
        insertafter: EOF
      notify: restart sysstat

    - name: 05. Install /etc/cron.d/sysstat
      ansible.builtin.copy:
        dest: /etc/cron.d/sysstat
        owner: root
        group: root
        mode: '0644'
        content: |
          */1 * * * * root /usr/lib64/sa/sa1 1 1
          53 23 * * * root /usr/lib64/sa/sa2 -A
      notify: restart sysstat

    - name: 06. Ensure sysstat service is enabled and started
      ansible.builtin.systemd:
        name: sysstat
        enabled: true
        state: started
```

- task01: sysstat 패키지 설치 확인
- task02, 03, 04: `/etc/sysconfig/sysstat` 설정
- task05: `/etc/cron.d/sysstat` 생성
- task06: sysstat 활성화
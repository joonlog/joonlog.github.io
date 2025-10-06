---
title : Cacti Conainer의 Script Tools
date : 2025-07-29 09:00:00 +09:00
categories : [Docker, Cacti]
tags : [docker, cacti, cacti container, cacti tools]  #소문자만 가능
---

### Backup

- 백업
    
    ```bash
    docker exec <docker image ID or name> ./backup.sh
    ```
    
    - 루트 Cacti/Spine 디렉토리 복사 및 cacti DB 덤프 수행
    - cacti container의 /backups 디렉토리 아래에 tar.gz 형식으로 압축 백업 저장
    - `BACKUP_RETENTION` 값 조정. 기본값 7개
- 백업 복원
    
    ```bash
    docker exec <docker image ID or name> ./restore.sh /backups/<filename>
    ```
    
- 백업 목록 조회
    
    ```bash
    docker exec <docker image ID or name> ls /backups
    ```
    

### Update

```
docker exec <docker image ID or name> ./upgrade.sh
```

- 버전 지정 시 upgrade.sh 수정
    - [Cacti Version Links](http://www.cacti.net/downloads)
    - [Spine Version Links](http://www.cacti.net/downloads/spine)

### Template

- 시작 시 스크립트가 가져와서 설치됨

```bash
├── templates
│   ├── template_name.xml
│   ├── resource
│   │   └── script_queries
│   │       └── ...
│   │   └── script_server
│   │       └── ...
│   │   └── snmp_queries
│   │       └── ...
│   ├── scripts
│   │   └── ...
```

### Plugins

- plugins 디렉토리에 플러그인을 넣는 것으로 부팅 시 자동 로드
- Cacti GUI에서 플러그인 활성화 필요
- 컨테이너 동작 후 플러그인 추가하려면 docker volume을 통해 마운트 필요

### Settings

- sql 변경사항을 settings 디렉토리 아래에 배치하여 초기 설치 시 설정 가능
    - start.sh가 설치 시 모든 *.sql을 병합하기 때문
    - ex) settings/spine.sql은 spine 활성화하기 위한 sql
        
        ```bash
        --
        -- Enable spine poller from docker installation
        --
        
        REPLACE INTO `%DB_NAME%`.`settings` (`name`, `value`) VALUES('path_spine', '/spine/bin/spine');
        REPLACE INTO `%DB_NAME%`.`settings` (`name`, `value`) VALUES('path_spine_config', '/spine/etc/spine.conf');
        REPLACE INTO `%DB_NAME%`.`settings` (`name`, `value`) VALUES('poller_type', '2');
        ```
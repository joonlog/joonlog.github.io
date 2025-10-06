---
title : Cacti Conainer 구조
date : 2025-07-27 09:00:00 +09:00
categories : [Docker, Cacti]
tags : [docker, cacti, cacti container]  #소문자만 가능
---

> https://github.com/scline/docker-cacti
> 

### Dockerfile

- https://github.com/scline/docker-cacti/blob/master/Dockerfile
- Rocky 9.0을 베이스 이미지로 사용
- 환경 변수로 DB 접속 정보, URL 경로, PHP 기본값 설정
- start.sh 복사 및 실행
- cacti + spine tar.gz 복사
- 서포트 파일(configs, configs/crontab, backup, upgrade, restore) 복사
- **RUN**
    - sh 파일 실행 권한 추가
    - 볼륨 디렉토리 /backups, /cacti, /spine 생성
    - 패키지 업데이트 및 설치(php, mariadb, 소스 설치 도구 등)

### start.sh

- https://github.com/scline/docker-cacti/blob/master/start.sh
- 초기 환경 설정
    - PHP 설정, 시스템 TZ 설정
- 새 설치 여부 판단
    - `/cacti/install.lock` 파일의 존재 여부로 이미 설치된 상태인지 판단
    - Cacti와 Spine 파일 압축 해제, 필요한 디렉토리로 이동
    - 템플릿 파일, 플러그인, 설정 파일들을 적절한 위치에 복사하여 기본 환경을 구성
    - DB가 준비될 때까지 대기 후, 데이터베이스 생성 및 사용자 권한 설정, 그리고 cacti.sql 파일을 DB에 적용하여 초기 데이터를 세팅합니다.
    - `REMOTE_POLLER`가 1이면, 마스터와의 설정 연동을 위한 추가 작업(예, config.php 수정, 파일 권한 조정 등)을 수행합니다.
- 서비스 시작
    - crond, snmpd, php-fpm, apache 실행

### Compose 구조

![CactiContainerArchitecture1.png](/assets/img/linux/CactiContainerArchitecture1.png)

- https://github.com/scline/docker-cacti/blob/master/docker-compose/cacti_multi_shared.yml
- Cacti Master+ Cacti Master DB + Poller + Poller DB
    
    **Cacti**
    
    - 위 Dockerfile로 생성한 cacti 이미지 사용
    - 웹 인터페이스 및 기본 Cacti 기능용 주 컨테이너
    - 환경변수로 db 접속 정보 설정
    - /cacti, /spine, /backups 볼륨 공유
    
    **Poller**
    
    - 위 Dockerfile로 생성한 cacti 이미지 사용
    - 환경변수 `REMOTE_POLLER=1`가 설정되어 원격 풀러(데이터 수집 전용)로 동작
    - cacti_poller DB 접속 정보와 동시에 마스터 DB 접속 정보를 전달받아, 마스터로부터 설정을 가져와서 폴링 작업만 수행하도록 구성
    - /cacti, /spine, /backups 볼륨 공유
    
    **DB**
    
    - MariaDB 10.3 이미지 사용하며, Cacti와 Poller가 각각 DB 사용
    
    **Network**
    
    - 4개 컨테이너가 같은 네트워크공유
    
    **Volume**
    
    - Cacti와 Poller가 호스트의 /cacti, /spine, /backups 볼륨 공유
    - /var/lib/mysql를 db volume으로 사용
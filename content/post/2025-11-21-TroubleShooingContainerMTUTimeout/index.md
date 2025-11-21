---
title: "컨테이너 환경에서 MTU 불일치로 인한 네트워크 타임아웃 이슈 TroubleShooting"
date: 2025-11-21T09:00:00+09:00
categories: ["Container", "Network"]
tags: ["container", "network", "mtu", "ncp mtu"]
---


> 컨테이너 환경에서 서버와 컨테이너의 MTU 불일치로 인한 네트워크 이슈를 정리한 글
> 

### 환경

- NCP 서버 3대
- 구조: Nginx 컨테이너 서버 → PHP-FPM 컨테이너 서버 → DB 서버
- MariaDB 10.3

### 문제 상황

1. 구축한 사이트 이용 간 아무 패턴 없이 간헐적으로 타임아웃 발생
    - 이슈 발생 주기는 완전 랜덤
    - 사이트 내 모든 경로에서 간헐적으로 발생
2. 서버 리소스나 부하는 모두 정상

### 에러 분석

- 처음에는 애플리케이션 단 문제로 생각해서 점검
    - 세션 분리는 이미 되어 있는 상태
    - 가비지 콜렉션 설정은 없어서 GC 설정 했으나 증상 동일
- DB쪽 지표 모두 정상

### 원인

- 네트워크 쪽을 점검해보니 서버와 컨테이너의 MTU가 차이가 많이 나는 것을 확인
    - NCP는 서버의 기본 MTU가 8590이라고 함
    - 컨테이너는 기본값 사용 중이었으니 1500
    - MTU 값 불일치로 발생했던 이슈로 확인
        - 크기가 작은 요청은 정상이다가 일부 큰 패킷들이 드랍 → php 세션 끊어짐 → PHP가 요청을 끝까지 처리하지 못해서 타임아웃

### 해결

1. Docker 데몬 MTU 수정
    - `/etc/docker/daemon.json`
        - 수정 후 docker 재기동 필요
        
        ```bash
        {
          "mtu": 8950
        }
        ```
        
2. 컨테이너 네트워크 MTU 수정
    - compose의 경우 mtu 값 명시
        
        ```bash
        networks:
          default:
            driver: bridge
            driver_opts:
              com.docker.network.driver.mtu: 8950
        ```
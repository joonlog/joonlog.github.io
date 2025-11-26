---
title: "Portfolio"
slug: "portfolio"
layout: "page"
menu:
    main:
        weight: 6
        params:
            icon: user
comments: false
---

경험을 기록으로 남기며 성장하는 1년 차 클라우드 엔지니어입니다.
레거시 환경의 클라우드 마이그레이션부터 컨테이너 오케스트레이션까지, **실무에서 직접 설계하고 구축한 경험**을 중심으로 정리했습니다.

---

## 주요 프로젝트

### 공공 G사 온프레미스 → NCP 클라우드 마이그레이션 & 컨테이너화
*2024년 10월 - 2024년 12월 (예정) | Naver Cloud Platform*

**레거시 온프레미스 환경을 NCP 기반 Docker 컨테이너 환경으로 전환하고, 전사 CI/CD 파이프라인을 구축한 프로젝트**

#### 프로젝트 배경
- 온프레미스 환경의 유지보수 비용 증가 및 확장성 한계
- SVN 기반 수동 배포로 인한 휴먼 에러 및 배포 시간 지연
- 클라우드 환경으로의 전환 및 개발 프로세스 현대화 필요

#### 프로젝트 규모
- **서비스**: 6개 (Java Spring 1.8, PHP-FPM 8.3)
- **서버**: Web/WAS/DB 각 2대 HA 구성 × 6개 사이트 = **18대** + GitLab, Jenkins, ELK, Zabbix, Matomo 등 솔루션 서버
- **아키텍처**: Nginx - WAS - MariaDB 3티어 컨테이너 구조
- **일정**: 10월 KT Cloud 테스트 환경 구축 → 11월 NCP 본격 마이그레이션 → 12월 중순 기관 최종 발표

#### 담당 역할 (95% 단독 수행)

**1. 인프라 설계 및 구축**
- NCP 리소스 설계 및 구축 (Server, Load Balancer, ACG, VPC)
- GitLab + GitLab Container Registry 구축 및 리포지토리 관리
- Jenkins 구축 및 CI/CD 파이프라인 설계
- SVN → Git 형상 관리 체계 전환

**2. 컨테이너화 및 CI/CD 구현**
- PHP-FPM 8.3, Spring Java 1.8 애플리케이션 Dockerfile 작성
- docker-compose 기반 멀티 컨테이너 환경 구성
- GitLab Webhook 기반 자동 CI/CD 파이프라인 구축
  - Main 브랜치 Protected 설정
  - Push 트리거 → CI (소스 빌드 → 이미지 빌드 → Registry 푸시) → CD (자동 배포)
- GitLab Container Registry를 활용한 이미지 버전 관리

**3. 자동화**
- Ansible Playbook을 통한 서버 초기 설정 자동화
  - sar, logrotate 설정
  - Docker 설치 및 daemon.json 관리
  - /etc/hosts 파일 관리

#### 기술적 챌린지

**1. MTU 불일치로 인한 간헐적 네트워크 타임아웃**
- **문제**: 아무 패턴 없이 간헐적으로 사이트 타임아웃 발생 (서버 리소스는 정상)
- **원인**: 호스트(NCP MTU 8590) vs 컨테이너(기본 MTU 1500) 불일치로 큰 패킷 손실
- **해결**: Docker daemon.json 및 docker-compose에서 MTU 8950으로 통일
- **상세**: [컨테이너 환경에서 MTU 불일치로 인한 네트워크 타임아웃 이슈 Troubleshooting](https://joonlog.github.io/p/%EC%BB%A8%ED%85%8C%EC%9D%B4%EB%84%88-%ED%99%98%EA%B2%BD%EC%97%90%EC%84%9C-mtu-%EB%B6%88%EC%9D%BC%EC%B9%98%EB%A1%9C-%EC%9D%B8%ED%95%9C-%EB%84%A4%ED%8A%B8%EC%9B%8C%ED%81%AC-%ED%83%80%EC%9E%84%EC%95%84%EC%9B%83-%EC%9D%B4%EC%8A%88-troubleshooting/)

**2. GitLab Container Registry 이미지 푸시 실패**
- **문제**: 사설 IP로 Registry 접근 시 푸시 중단 (GitLab 웹콘솔은 공인 IP 사용)
- **원인**: GitLab이 내부적으로 공인 IP로 redirect하면서 연결 끊김
- **해결**: 도메인 기반 접근으로 전환
  - DNS A 레코드: 도메인 → 공인 IP
  - 이미지 푸시 서버 /etc/hosts: 도메인 → 사설 IP
  - 불필요한 공인 IP 트래픽 절감 (이미지당 1GB+ 절약)

#### 성능 최적화

**웹 애플리케이션 응답 시간 최적화 (개발팀 협업)**
- **개선 전**: 평균 3초 이상
- **개선 후**: 평균 150ms
- **개선율**: **약 20배 성능 향상**

**최적화 작업**:
- Nginx gzip 압축 설정 (4096+ 파일 대상)
- Nginx keepalive 설정
- MariaDB 인덱스 추가 (개발팀 협업)
- WAS/DB connection pool 및 buffer 튜닝
- hey, jmeter 부하 테스트를 통한 CPU, Memory, Connection 수치 기반 튜닝

#### 사용 기술
- **Cloud**: Naver Cloud Platform (Server, Load Balancer, ACG, VPC)
- **Container**: Docker, docker-compose
- **CI/CD**: GitLab (SCM, Container Registry, Webhook), Jenkins (Pipeline)
- **Automation**: Ansible
- **Monitoring**: ELK, Zabbix, Matomo
- **Load Testing**: hey, jmeter

**관련 포스트**:
- [컨테이너 환경에서 Jenkins의 Java CICD 파이프라인 구축하기](https://joonlog.github.io/p/%EC%BB%A8%ED%85%8C%EC%9D%B4%EB%84%88-%ED%99%98%EA%B2%BD%EC%97%90%EC%84%9C-jenkins%EC%9D%98-java-cicd-%ED%8C%8C%EC%9D%B4%ED%94%84%EB%9D%BC%EC%9D%B8-%EA%B5%AC%EC%B6%95%ED%95%98%EA%B8%B0/)
- [컨테이너 환경에서 Jenkins의 PHP CICD 파이프라인 구축하기](https://joonlog.github.io/p/%EC%BB%A8%ED%85%8C%EC%9D%B4%EB%84%88-%ED%99%98%EA%B2%BD%EC%97%90%EC%84%9C-jenkins%EC%9D%98-php-cicd-%ED%8C%8C%EC%9D%B4%ED%94%84%EB%9D%BC%EC%9D%B8-%EA%B5%AC%EC%B6%95%ED%95%98%EA%B8%B0/)
- [컨테이너 환경에서 MTU 불일치로 인한 네트워크 타임아웃 이슈 Troubleshooting](https://joonlog.github.io/p/%EC%BB%A8%ED%85%8C%EC%9D%B4%EB%84%88-%ED%99%98%EA%B2%BD%EC%97%90%EC%84%9C-mtu-%EB%B6%88%EC%9D%BC%EC%B9%98%EB%A1%9C-%EC%9D%B8%ED%95%9C-%EB%84%A4%ED%8A%B8%EC%9B%8C%ED%81%AC-%ED%83%80%EC%9E%84%EC%95%84%EC%9B%83-%EC%9D%B4%EC%8A%88-troubleshooting/)

---

### 공공 H사 KT Cloud 운영 자동화
*2024년 11월 - 현재 | KT Cloud*

**3티어 인프라 운영 및 Ansible 기반 자동화**

- **규모**: 10개 시스템, 40여대 서버
- **주요 성과**:
  - Ansible 기반 보안 취약점 점검 및 조치 자동화
  - Shell Script를 활용한 WEB/WAS 로그 관리 및 백업 이중화
  - Nagios/Munin 기반 메트릭 모니터링 및 장애 알람 운영

---

### 셀프 매니지드 Kubernetes 클러스터 구축 (개인 프로젝트)
*2024년 8월 - 현재*

**CSP 없이 자체 관리형 Kubernetes 클러스터를 구축하며 내부 동작 원리 학습**

**구축 환경**:
- MetalLB + HAproxy + Nginx Ingress Controller (네트워크)
- Longhorn (분산 스토리지, 고가용성 PV)
- Jenkins, GitLab, Prometheus, Grafana (Helm Chart 기반 배포)

**주요 성과**:
- CSP의 ALB/EBS를 대체하는 오픈소스 기반 네트워크/스토리지 구성
- 115개 이상의 기술 블로그 포스트 작성 (학습 과정 문서화)

**관련 포스트**: [Kubernetes 카테고리](https://joonlog.github.io/categories/kubernetes/)

---

## 기술 블로그

**115개 이상의 기술 포스트** (2024년 11월 - 현재)

실무 경험과 학습 내용을 체계적으로 문서화하고 있습니다.

- [Kubernetes](https://joonlog.github.io/categories/kubernetes/) - 셀프 매니지드 클러스터 구축
- [Container](https://joonlog.github.io/categories/container/) - Docker, CI/CD 파이프라인
- [Linux](https://joonlog.github.io/categories/linux/) - 시스템 관리 및 자동화

**블로그**: [joonlog.github.io](https://joonlog.github.io)

---

## 연락처

- **Email**: ksi05298@gmail.com
- **Blog**: [joonlog.github.io](https://joonlog.github.io)
- **GitHub**: [github.com/joonlog](https://github.com/joonlog)
- **Resume**: [이력서 보기](https://joonlog.github.io/resume/)

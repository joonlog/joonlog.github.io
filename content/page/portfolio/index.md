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
실무 프로젝트와 개인 학습을 통해 **컨테이너 오케스트레이션**과 **CI/CD 자동화** 역량을 쌓아왔습니다.

---

## 주요 프로젝트

### 공공 G사 NCP 컨테이너 마이그레이션 프로젝트
*2025년 10월 - 2025년 12월 (예정) | Naver Cloud Platform*

**온프레미스 환경을 NCP 기반 Docker 컨테이너 환경으로 마이그레이션**

- **규모**: 6개 서비스, 30여대 서버
- **아키텍처**: Nginx-WAS-MariaDB 3티어 컨테이너 구조
- **기술 스택**: Docker, Jenkins, GitLab, Ansible

**주요 성과**:
- PHP-FPM 8.3 및 Spring Java 1.8 애플리케이션을 컨테이너화
- GitLab Container Registry를 활용한 이미지 버전 관리 체계 구축
- Jenkins 기반 CI/CD 파이프라인 구현 (소스 빌드 → 이미지 빌드 → 자동 배포)
- Ansible을 통한 서버 환경 설정 표준화 및 배포 자동화

**관련 포스트**:
- [컨테이너 환경에서 Jenkins의 Java CICD 파이프라인 구축하기](https://joonlog.github.io/p/%EC%BB%A8%ED%85%8C%EC%9D%B4%EB%84%88-%ED%99%98%EA%B2%BD%EC%97%90%EC%84%9C-jenkins%EC%9D%98-java-cicd-%ED%8C%8C%EC%9D%B4%ED%94%84%EB%9D%BC%EC%9D%B8-%EA%B5%AC%EC%B6%95%ED%95%98%EA%B8%B0/)
- [컨테이너 환경에서 Jenkins의 PHP CICD 파이프라인 구축하기](https://joonlog.github.io/p/%EC%BB%A8%ED%85%8C%EC%9D%B4%EB%84%88-%ED%99%98%EA%B2%BD%EC%97%90%EC%84%9C-jenkins%EC%9D%98-php-cicd-%ED%8C%8C%EC%9D%B4%ED%94%84%EB%9D%BC%EC%9D%B8-%EA%B5%AC%EC%B6%95%ED%95%98%EA%B8%B0/)

---

### 공공 H사 KT Cloud 운영 프로젝트
*2024년 11월 - 현재 | KT Cloud*

**KT Cloud 기반 3티어 인프라 운영 및 자동화**

- **규모**: 10개 시스템, 40여대 서버
- **기술 스택**: KT Cloud, Ansible, Nagios/Munin, Shell Script

**주요 성과**:
- Ansible 기반 서버 환경 설정 및 모니터링 에이전트 배포 자동화
- 보안 취약점 점검 및 조치 자동화 (공공기관 보안 점검 대응)
- Shell Script를 활용한 WEB/WAS 로그 관리 및 백업 이중화 자동화
- Nagios/Munin 기반 메트릭 모니터링 및 장애 알람 체계 운영

---

### 셀프 매니지드 Kubernetes 클러스터 구축 (개인 프로젝트)
*2025년 8월 - 현재*

**CSP 없이 자체 관리형 Kubernetes 클러스터를 구축하고 운영**

**목표**: 클라우드 환경 없이도 프로덕션급 Kubernetes 환경을 직접 구축하며 내부 동작 원리 학습

**구축 환경**:
- **네트워크**: MetalLB (L4 LB) + HAproxy (외부 라우팅) + Nginx Ingress Controller
- **스토리지**: Longhorn (분산 블록 스토리지, 고가용성 PV)
- **CI/CD**: Jenkins (Helm Chart 기반)
- **SCM**: GitLab (Container Registry 포함)
- **모니터링**: Prometheus + Grafana (kube-prometheus-stack)
- **접근 제어**: Teleport (클러스터 접근 관리)

**주요 성과**:
- CSP의 ALB/EBS를 대체하는 오픈소스 기반 네트워크/스토리지 구성
- Helm Chart 기반 애플리케이션 배포 및 GitOps 파이프라인 구축
- 115개 이상의 기술 블로그 포스트 작성 (학습 과정 문서화)

**관련 포스트**:
- [외부에서 자체 관리형 Kubernetes 접근을 위한 MetalLB/HAproxy 설정과 통신 구조](https://joonlog.github.io/p/%EC%99%B8%EB%B6%80%EC%97%90%EC%84%9C-%EC%9E%90%EC%B2%B4-%EA%B4%80%EB%A6%AC%ED%98%95-kubernetes-%EC%A0%91%EA%B7%BC%EC%9D%84-%EC%9C%84%ED%95%9C-metallbhaproxy-%EC%84%A4%EC%A0%95%EA%B3%BC-%ED%86%B5%EC%8B%A0-%EA%B5%AC%EC%A1%B0/)
- [자체 관리형 Kubernetes에서의 Jenkins 구축](https://joonlog.github.io/p/%EC%9E%90%EC%B2%B4-%EA%B4%80%EB%A6%AC%ED%98%95-kubernetes%EC%97%90%EC%84%9C%EC%9D%98-jenkins-%EA%B5%AC%EC%B6%95/)
- [자체 관리형 Kubernetes에서의 분산 스토리지 Longhorn 구축](https://joonlog.github.io/p/%EC%9E%90%EC%B2%B4-%EA%B4%80%EB%A6%AC%ED%98%95-kubernetes%EC%97%90%EC%84%9C%EC%9D%98-%EB%B6%84%EC%82%B0-%EC%8A%A4%ED%86%A0%EB%A6%AC%EC%A7%80-longhorn-%EA%B5%AC%EC%B6%95/)
- [자체 관리형 Kubernetes에서의 Prometheus와 Grafana 구축](https://joonlog.github.io/p/%EC%9E%90%EC%B2%B4-%EA%B4%80%EB%A6%AC%ED%98%95-kubernetes%EC%97%90%EC%84%9C%EC%9D%98-prometheus%EC%99%80-grafana-%EA%B5%AC%EC%B6%95/)
- [자체 관리형 Kubernetes에서의 GitLab 구축](https://joonlog.github.io/p/%EC%9E%90%EC%B2%B4-%EA%B4%80%EB%A6%AC%ED%98%95-kubernetes%EC%97%90%EC%84%9C%EC%9D%98-gitlab-%EA%B5%AC%EC%B6%95/)

**카테고리 전체 보기**: [Kubernetes](https://joonlog.github.io/categories/kubernetes/)

---

## 기술 블로그

**115개 이상의 기술 포스트 작성** (2024년 11월 - 현재)

학습한 내용과 실무 경험을 체계적으로 문서화하고 있습니다.

**주요 카테고리**:
- [Kubernetes](https://joonlog.github.io/categories/kubernetes/) - 셀프 매니지드 클러스터 구축 및 운영
- [Linux](https://joonlog.github.io/categories/linux/) - 시스템 관리 및 자동화
- [Container](https://joonlog.github.io/categories/container/) - Docker, CI/CD 파이프라인
- [SCM](https://joonlog.github.io/categories/scm/) - GitLab, GitHub Actions

**블로그**: [joonlog.github.io](https://joonlog.github.io)

---

## 기술 스택

**Linux & Automation**
- RHEL 계열 시스템 관리 및 운영
- Ansible 기반 인프라 자동화 (배포, 설정, 보안 점검)
- Shell Script를 활용한 로그 관리 및 백업 자동화

**Cloud & Container**
- CSP: AWS, KT Cloud, Naver Cloud Platform
- Container: Docker, Containerd (nerdctl, buildkit)
- Orchestration: Kubernetes (자체 관리형 클러스터 구축 경험)

**CI/CD**
- Jenkins (파이프라인 설계 및 컨테이너 에이전트 구성)
- GitLab (SCM, Container Registry, GitOps)
- Maven, Docker 기반 빌드/배포 자동화

**Monitoring & Storage**
- Monitoring: Nagios, Munin, Prometheus, Grafana
- Storage: Longhorn (분산 블록 스토리지)

**Network**
- Load Balancer: MetalLB, HAproxy, Nginx Ingress Controller
- CSP LB: AWS ALB, KT/NCP Load Balancer

---

## 연락처

- **Email**: ksi05298@gmail.com
- **Blog**: [joonlog.github.io](https://joonlog.github.io)
- **GitHub**: [github.com/joonlog](https://github.com/joonlog)
- **Resume**: [Resume 페이지](https://joonlog.github.io/resume/)

---
title : Apache OPTIONS * 루프백 요청으로 인한 CPU 과부하 TroubleShooting
date : 2025-04-25 09:00:00 +09:00
categories : [Linux, Troubleshooting]
tags : [linux, apache, mpm, troubleshooting] #소문자만 가능
---

### 문제 상황

- 운영 중인 CentOS 서버에서 너무 많은 수의 `httpd` 프로세스가 CPU를 과점유하고, Load Average가 30 이상까지 치솟는 현상이 발생

### 원인 분석

- `access_log` 확인 결과, 루프백 주소(`::1`)에서 `OPTIONS * HTTP/1.0` 요청이 수 초 간격으로 다수 발생
- User-Agent는 `"Apache/... (internal dummy connection)"`으로, Apache 자체가 생성하는 dummy connection임을 확인
    - 이 요청은 Apache의 internal dummy connection이며, 일반적으로는 MPM이 유휴 프로세스를 관리하기 위해 보내는 내부 테스트용 요청
    - 하지만 일정 이상 발생하게 되면, 해당 요청을 처리하기 위한 httpd 프로세스가 과도하게 생성되고 정리되지 않아, 리소스를 점유

```bash
# tail -n 1000 /var/log/httpd/access_log | grep "::1"

::1 - - [22/Apr/2025:07:01:55 +0000] "OPTIONS * HTTP/1.0" 200 - "-" "Apache/2.4.6 (CentOS) OpenSSL/1.0.2k-fips (internal dummy connection)"
::1 - - [22/Apr/2025:07:02:34 +0000] "OPTIONS * HTTP/1.0" 200 - "-" "Apache/2.4.6 (CentOS) OpenSSL/1.0.2k-fips (internal dummy connection)"
::1 - - [22/Apr/2025:07:02:35 +0000] "OPTIONS * HTTP/1.0" 200 - "-" "Apache/2.4.6 (CentOS) OpenSSL/1.0.2k-fips (internal dummy connection)"
::1 - - [22/Apr/2025:07:02:36 +0000] "OPTIONS * HTTP/1.0" 200 - "-" "Apache/2.4.6 (CentOS) OpenSSL/1.0.2k-fips (internal dummy connection)"
::1 - - [22/Apr/2025:07:02:42 +0000] "OPTIONS * HTTP/1.0" 200 - "-" "Apache/2.4.6 (CentOS) OpenSSL/1.0.2k-fips (internal dummy connection)"
::1 - - [22/Apr/2025:07:02:56 +0000] "OPTIONS * HTTP/1.0" 200 - "-" "Apache/2.4.6 (CentOS) OpenSSL/1.0.2k-fips (internal dummy connection)"
```

```bash
# httpd -V | grep MPM

[Fri Apr 25 04:35:44.665826 2025] [so:warn] [pid 31382] AH01574: module rewrite_module is already loaded, skipping
Server MPM:     prefork
```

- MPM 관련 설정(StartServers, MaxRequestWorkers)이 명시되지 않은 상태
    - Apache 기본값(최대 256개 프로세스)이 적용된 것으로 추정

### 조치 내용

> 하기 내용은 Apache 설정을 수정하지 않고 즉시 대응 가능한 임시 조치
> 
1. IPv6 루프백 dummy connection 차단을 위한 ip6tables 규칙 추가
- 루프백(`::1`)에서 들어오는 `OPTIONS *` 요청만 차단
- 외부 요청이나 정상적인 HTTP 트래픽에는 영향 없음

```bash
ip6tables -I INPUT -s ::1 -p tcp --dport 80 -m string --algo bm --string "OPTIONS *" -j DROP
```

2. Apache 재기동
- 이미 생성되어 남아 있는 과도한 `httpd` 프로세스를 정리하고 초기화
- 신규 요청부터는 차단 규칙이 적용된 상태로 운영

```bash
systemctl restart httpd
```

3. CPU 사용률 및 Load Average 수 분 내로 정상화

### 결론

- `httpd` 프로세스 수 안정화, 더 이상 루프백 `OPTIONS *` 요청 발생하지 않음

### 추가 조치 방안

- 문제의 근본 원인을 해결하기 위해서는 Apache 설정 최적화 필요
- 특히, 현재 서버는 prefork MPM을 사용하고 있으며, 기본 설정값으로 운영중
- 이로 인해 부하 상황에서 과도한 프로세스 생성 및 dummy connection 폭주 현상이 발생한 것으로 분석됨
- mpm 설정 추가
    - MaxRequestWorkers를 서버 사양에 맞게 제한

```bash
# /etc/httpd/conf.modules.d/00-mpm.conf

<IfModule mpm_prefork_module>
    StartServers        5
    MinSpareServers     5
    MaxSpareServers    10
    MaxRequestWorkers 150
    ServerLimit        150
</IfModule>
```
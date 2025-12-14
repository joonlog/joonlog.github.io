---
title: "NCP에서 NLB 구성 후 503 발생 TroubleShooting"
date: 2025-12-14T09:00:00+09:00
categories: ["NCP", "Load Balancer"]
tags: ["ncp", "load balancer", "nlb", "network load balancer"]
---


> NLB 사용 시 Source IP 유지로 인해 발생한 ACG 방화벽 정책 이슈 해결
> 

### 문제 상황

- NCP에서 WEB → NLB → WAS 구조로 서비스 구성
- 구축 단계에서는 아웃바운드 ANY OPEN으로 사이트 정상 동작
- 운영 전환을 위해 아웃바운드 제한하니 사이트에서 503 오류 발생
- 인바운드는 그대로였으므로 아웃바운드 정책 문제로 확인
    - WEB 서버에서 NLB의 Private IP로 아웃바운드를 허용했던 상황

### 원인 분석

- WAS에서 tcpdump 확인
    - Source IP가 NLB Private IP일 줄 알았으나, 덤프 확인 시 WEB 서버 Private IP로 출력됨
        - 트래픽을 단순 포워딩만 하고 클라이언트 IP를 그대로 전달
- 결과적으로 WAS는 NLB가 아닌 WEB 서버와 직접 통신하는 형태가 됨
- 동작 방식은 AWS에서의 NLB도 동일하며 L4 LB의 정상적인 구조로 확인

### 조치

- WEB 서버의 아웃바운드에 WAS 서버의 Private IP 추가 후 정상 동작 확인

### 결론

L4 LB인 NLB NAT를 하지 않기 때문에 WAS는 항상 앞단 서버의 Real IP를 보게 된다.

따라서 ACG 정책도 NLB가 아닌 실제 통신 경로(WEB → WAS) 기반으로 설정해야 한다.
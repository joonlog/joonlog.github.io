---
title: "kube-proxy 모드(iptables vs IPVS)"
slug: "kube-proxy-iptables-vs-ipvs"
aliases: ["/p/kube-proxy-모드iptables-vs-ipvs/"]
date: 2025-08-19T09:00:00+09:00
categories: ["Kubernetes", "Architecture"]
tags: ["Kubernetes", "iptables", "IPVS"]
---

# Kubernetes kube-proxy: iptables vs IPVS

- K8S에서 Service를 만들면, 클러스터 내부 Pod는 `ClusterIP`를 통해 접근 가능

    - 이때 트래픽을 실제 Pod로 보내주는 역할을 하는 게 `kube-proxy`

- 리눅스에서 kube-proxy는 iptables, IPVS, nftables 세 모드 중 하나로 동작

    - 기본값은 아직 iptables이고, 앞으로 nftables로 바뀔 예정

    - IPVS는 1.35에서 deprecated, 1.40부터 기본 비활성, 1.43에서 제거 예정

---

## iptables 모드

K8S 초창기부터 기본으로 사용된 방식

- 동작 원리

    - service/endpoint 정보를 iptables 규칙에 직접 기록

    - 커널이 NAT 체인을 따라가면서 Pod로 트래픽을 전송

- 장점

    - 별도 모듈 설치 없이 바로 사용 가능

    - 오래 쓰여온 만큼 안정적이고, 문제 생겼을 때 `iptables -t nat -L -n` 으로 바로 확인 가능

- 단점

    - 엔드포인트 수가 많아지면 규칙도 폭발적으로 늘어나 성능 저하가 생길 수 있음

    - 규칙을 갱신할 때 전체 체인을 다시 써야 해서 순간 부하 발생

---

## IPVS 모드

리눅스 커널의 IPVS 기능을 사용

- 동작 원리

    - `ip_vs` 모듈을 이용해 L4 로드밸런서처럼 동작

    - kube-proxy는 service/endpoint를 커널에 등록만 해두고, 실제 트래픽 처리는 커널이 빠르게 처리

- 장점

    - 수천~수만 엔드포인트까지 안정적으로 처리 가능

    - 엔드포인트 변경 시 필요한 부분만 업데이트 → 부하 적음

    - 라운드로빈, 세션 해시 등 다양한 로드밸런싱 알고리즘 지원

- 단점

    - 커널 모듈(`ip_vs`)과 유틸리티(`ipvsadm`, `ipset`) 설치 필요

    - iptables에 비해 디버깅이 낯설 수 있음 (`ipvsadm -Ln` 확인 필요)

---

## nftables 모드

리눅스 nftables(iptables의 후속 패킷 필터)를 사용. 1.33에서 GA

- iptables 모드와 IPVS 모드를 모두 대체하는 방식으로, 둘보다 성능이 좋고 IPVS 대체로 권장됨

- 규칙을 서비스 단위로 갱신해서, 서비스가 많아도 iptables처럼 전체 규칙을 다시 쓰지 않음

- 커널 5.13 이상 필요, 확인은 `nft list ruleset`

---

### 선택 기준

- iptables 모드

    - 수십~수백 노드, 수천 개 엔드포인트 수준

    - 운영 단순성, 문제 대응 속도 중시

- nftables 모드

    - 새로 구축하는 클러스터(커널 5.13 이상)

    - IPVS를 쓰던 대규모 환경의 이전 대상

- IPVS 모드 (deprecated, 기존 클러스터만)

    - 수천 노드 이상, 서비스 엔드포인트 수가 수만 단위

    - 대규모 SaaS 환경, 트래픽이 매우 많은 서비스

    - 고급 로드밸런싱(세션 어피니티 등) 필요한 경우

> kubespray는 기본값이 지금도 IPVS(`kube_proxy_mode: ipvs`)라, 새로 구축한다면 `nftables`로 바꾸는 것을 검토

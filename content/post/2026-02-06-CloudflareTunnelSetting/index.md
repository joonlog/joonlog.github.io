---
title: "홈서버 접근 용도로 Cloudflare Tunnel을 쓰는 이유"
date: 2026-02-05T09:00:00+09:00
categories: ["Home Server", "Cloudflare"]
tags: ["home server", "cloudflare", "cloudflare tunnel", "promox"]
---


> 외부에서 홈서버에 안전하게 접근하기 위한 설정
> 

홈서버를 디폴트 상태로 외부에서 접근하기 위해선 IP와 Port로 접근이 필요하다.

클라우드라면 Security Group / ACL 등등으로 제어가 쉽기 때문에 상관이 없지만, 홈서버는 집에서 사용하는 공유기 공인 IP를 외부에 오픈해야하기 때문에 보안적으로 위험하다.

기본적으로 홈 IP는 AnyOpen되어 있고, 외부에서 접근한다는 것 자체가 Wifi든 모바일이든 매번 접근하는 IP가 가변적이기 때문에 IP로 접근을 차단하는 것에도 무리가 있다.

특히 Proxmox 관리 콘솔(PVE)은 관리자 권한 탈취될 시 홈 네트워크 전체에 대한 통제권을 넘길 수 있기 때문에 매우 위험하다.

따라서 `Cloudflare Tunnel`로 공인 IP를 숨기고, `Cloudflare Access`로 인증된 사용자만 접근할 수 있게 하려고 한다.

## **1. Cloudflare Tunnel**

- IP와 포트를 공개하지 않고도 외부에서 서버로 접속할 수 있게 하는 서비스
- 외부에서 Cloudflare를 거쳐서 서버에 접속
- 기존 인바운드+포트 포워딩 기반 접근과 다른 아웃바운드 기반 리버스 프록시 연결

### Cloudflare Tunnel을 쓰는 이유

- 홈서버는 가정용 공유기에 부여된 공인 IP로 접근해야 하기 때문에 보안에 취약
- 기업에서 쓰는 서버들의 경우 WAF, IPS, F/W 등 방어 솔루션이 많지만, 홈서버는 사실상 방어 수단이 전무
- Cloudflare를 거치면서 레이턴시가 늘어나더라도, 홈서버는 보안이 최우선
- 따라서 Zero-Inbound를 유지하고도 운영이 가능한 Cloudflare Tunnel 사용

### **기존 포트포워딩 방식**

1. 포트포워딩은 외부에서 들어오는 Inbound 연결을 서버가 직접 받아들이는 구조
    
    ```
    외부 클라이언트 → 공유기(포트포워딩) → 홈서버
    ```
    
    - 이 방식에서는 서버가 외부에서 접근 가능한 상태가 되며, 열려 있는 포트는 스캔이나 공격의 대상이 될 수 있어서 보안적으로 매우 취약

### Cloudflare Tunnel 방식

1. Cloudflare Tunnel은 서버가 먼저 Cloudflare로 Outbound 연결을 생성
    
    ```
    [항상 유지되는 연결]
    홈서버(cloudflared) → Cloudflare
    ```
    
    - `cloudflared` 데몬은 Cloudflare Edge와 여러 개의 TLS 연결을 지속적으로 유지
    - 이 시점에는 외부 요청이 없어도 연결은 계속 생존
    - 즉, 요청이 올 때마다 연결을 여는 구조가 아니다
2. 외부 사용자가 서비스에 접근하면 요청은 먼저 Cloudflare에 도달
    
    ```
    외부 클라이언트 → Cloudflare
    ```
    
    이 단계에서 Cloudflare는 아래 접근 제어를 수행
    
    - Cloudflare Access 기반 인증 (로그인, OTP, SSO, 서비스 토큰 등)
    - IP, 국가, 정책 기반 접근 제어
    - 인증 실패 시 요청 차단
    - 인증과 접근 제어는 홈서버와 분리된 상태에서 Cloudflare가 전담
3. 인증 이후 실 연결
    
    ```
    Cloudflare → 기존 Tunnel → 홈서버
    ```
    
    - 인증을 통과한 요청만 이미 열려 있는 터널 연결을 통해 홈서버로 전달됨

## 2. Cloudflare Access

Cloudflare tunnel만 사용한다면, 최초 목적 중 공인 IP 숨기기만 가능하다.

PVE 콘솔에 대한 접근 자체를 막기 위해선 Cloudflare Access를 사용한 인증 절차 도입이 필요하다.

### Cloudflare Access 도입 후의 구조

- 위 Cloudflare Tunnel의 2번째 단계에서 Cloudflare Access 추가
    
    ```jsx
    외부 클라이언트 → Cloudflare Access → Cloudflare
    ```
    

### 정리

- Cloudflare Tunnel은 요청을 역방향으로 전달하는 Outbound 기반 리버스 프록시 구조
- 홈서버 입장에서는 외부 클라이언트와 직접 통신하지 않는다
- Cloudflare와 유지 중인 Outbound 연결만 사용하며, 외부에서 서버로 들어오는 Inbound 포트는 존재하지 않는다
- 따라서 Zero-Inbound로 외부에서 서비스 접근이 가능
- 내가 외부에서 접근해야 하지만 공개적으로 오픈하면 안 될 경우 Cloudflare Access 사용
- 실제 작업 과정은 다음 글 확인
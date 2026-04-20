---
title: "Cloudflare Tunnel을 사용한 홈서버 Promox PVE 콘솔 접근 설정"
date: 2026-02-06T09:00:00+09:00
categories: ["Home Server", "Cloudflare"]
tags: ["home server", "cloudflare", "cloudflare tunnel", "cloudflare access", "promox"]
---


### 목표

- 현재 외부에서 홈서버 Promox PVE 콘솔에 접근 시 공인 IP로 직접 접근하는 중
- PVE 접속용 도메인 설정
- 외부에서 도메인 접근 시 Cloudflare 인증을 요구하도록 설정하고, 인증 성공 시에만 PVE 콘솔에 접근하도록 설정

### 0. 사전 준비

- Cloudflare에 등록할 도메인 준비
- 도메인 네임서버를 Cloudflare로 설정
    - Cloudflare 회원가입 및 로그인
    - 로그인 시 나오는 네임서버 도메인 2개를 도메인 등록기관의 네임서버 설정에 등록
    - 초록색 체크 표시 뜨면 OK!
        
        ![CloudflareTunnelSetting01.png](CloudflareTunnelSetting01.png)
        

### 1. Cloudflared 설치 및 연동

![CloudflareTunnelSetting02.png](CloudflareTunnelSetting02.png)

- `Networks` - `Connectors` - `Add a tunnel`

![CloudflareTunnelSetting03.png](CloudflareTunnelSetting03.png)

- `Cloudflared` 선택

![CloudflareTunnelSetting04.png](CloudflareTunnelSetting04.png)

- tunnel 이름 기입

![CloudflareTunnelSetting05.png](CloudflareTunnelSetting05.png)

- OS 확인 및 Cloudflared 설치
    - Promox는 Debian 리눅스 64bit
        - `cat /etc/os-release` + `uname -m` 명령어로 확인 가능
    - `root` 사용자로 진행
        - Promox는 sudo가 설치되어 있지 않음
    
    ```jsx
    # Add cloudflare gpg key
    mkdir -p --mode=0755 /usr/share/keyrings
    curl -fsSL https://pkg.cloudflare.com/cloudflare-public-v2.gpg | tee /usr/share/keyrings/cloudflare-public-v2.gpg >/dev/null
    
    # Add this repo to your apt repositories
    echo 'deb [signed-by=/usr/share/keyrings/cloudflare-public-v2.gpg] https://pkg.cloudflare.com/cloudflared any main' | tee /etc/apt/sources.list.d/cloudflared.list
    
    # install cloudflared
    apt-get update && apt-get install cloudflared
    ```
    
    - apt-get update 시 `enterprise.proxmox.com 패키지 Unauthorized` 오류가 나는 경우
        - Promox의 기본 패키지는 유료 구독형이기 때문에 비활성화 필요
        
        ```jsx
        cd /etc/apt/sources.list.d
        
        mv -v pve-enterprise.sources pve-enterprise.sources.disabled
        mv -v ceph.sources ceph.sources.disabled
        ```
        
        - 일반 패키지 리포 추가
        
        ```jsx
        cat > /etc/apt/sources.list.d/pve-no-subscription.sources <<'EOF'
        Types: deb
        URIs: http://download.proxmox.com/debian/pve
        Suites: trixie
        Components: pve-no-subscription
        Signed-By: /usr/share/keyrings/proxmox-archive-keyring.gpg
        EOF
        ```
        
    - 하단의 토큰 값 유출되지 않도록 유의
        - 서버가 시작될때마다 Tunnel이 자동으로 기동되게 할 경우
            
            ```jsx
            cloudflared service install <토큰>
            ```
            
        - 현재 터미널에서만 Tunnel 연결할 경우
            
            ```jsx
            cloudflared tunnel run --token <토큰>
            ```
            

### 2. Cloudflare Routes 설정

- 연동한 tunnel을 통해서 연결할 도메인 경로 및 라우팅 경로 설정
    - Hostname: 보유 도메인 기입
    - Type: HTTPS
    - URL: localhost:8006
    
    ![CloudflareTunnelSetting06.png](CloudflareTunnelSetting06.png)
    
- Promox 서버에서 https://localhost:8006로 접근 시 TLS 해제 필요
    
    ![CloudflareTunnelSetting07.png](CloudflareTunnelSetting07.png)
    
- 여기까지 설정 완료 시 기입한 서브도메인으로 PVE 접근 가능

### 3. Cloudflare Access 설정

- 서브도메인에 대한 인증 설정

![CloudflareTunnelSetting08.png](CloudflareTunnelSetting08.png)

![CloudflareTunnelSetting09.png](CloudflareTunnelSetting09.png)

- Self-hosted 설명에 나와 있듯이 Policy 먼저 생성
    
    ![CloudflareTunnelSetting10.png](CloudflareTunnelSetting10.png)
    
    - PVE 접속 도메인 접근 시 Email 인증하도록 Policy 설정

![CloudflareTunnelSetting11.png](CloudflareTunnelSetting11.png)

- 설정한 Policy 할당

![CloudflareTunnelSetting12.png](CloudflareTunnelSetting12.png)

- PVE 접속 도메인 접근 시 이메일 인증 요구하는 화면 출력 확인
- 설정 완료!

### 3. Promox 인바운드 접근 차단

- Cloudflare Tunnel이 뚫렸기 때문에, 기존에 anyopen되어 있던 Promox 공인 IP에 대한 접근을 전부 차단
- Zero-Inbound 설정
- PVE 콘솔 → Datacenter → Firewall → Options → Firewall 활성화
    
    ![CloudflareTunnelSetting13.png](CloudflareTunnelSetting13.png)
    
- 설정 즉시 공인 IP를 통한 PVE 콘솔 차단 완료!

### 결과

- 홈서버로서의 취약한 보안 조치 완료
    - 홈 공인IP 숨기기
    - 관리자용 페이지 접근 시 이메일 인증 적용
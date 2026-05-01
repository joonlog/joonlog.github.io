---
title: "SSH Config 옵션 정리"
date: 2026-05-01T09:00:00+09:00
categories: ["Linux", "SSH"]
tags: ["linux", "ssh", "ssh config", "proxyjump", "proxycommand"]
---


> `~/.ssh/config` SSH 접속 편의성
> 

매번 `ssh -i key user@host -p port` 처럼 옵션을 직접 입력하는 대신, `~/.ssh/config`에 설정을 저장해두면 `ssh <Host명>` 만으로 접속 가능하다.

- Windows: `C:\\Users\\<계정명>\\.ssh\\config`
- Linux/Mac: `~/.ssh/config`

## 1. 기본 구조

```
Host <별칭>
    옵션 값
```

- `Host`에 지정한 별칭으로 `ssh <별칭>` 접속
- 와일드카드 사용 가능 (`Host *`, `Host node*`)

## 2. 주요 옵션

### HostName

- 실제 접속 대상 주소 (IP 또는 도메인)
- `Host`에 별칭을 쓰고 실제 주소는 여기에 기입

```
Host mgt
    HostName ssh-mgt.joonlog.store
```

### User

- 접속할 계정명
- 설정하지 않으면 현재 로컬 계정으로 접속 시도

```
    User root
```

### IdentityFile

- 사용할 개인키 파일 경로
- 설정하지 않으면 `~/.ssh/id_rsa` 등 기본 키 파일을 순서대로 시도

```
    IdentityFile ~/.ssh/pve
```

### Port

- 접속할 SSH 포트
- 기본값은 22이며 커스텀 포트 사용 시 지정

```
    Port 2222
```

### ServerAliveInterval / ServerAliveCountMax

- 연결 유지 설정
- `ServerAliveInterval`: 지정한 초마다 서버에 alive 패킷 전송
- `ServerAliveCountMax`: 응답 없을 때 최대 재시도 횟수, 초과 시 연결 종료
- idle 상태에서 연결이 끊기는 것을 방지

```
    ServerAliveInterval 60
    ServerAliveCountMax 3
```

### ProxyCommand

- SSH 연결 전에 실행할 명령어로 연결을 중계
- Cloudflare Tunnel SSH 접속 시 사용

```
    ProxyCommand cloudflared access ssh --hostname %h
```

- `%h`는 `HostName` 값으로 자동 치환

### ProxyJump

- 점프 호스트를 거쳐서 접속할 때 사용
- VSCode Remote SSH나 내부망 서버 접속 시 유용
- `ProxyJump <Host별칭>` 형식으로 이미 config에 정의된 Host 재사용 가능

```
Host node1
    HostName <IP>
    User root
    IdentityFile ~/.ssh/pve
    ProxyJump mgt
```

- 위 설정 시 `ssh node1` 하면 mgt를 경유해서 node1에 접속
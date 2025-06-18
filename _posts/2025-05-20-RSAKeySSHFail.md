---
title : RSA Key를 사용한 SSH 불가 TroubleShooting
date : 2025-05-20 09:00:00 +09:00
categories : [Linux, SSH]
tags : [linux, ssh, rsa, key, esdsa, troubleshooting] #소문자만 가능
---

## 문제 요약

- Ubuntu 24.04, Amazon Linux 2023 등 최신 리눅스에서 기존 RSA 키로 접속 불가
- 특히 SecureCRT는 `ssh-rsa` 알고리즘만 사용해서 오류 발생
- Putty는 `rsa-sha2-256` 지원 → 같은 키로도 접속 가능

## 원인

- 최신 OpenSSH는 `ssh-rsa` (SHA-1 기반) 비활성화 → 보안상 위험
- Amazon Linux 2023은 `/etc/crypto-policies/back-ends/opensshserver.config`를 통해 강제로 막음
- Ubuntu는 기본 정책에선 막히진 않지만, OpenSSH 자체가 기본으로 제외시킴

## 해결 방법 (서버 측 설정 변경)

```bash
vim /etc/ssh/sshd_config
```

- 아래 두 줄을 추가
    - **Amazon Linux 2023**: 반드시 `Include`보다 위쪽에
    - **Ubuntu 24.04**: 위치는 무관 (보통 맨 아래 추가)

```bash
HostkeyAlgorithms +ssh-rsa
PubkeyAcceptedAlgorithms +ssh-rsa
```

```bash
# 설정 확인 및 sshd 재시작
sshd -t
systemctl restart sshd
```

## 주의사항

- `ssh-rsa`는 **보안 취약점 있음** (SHA-1 기반)
- 이 설정은 **임시 방편**일 뿐, 장기적으론 권장되지 않음

## 대안: 새 키 생성 시 최신 알고리즘 사용

```bash

ssh-keygen -t ed25519 -C "your@email.com"
```

> 추후 key 생성 시 RSA 보단 ESDSA 같은 알고리즘 사용 필요
>
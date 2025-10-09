---
title : Teleport 18 이미지 빌드 및 배포
date : 2025-10-09 09:00:00 +09:00
categories : [Go, Teleport]
tags : [go, teleport, teleport 18, dockerfile, docker compose]  #소문자만 가능
---

## Teleport 18 이미지 빌드

### 도커 이미지 빌드

- Dockerfile
    - 빌드한 바이너리 파일의 teleport/tctl/tsh/tbot을 복사해서 이미지 빌드

```bash
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y \
    ca-certificates dumb-init libfido2-1 && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY build/teleport /usr/local/bin/teleport
COPY build/tctl /usr/local/bin/tctl
COPY build/tsh /usr/local/bin/tsh
COPY build/tbot /usr/local/bin/tbot

ENTRYPOINT ["/usr/bin/dumb-init", "teleport", "start", "-c", "/etc/teleport/teleport.yaml"]
```

- 이미지 빌드 및 도커허브 푸시

```bash
docker build -t teleport18-custom .
docker login -u <username>
docker tag teleport18-custom <username>/teleport18:v1.0.0
docker push <username>/teleport18:v1.0.0
```

## Teleport 18 배포

> Auth+Proxy 분리 구축
> 

### 디렉토리 구조

```
# tree
.
├── auth
│   ├── config
│   │   └── teleport.yaml
│   └── data
├── docker-compose.yaml
└── proxy
    ├── config
    │   └── teleport.yaml
    └── data
```

---

### 1. docker-compose.yaml

```yaml
version: '3.8'

services:
  auth:
    image: <username>/teleport18:v1.0.0
    container_name: teleport-auth
    hostname: teleport-auth
    restart: unless-stopped
    volumes:
      - ./auth/config:/etc/teleport
      - ./auth/data:/var/lib/teleport
    entrypoint: ["/usr/local/bin/teleport"]
    command: ["start", "--config=/etc/teleport/teleport.yaml"]
    expose:
      - "3025"
    networks:
      - teleport-net

  proxy:
    image: <username>/teleport18:v1.0.0
    container_name: teleport-proxy
    hostname: teleport-proxy
    restart: unless-stopped
    ports:
      - "3080:3080"      # Web UI & public HTTPS
      - "3023:3023"    # Proxy SSH
      - "3024:3024"    # Reverse tunnel
      - "3026:3026"    # Kube reverse tunnel
    volumes:
      - ./proxy/config:/etc/teleport
      - ./proxy/data:/var/lib/teleport
      - /etc/letsencrypt:/etc/letsencrypt:ro
    environment:
      - TELEPORT_CDN_BASE_URL=https://cdn.cloud.gravitational.io
    entrypoint: ["/usr/local/bin/teleport"]
    command: [
      "start",
      "--config=/etc/teleport/teleport.yaml",
      "--roles=proxy",
      "--token=teleportclustertoken",
      "--auth-server=teleport-auth:3025"
    ]
    networks:
      - teleport-net
    depends_on:
      - auth

networks:
  teleport-net:
    driver: bridge
```

---

### 2. auth/teleport.yaml

```yaml
teleport:
  nodename: teleport-auth
  data_dir: /var/lib/teleport
  log:
    output: stderr
    severity: INFO

auth_service:
  enabled: yes
  cluster_name: cluster
  listen_addr: 0.0.0.0:3025
  tokens:
  - "proxy:teleportclustertoken"

proxy_service:
  enabled: no

ssh_service:
  enabled: no
```

---

### 3. proxy/teleport.yaml

```yaml
teleport:
  nodename: teleport-proxy
  data_dir: /var/lib/teleport
  log:
    output: stderr
    severity: INFO

  auth_servers:
    - teleport-auth:3025

auth_service:
  enabled: no

proxy_service:
  enabled: yes
  web_listen_addr: 0.0.0.0:3080
  tunnel_listen_addr: 0.0.0.0:3024
  kube_listen_addr: 0.0.0.0:3026
  public_addr: <도메인>:3080
  https_keypairs:
    - cert_file: /etc/letsencrypt/live/<도메인>/fullchain.pem
      key_file: /etc/letsencrypt/live/<도메인>/privkey.pem

ssh_service:
  enabled: no
```
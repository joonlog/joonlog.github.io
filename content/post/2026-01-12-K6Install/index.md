---
title: "K6 설치 방법"
date: 2026-01-12T09:00:00+09:00
categories: ["Testing", "k6"]
tags: ["k6", "performance test", "grafana k6", "k6 install"]
---


### K6

- VU(가상 사용자) 기반의 HTTP 부하 생성 툴
- 스크립트에 정의된 시나리오를 수행

## K6 설치

K6 GitHub 최신 Release에서 실행 파일을 다운로드

https://github.com/grafana/k6/releases

### Window

- winget을 사용해 설치 가능
    
    ```powershell
    winget install k6
    k6 version
    ```
    
- GitHub Release로도 설치 가능

### Linux

- GitHub Release에 있는 OS에 맞는 파일을 설치
    - ex) k6-v1.5.0-linux-amd64.tar.gz
    
    ```jsx
    tar -xvzf k6-v1.5.0-linux-amd64.tar.gz
    mv k6-v1.5.0-linux-amd64/k6 /usr/local/bin/k6
    k6 version
    ```
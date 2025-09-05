---
title : Containerd 설치 및 환경 세팅하기
date : 2025-09-05 09:00:00 +09:00
categories : [Container, Containerd]
tags : [container, containerd, nerdctl, buildkit] #소문자만 가능
---

> containerd만 단독으로 설치 했을 경우 CRI의 기능만 제공하며, docker처럼 기능을 사용하기 위해선 플러그인 설치가 필요
> 
- https://containerd.io/downloads/
- https://medium.com/@DannielWhatever/using-containerd-without-docker-9d08332781b4
1. containerd 설치
    
    ```bash
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo dnf remove $pkg -y; done
    dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    dnf install -y containerd.io
    containerd config default | sudo tee /etc/containerd/config.toml
    ```
    
2. containerd 디렉토리 설정(선택)
    
    ```bash
    ### containerd root 디렉토리 외에 data 디렉토리에서 동작하도록 설정
    sed -i 's|^root = .*|root = "/data/container/containerd"|' /etc/containerd/config.toml
    sed -i 's|^state = .*|state = "/data/container/containerd-state"|' /etc/containerd/config.toml
    mkdir -p /data/container/containerd
    mkdir -p /data/container/containerd-state
    mkdir -p /data/container/nerdctl
    echo 'export NERDCTL_DATA_ROOT=/data/container/nerdctl' >> ~/.bashrc
    source ~/.bashrc
    ```
    
3. nerdctl 설치
    - 도커 명령어처럼 containerd를 관리하기 위한 툴
        - ex) nerdctl ps -a, nerdctl images
    
    ```bash
    # nerdctl 설치
    wget https://github.com/containerd/nerdctl/releases/download/v0.17.0/nerdctl-0.17.0-linux-amd64.tar.gz
    tar Cxzvvf /usr/bin nerdctl-0.17.0-linux-amd64.tar.gz
    echo "source <(nerdctl completion bash)" >> ~/.bashrc
    ```
    
4. CNI 설치
    - 컨테이너 네트워크를 관리하기 위한 툴
        - 컨테이너의 ip 등의 네트워크를 지정해서 설정 가능
    
    ```bash
    # CNI 설치
    wget https://github.com/containernetworking/plugins/releases/download/v1.0.1/cni-plugins-linux-amd64-v1.0.1.tgz
    mkdir -p /opt/cni/bin
    tar Cxzvvf /opt/cni/bin cni-plugins-linux-amd64-v1.0.1.tgz
    ```
    
5. buildkit 설치
    - containerd에서 image를 빌드하기 위한 툴
    
    ```bash
    wget https://github.com/moby/buildkit/releases/download/v0.10.0-rc1/buildkit-v0.10.0-rc1.linux-amd64.tar.gz 
    tar Cxzvvf /usr/bin buildkit-v0.10.0-rc1.linux-amd64.tar.gz 
    nohup /usr/bin/builkitd < /dev/null > /var/log/buildkitd 2>&1 &
    ```
    
6. containerd 실행
    
    ```bash
    systemctl enable --now containerd
    ```
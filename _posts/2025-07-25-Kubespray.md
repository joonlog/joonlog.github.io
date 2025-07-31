---
title : Kubespray를 사용한 Kubernetes 클러스터 배포
date : 2025-07-25 09:00:00 +09:00
categories : [Kubernetes, Deploy]
tags : [kubernetes, k8s, kubespray, docker]  #소문자만 가능
---

https://github.com/kubernetes-sigs/kubespray

- Kubespray를 사용해 설치
    - kubespray 2.27 버전을 사용하여 Kubernetes 1.31버전 클러스터 구성
- 최소 4대가 필요
    - 배포 서버 1EA
        - 도커로 Kubespray를 배포하고, Kubernetes CRI로 containerd를 사용할 거기 때문에 Master Node에서 배포 불가해서 별도의 배포 서버 필요
    - Master Node 1EA
    - Worker Node 2EA

## 사전작업

- 모든 Master, Worker Node에 작업
    
    ```bash
    # vim 설정
    cat << EOF >> ~/.bashrc
    alias vi='/usr/bin/vim'
    EOF
    
    cat << EOF >> ~/.vimrc
    syntax on
    autocmd FileType yaml setlocal ai nu ts=2 sw=2 et
    autocmd Filetype python setlocal ai nu ts=2 sw=2 et
    EOF
    
    . ~/.bashrc
    
    # DNS 서버가 없으므로 /etc/hosts 파일 사용
    cat << EOF >> /etc/hosts
    172.27.1.9    master.k8s.com  master
    172.27.0.182  node1.k8s.com   node1
    172.27.1.35   node2.k8s.com   node2
    EOF
    
    # cri 전부 제거
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo dnf remove $pkg -y; done
    ```
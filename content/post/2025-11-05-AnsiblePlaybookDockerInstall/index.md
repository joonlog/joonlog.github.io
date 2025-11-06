---
title: "Ansible Playbook으로 Docker 설치"
date: 2025-11-04T09:00:00+09:00
categories: ["Ansilbe", "Playbook"]
tags: ["ansible", "playbook", "docker", "redhat"]
---


> RedHat 계열 리눅스에서 Ansible로 도커 설치하는 방법
> 

## Ansible 설정

- install_docker.yaml
    - 10대 이상의 리눅스에 도커를 설치하는 playbook
    
    ```bash
    - name: Install Docker CE on web/was hosts
      hosts: webwas
      become: true
      gather_facts: yes
    
      vars:
        conflict_pkgs:
          - docker.io        
          - docker-doc
          - docker-compose
          - docker-compose-v2
          - podman-docker
          - containerd
          - runc
    
        docker_pkgs:
          - docker-ce
          - docker-ce-cli
          - containerd.io
          - docker-buildx-plugin
          - docker-compose-plugin
    
      tasks:
        - name: 01. Remove conflicting packages (ignore if not present)
          package:
            name: "{{ conflict_pkgs }}"
            state: absent
    
        - name: 02. Install prerequisites
          package:
            name: yum-utils
            state: present
          when: ansible_os_family == "RedHat"
    
        - name: 03. Add Docker CE yum repository (official)
          yum_repository:
            name: docker-ce-stable
            description: Docker CE Stable - $basearch
            baseurl: "https://download.docker.com/linux/centos/$releasever/$basearch/stable"
            gpgcheck: yes
            gpgkey: https://download.docker.com/linux/centos/gpg
            enabled: yes
          when: ansible_os_family == "RedHat"
    
        - name: 04. Install Docker CE packages
          package:
            name: "{{ docker_pkgs }}"
            state: present
    
        - name: 05. Enable and start Docker
          service:
            name: docker
            state: started
            enabled: yes
    ```
    
    - task01: 기존 충돌 패키지 제거
    - task02: 필수 패키지 설치
    - task03: Docker CE 패키지 리포지토리 추가
        - OS 버전에 맞는 리포지토리 자동 설정
    - task04: 도커 설치
    - task05: 도커 서비스 시작 및 자동 실행
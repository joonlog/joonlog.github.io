---
title: "Linux 환경에서 Jenkins 설치"
date: 2025-10-23T09:00:00+09:00
categories: ["Linux", "Jenkins"]
tags: ["linux", "jenkins", "rocky 8.10", "java21"]
---


> 로컬 리눅스 환경에서 Jenkins 설치하는 방법
> 

### 환경

- OS: Rocky 8.10
- java: OpenJDK 21

### 설치

1. 패키지 설치
    
    ```bash
    wget -O /etc/yum.repos.d/jenkins.repo https://pkg.jenkins.io/redhat-stable/jenkins.repo
    rpm --import https://pkg.jenkins.io/redhat-stable/jenkins.io-2023.key
    
    yum install fontconfig java-21-openjdk
    ```
    
    - (선택) OpenJDK 21이 아닌 다른 버전으로 선택되어 있을 시 변경 필요
    
    ```bash
    java --version
    # java 21 선택
    alternatives --config java
    java --version
    ```
    
2. jenkins 설치
    
    ```bash
    yum install jenkins
    systemctl daemon-reload
    systemctl enable --now jenkins
    ```
    
3. http://<ip>:8080으로 콘솔 접근
4. Install wizard 진행
    - jenkins는 최초 접근 시 unlock 되어 있음
        - 비밀번호 확인 필요
        
        ```bash
        cat /var/lib/jenkins/secrets/initialAdminPassword
        ```
        
        ![LinuxJenkinsInstall01.png](LinuxJenkinsInstall01.png)
        
    - 이후 자동 설치 진행

### 참고

- Jenkins Install 공식 문서: 
https://www.jenkins.io/doc/book/installing/linux/
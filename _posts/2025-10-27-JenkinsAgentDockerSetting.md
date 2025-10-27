---
title : Jenkins Agent를 Docker로 동작시키기 위한 방법
date : 2025-10-27 09:00:00 +09:00
categories : [Container, Jenkins]
tags : [docker, container, jenkins, dockerfile, docker compose]  #소문자만 가능
---

> Jenkins Agent가 도커로 사용되도록 설정
Jenkins Server-Agent 분리
> 
1. Built-In-Node excutor 설정
    - Jenkins Agent를 도커에서 동적으로 실행할거라 `built-in-node`의 `Number of executors`값을 0으로 설정
        
        ![JenkinsAgentDocker01.png](/assets/img/container/JenkinsAgentDocker01.png)
        
2. Agent Node 추가
    - Manage Jenkins - Nodes - System Configuration - New Node
    - Node 설정
        - `Number of excutors`: Jenkins Agent가 최대로 실행 가능한 파이프라인 수
        - `Remote root directory`: Jenkins Agent 컨테이너 내부의 루트 디렉토리
        
        ![JenkinsAgentDocker02.png](/assets/img/container/JenkinsAgentDocker02.png)
        
3. Agent 설치 스크립트 확인
    - 이후에 컨테이너 기동에 필요한 secret 값이 포함되어 있으니 확인 필수
    
    ![JenkinsAgentDocker03.png](/assets/img/container/JenkinsAgentDocker03.png)
    
4. Dockerfile
    - 파이프라인용 에이전트이므로 최소로 컨테이너 기동
    
    ```bash
    FROM jenkins/inbound-agent:latest-jdk21
    USER root
    RUN apt-get update \
     && apt-get install -y docker.io ca-certificates openssh-client \
     && rm -rf /var/lib/apt/lists/*
    USER jenkins
    ```
    
5. docker-compose.yaml
    - 변수들은 .env 파일을 통해 별도로 설정
        - `JENKINS_URL`:  도커 네트워크 인터페이스 docker0의 ip인 `172.17.0.1`를 기입
        - `JENKINS_SECRET`:  위 콘솔에 나온 secret 값을 기입
        - `JENKINS_AGENT_NAME`: 위에서 설정한 이름
        - `JENKINS_AGENT_WORKDIR`: 위에서 설정한 Jenkins Agent가 컨테이너 내부에서 사용할 경로
        - `AGENT_WORKDIR`: Jenkins Agent 컨테이너가 호스트에서 사용하는 실제 경로
    
    ```bash
    version: "3.8"
    
    services:
      agent:
        build:
          context: .
          dockerfile: Dockerfile
        container_name: jenkins-agent-docker
        restart: unless-stopped
        init: true
        environment:
          JENKINS_URL: "${JENKINS_URL}"
          JENKINS_SECRET: "${JENKINS_SECRET}"
          JENKINS_AGENT_NAME: "${JENKINS_AGENT_NAME}"
          JENKINS_AGENT_WORKDIR: "${JENKINS_AGENT_WORKDIR:-/home/jenkins/agent}"
          DOCKER_BUILDKIT: "1"
          JENKINS_WEB_SOCKET: "true"
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock
          - ${AGENT_WORKDIR:-/var/lib/jenkins_agent}:${JENKINS_AGENT_WORKDIR:-/home/jenkins/agent}
        group_add:
          - "${DOCKER_GID}"
    ```
    
6. 컨테이너 권한 설정
    - 호스트의 도커 소켓을 사용하는 DooD 방식이기 때문에 컨테이너에서도 호스트의 도커 GID 동기화 필요
        - `/var/lib/jenkins_agent`: 호스트에 컨테이너 데이터가 저장될 경로.
            - 컨테이너 내부에서도 기록해야하니 기본 사용자인 1000으로 설정
    
    ```bash
    echo "DOCKER_GID=$(stat -c %g /var/run/docker.sock)" >> .env
    mkdir -p /var/lib/jenkins_agent && chown 1000:1000 /var/lib/jenkins_agent
    ```
    
7. Jenkins Agent 빌드
    
    ```bash
    docker compose up -d --build
    ```
    
8. Node 연결 확인
    
    ![JenkinsAgentDocker04.png](/assets/img/container/JenkinsAgentDocker04.png)
    
9. 파이프라인 테스트
    - hello world 테스트 성공!
    
    ![JenkinsAgentDocker05.png](/assets/img/container/JenkinsAgentDocker05.png)
    

### 참고

- Jenkins Agent 노드 추가:
https://github.com/jenkinsci/docker-agent/blob/master/README_inbound-agent.md
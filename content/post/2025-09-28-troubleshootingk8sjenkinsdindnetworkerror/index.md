---
title: "K8S에 구축한 Jenkins CI 파이프라인의 Dind 네트워크 에러 TroubleShooting"
date: 2025-09-27T09:00:00+09:00
categories: ["Kubernetes", "Jenkins"]
tags: ["Kubernetes", "k8s", "self managed k8s", "troubleshooting", "jenkins", "dind"]
---


> K8S 환경의 Jenkins에서 docker build가 Jenkins 파이프라인 파드의 dind 컨테이너에서 실행되기 때문에 발생한 이슈
> 

### 환경

- Jenkins 서버를 K8S에서 구축하고, 파이프라인을 K8S 파드에서 동작하도록 설정한 환경
    - Jenkinsfile에서 K8S 파드 템플릿으로 파이프라인 실행하게 명시

### 문제 상황

- Jenkinsfile의 `docker build` 단계에서 파이프라인 빌드 실패
    - Dockerfile의 composer install 레이어에서 `Could not fetch https://api.github.com/repos/` 에러 발생
    - 외부에서 데이터를 받아오는 명령어에서 전부 timeout 에러 발생

### 에러 분석

- Jenkins 파이프라인이 K8S 파드의 `docker` 컨테이너에서 실행되고, docker build 작업은 `dind` 컨테이너에서 실행됨
    - `dind` 컨테이너 입장에서 네트워크가 격리되어 있어 GitHub API 서버로의 아웃바운드 연결 불가
- 따라서 `composer install` 와 같은 명령어 실행 시 연결 실패 에러 발생

### 해결 방법

- Jenkinsfile의 k8s 파드 템플릿에 `hostNetwork: true`를 설정해서 파드의 컨테이너가 호스트 노드(현재 워커 노드)의 네트워크 네임스페이스를 사용할 수 있도록 설정
    - 파드 템플릿이 아닌 `docker build --network=host`로 명령어 마다 옵션을 명시해서 해결도 가능하지만, 파드 레벨에서 해결하는게 관리에 용이

```bash
kubernetes {
            yaml """
            apiVersion: v1
            kind: Pod
            spec:
              serviceAccountName: jenkins
              hostNetwork: true
              hostAliases:
              - ip: "172.27.1.9"
                hostnames:
                - "registry.<도메인>"
                - "gitlab.<도메인>"
              containers:
              - name: docker
                image: docker
                command:
                - sleep
                args:
                - 99d
                tty: true
                volumeMounts:
                - name: workspace-volume
                  mountPath: /home/jenkins/agent
                - name: docker-socket
                  mountPath: /var/run
              - name: docker-daemon
                image: docker:dind
                securityContext:
                  privileged: true
                env:
                - name: DOCKER_TLS_CERTDIR
                  value: ""
                args:
                - --insecure-registry=registry.<도메인>:8080
                volumeMounts:
                - name: docker-socket
                  mountPath: /var/run
              volumes:
              - name: workspace-volume
                emptyDir: {}
              - name: docker-socket
                emptyDir: {}
            """
        }
```
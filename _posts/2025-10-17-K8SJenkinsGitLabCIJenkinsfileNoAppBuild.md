---
title : K8S에 구축된 Jenkins-GitLab CI Jenkinsfile - 애플리케이션 빌드 필요 없는 경우
date : 2025-10-17 09:00:00 +09:00
categories : [Kubernetes, Jenkins]
tags : [Kubernetes, k8s, self managed k8s, jenkins, gitlab, jenkins pipeline, ci, jenkinsfile] #소문자만 가능
---

> Kubernetes 환경에서 Jenkins 파이프라인으로 CI를 수행하도록 구성한 Jenkinsfile 정리
> 

> Jenkins와 GitLab이 같은 K8S에 있는 환경인 가정하에 유효
> 

> PHP처럼 별도의 애플리케이션 빌드 과정이 필요 없는 경우에 유효
> 

### 1. Kubernetes Pod 템플릿

- Jenkins 에이전트를 파드 단위로 생성하는 파드 템플릿 정의
- docker, dind 컨테이너로 구성
    - `docker-daemon` 컨테이너는 Docker in Docker(DIND)를 사용해서 도커 데몬 실행
- `workspace-volume`와 `docker-socket` 볼륨 공유
- `hostNetwork: true` 설정으로 파드의 컨테이너가 호스트 노드(현재 워커 노드)의 네트워크 네임스페이스를 사용할 수 있도록 설정
    - `hostAliases` 설정으로 사설 IP로 GitLab 서버에 접근할 수 있도록 설정
    
    ```bash
        agent {
            kubernetes {
                yaml """
                apiVersion: v1
                kind: Pod
                spec:
                  serviceAccountName: jenkins
                  hostNetwork: true
                  hostAliases:
                  - ip: "<마스터 서버 사설 IP>"
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
        }
    ```
    

### 2. Environment

- Jenkinsfile에서 사용될 전역 환경 변수 정의
    - `GITLAB_REGISTRY`: GitLab Registry 주소
    - `GITLAB_TOKEN`: Jenkins Credentials에 저장된 GitLab 토큰
    
    ```bash
        environment {
            APP_NAME = "<앱 이름>"
            RELEASE = "1.0.0"
            GITLAB_REGISTRY = "registry.<도메인>:8080"
            GITLAB_TOKEN = credentials('<credentials 이름>')
            IMAGE_NAME = "${GITLAB_REGISTRY}/root/${APP_NAME}"
            IMAGE_TAG = "${RELEASE}-${BUILD_NUMBER}"
        }
    ```
    

### 3. Stages

- Stage 1
    - SCM에서 소스코드 체크아웃
    - SCM의 url은 K8S Service의 FQDN으로 사용
    
    ```bash
            stage("Checkout from SCM") {
                steps {
                    script {
                        sh 'pwd'
                        sh 'ls -al'
                        git credentialsId: '<credentials 이름>',
                            url: 'http://gitlab-webservice-default.gitlab.svc.cluster.local:8181/root/<리포지토리명>.git',
                            branch: 'main'
                        sh 'ls -al'
                    }
                }
            }
    ```
    
- Stage 2
    - Docker 이미지 빌드
        - DinD 컨테이너에서 이미지 빌드 및 태그 추가
    
    ```bash
            stage('Build Docker Image') {
                steps {
                    container('docker') {
                        script {
                            sh 'pwd'
                            sh 'ls -al'
                            sh 'docker build --no-cache -t ${IMAGE_NAME}:${IMAGE_TAG} .'
                            sh 'docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:latest'
                        }
                    }
                }
            }
    ```
    
- Stage 3
    - GitLab Registry 이미지 푸시
        - GitLab Container Registry에 로그인 후 이미지 푸시
        - `GITLAB_TOKEN` 환경변수에서 주입받은 credential로 `GITLAB_TOKEN_PSW`, `GITLAB_TOKEN_USR` 변수 사용 가능
    
    ```bash
    stage('Push to GitLab Registry') {
                steps {
                    container('docker') {
                        script {
                            sh 'echo $GITLAB_TOKEN_PSW | docker login ${GITLAB_REGISTRY} -u $GITLAB_TOKEN_USR --password-stdin'
                            sh 'docker push ${IMAGE_NAME}:${IMAGE_TAG}' // Docker 이미지 푸시
                            sh 'docker push ${IMAGE_NAME}:latest' // Docker 이미지 푸시(latest 태그)
                        }
                    }
                }
            }
    ```
    

### 4. 정리

- 파이프라인 종료 전 도커 리소스 정리 및 로그아웃
    
    ```bash
        post {
            always {
                // 작업 완료 후 정리
                container('docker') {
                    script {
                        sh 'docker logout ${GITLAB_REGISTRY} || true'
                        sh 'docker system prune -f || true'
                    }
                }
            }
            success {
                echo "Pipeline succeeded! Image ${IMAGE_NAME}:${IMAGE_TAG} pushed to GitLab Registry"
            }
            failure {
                echo "Pipeline failed!"
            }
        }
    ```
    

### 전체 Jenkinsfile

```bash
pipeline {
    agent {
        kubernetes {
            yaml """
            apiVersion: v1
            kind: Pod
            spec:
              serviceAccountName: jenkins
              hostNetwork: true
              hostAliases:
              - ip: "<마스터 서버 사설 IP>"
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
    }

    environment {
        APP_NAME = "<앱 이름>"
        RELEASE = "1.0.0" // 릴리즈 버전
        GITLAB_REGISTRY = "registry.<도메인>:8080" // GitLab Registry 주소
        GITLAB_TOKEN = credentials('<credentials 이름>') // Jenkins Credentials에 저장된 GitLab 토큰
        IMAGE_NAME = "${GITLAB_REGISTRY}/root/${APP_NAME}" // 이미지 이름
        IMAGE_TAG = "${RELEASE}-${BUILD_NUMBER}" // 이미지 태그
    }

    stages {
        stage("Checkout from SCM") { // 소스 코드 관리(SCM)에서 체크아웃 단계
            steps {
                script {
                    sh 'pwd'
                    sh 'ls -al'
                    git credentialsId: '<credentials 이름>', // GitLab 자격증명 사용
                        url: 'http://gitlab-webservice-default.gitlab.svc.cluster.local:8181/root/<리포지토리명>.git', // GitLab 리포지토리 URL
                        branch: 'main' // 브랜치 이름
                    sh 'ls -al'
                }
            }
        }

        stage('Build Docker Image') { // Docker 이미지 빌드 단계
            steps {
                container('docker') {
                    script {
                        sh 'pwd'
                        sh 'ls -al'
                        sh 'docker build --no-cache -t ${IMAGE_NAME}:${IMAGE_TAG} .' // Docker 이미지 빌드
                        sh 'docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:latest' // latest 태그 추가
                    }
                }
            }
        }

        stage('Push to GitLab Registry') { // GitLab Registry에 이미지 푸시 단계
            steps {
                container('docker') {
                    script {
                        sh 'echo $GITLAB_TOKEN_PSW | docker login ${GITLAB_REGISTRY} -u $GITLAB_TOKEN_USR --password-stdin' // GitLab Registry 로그인
                        sh 'docker push ${IMAGE_NAME}:${IMAGE_TAG}' // Docker 이미지 푸시
                        sh 'docker push ${IMAGE_NAME}:latest' // Docker 이미지 푸시(latest 태그)
                    }
                }
            }
        }
    }

    post {
        always {
            // 작업 완료 후 정리
            container('docker') {
                script {
                    sh 'docker logout ${GITLAB_REGISTRY} || true'
                    sh 'docker system prune -f || true'
                }
            }
        }
        success {
            echo "Pipeline succeeded! Image ${IMAGE_NAME}:${IMAGE_TAG} pushed to GitLab Registry"
        }
        failure {
            echo "Pipeline failed!"
        }
    }
}
```
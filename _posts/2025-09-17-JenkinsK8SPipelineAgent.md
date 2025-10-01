---
title : Jenkins 파이프라인을 K8S Pod로 사용하기 위한 설정
date : 2025-09-17 09:00:00 +09:00
categories : [Kubernetes, Jenkins]
tags : [kubernetes, jenkins, cicd, jenkins pipeline, jenkins agent] #소문자만 가능
---

> Jenkins 파이프라인이 K8S pod로 사용되도록 설정
> 

> Helm으로 설치 시에는 하기 설정들이 자동으로 세팅되어 있으니 설정 불필요
> 
1. Jenkins에 Kubernetes 플러그인 설치 확인
    - Jenkins가 K8S API와 통신해서 pod(파이프라인 에이전트)를 생성/삭제 가능하도록 설정
2. Built-In-Node excutor 설정
    - 파이프라인 에이전트를 k8s pod에서 동적으로 실행할거라 built-in-node(jenkins 서버 pod)의 Number of executors을 0으로 설정
    
    ![JenkinsK8SPipelineAgent1.png](/assets/img/kubernetes/JenkinsK8SPipelineAgent1.png)
    
3. K8S Configure
    - Manage Jenkins > Clouds > New Cloud / Configure
    - K8S URL로 Test Connection 확인
    
    ![JenkinsK8SPipelineAgent2.png](/assets/img/kubernetes/JenkinsK8SPipelineAgent2.png)
    
    - Jenkins URL/Tunnel
        - Jenkins URL: http://jenkins.jenkins.svc.cluster.local:8080
        - Jenkins tunnel: jenkins-agent.jenkins.svc.cluster.local:50000
            - jenkins tunnel은 jnlp(파이프라인 에이전트 관리 pod)로 연결되는 FQDN 입력
            - jnlp 포트가 50000
        - url, tunnel에 도메인을 넣으면 IP가 바뀌더라도 연결 가능
            - K8S FQDN: 클러스터 DNS 규칙에 의해 만들어짐
                - FQDN 룰: <Service명>.<Namespace명>.svc.<cluster domain>
                - cluster domain 기본값은 cluster.local
            - FQDN 예시
                - jenkins FQDN: jenkins.jenkins.svc.cluster.local
                - jenkins-agent FQDN: jenkins-agent.jenkins.svc.cluster.local
                
                ```bash
                # kubectl get svc -n jenkins
                NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)     AGE
                jenkins         ClusterIP   10.233.54.119   <none>        8080/TCP    11d
                jenkins-agent   ClusterIP   10.233.57.212   <none>        50000/TCP   11d
                ```
                

4. Pod Template

- 파이프라인 에이전트가 K8S에서 Pod로 관리될 때의 Pod 설정을 정의

    ![JenkinsK8SPipelineAgent3.png](/assets/img/kubernetes/JenkinsK8SPipelineAgent3.png)

---
title : K8S에 구축된 Jenkins-GitLab 파이프라인 설정
date : 2025-10-16 09:00:00 +09:00
categories : [Kubernetes, Jenkins]
tags : [Kubernetes, k8s, self managed k8s, jenkins, gitlab, jenkins pipeline] #소문자만 가능
---

### 파이프라인 설정

- Definition: `Pipeline script from SCM`
    - SCM: `Git`
        - Jenkinsfile을 Git으로 읽도록 설정
- git repository url: `http://gitlab-webservice-default.gitlab.svc.cluster.local:8181/root/test.git`
    - Jenkins 서버와 GitLab 서버가 모두 같은 k8s 클러스터에 있으므로, 외부 서버에서 url 설정하던 것처럼 하긴 어려움
    - GitLab UI Service의 FQDN인 `gitlab-webservice-default.gitlab.svc.cluster.local` 을 사용
- Credentials: `GitLab PAT`
    - GitLab의 사용자/리포지토리에서 발급한 PAT를 Jenkins Credentials에 등록해서 사용
- Branches to build: `*/main`
    - 빌드할 브랜치 선택

![JenkinsK8SPipelineGitLab01.png](/assets/img/kubernetes/JenkinsK8SPipelineGitLab01.png)

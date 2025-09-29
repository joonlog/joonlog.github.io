---
title : K8S에 구축한 GitLab UI 접근 시 접근 포트가 제거된 상태로 강제 리다이렉션 TroubleShooting
date : 2025-09-25 09:00:00 +09:00
categories : [Kubernetes, GitLab]
tags : [Kubernetes, k8s, self managed k8s, troubleshooting, gitlab, gitlab helm chart] #소문자만 가능
---

> GitLab을 Helm Chart로 구축 + 브라우저에서 GitLab 접근 시 80 포트가 아닌 다른 포트로 접근해야 할 경우 발생하는 이슈
GitLab Helm Chart의 포트 설정이 80/443으로 고정되어 있기 때문에 발생
> 

### 환경

1. 자체 관리형 kubernetes이기에 외부용 L7 로드밸런서로 haproxy를 사용 중인 경우
2. helm chart로 GitLab 서버를 구축한 경우
3. 외부에서 브라우저로 GitLab UI 접근 시 80포트가 아닌 8080과 같은 별도의 포트로만 접근해야하는 제약사항이 있는 경우

### 문제 상황

- helm chart로 GitLab 서버 구축 후 <도메인>:<포트>로 접근 시 포트가 제거된 후 도메인으로만 강제 리다이렉션 되는 현상 발생
    - 404 에러
- 브라우저 → HAProxy(8080) → GitLab(80) → 리다이렉트 응답

### 해결 방법

- HAproxy 설정을 통해 포트를 포함해서 리다이렉트 하도록 설정
- `http-response replace-value Location ^http://gitlab\.<도메인>/(.*)$ http://gitlab.<도메인>:8080/\\1`
    - 정규식 분석:
        - ^http://gitlab\.<도메인>/(.*)$ - 매칭 패턴
        - http://gitlab.<도메인>:8080/\\1 - 치환 결과
        - \1 = 첫 번째 그룹 (.*) 내용 (URL 경로)
        - 변환 예시:
            - 원본: Location: http://gitlab.<도메인>/users/sign_in
            - 변환: Location: http://gitlab.<도메인>/users/sign_in
    - 외부 브라우저에서 GitLab 접근 시 포트를 붙여서 접근해야 될 경우 반드시 설정
    - 설정하지 않으면 접근 시 자동으로 포트가 제거된 상태로 리다이렉션되어 404 에러
    
    ```bash
    frontend unified_frontend_8080
        bind *:8080
        mode http
        option forwardfor
    
        # Host 기반 라우팅
        ...
        use_backend metallb_backend_gitlab if { hdr(host) -i gitlab.<도메인> }
        ...
        
    backend metallb_backend_gitlab
        http-response replace-value Location ^http://gitlab\.<도메인>/(.*)$ http://gitlab.<도메인>:8080/\\1
        server gitlab 172.27.1.102:80
    ```
    

### 결과

- 브라우저가 8080포트를 포함해서 리다이렉트
- HAProxy를 통한 접속 유지
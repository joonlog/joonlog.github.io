---
title: "GitLab Container Registry에 Docker Push 실패"
date: 2025-12-04T09:00:00+09:00
categories: ["Linux", "GitLab"]
tags: ["linux", "gitlab", "gitlab.rb", "external_url"]
---


> GitLab Omnibus 설정 간 URL에 공인/사설 IP 혼용 시 Docker Push 실패 문제 해결
> 

### 문제 상황

- GitLab 웹 콘솔은 외부에서 접근해야 하니 공인 IP로 접근 필요
- Jenkins나 다른 WAS 서버에서 이미지 Pull/Push할때 사설 IP로 접근 필요
    - 그래서 gitlab.rb 파일에 `external_url`은 공인 IP, `registry_external_url`은 사설 IP 로 설정한 상황
    - 여기서 docker push 테스트 시 푸시가 중간에서 끊김

### 원인 분석

- GitLab Container Registry는 push 과정에서 이미지 blob을 업로드할 때, GitLab의 canonical URL(=external_url)을 기준으로 내부 리다이렉트 동작을 수행
- 동작 과정
    1. 사용자는 사설 IP(`registry_external_url`)로 push 요청을 보냄
    2. Registry는 manifest 처리 과정에서 GitLab API 엔드포인트를 호출함 → 이때 GitLab의 external_url(공인 IP)을 참조
    3. 사설 → 공인으로 트래픽 루프가 발생
    4. 클라우드 서버라 루프백이 막혀 있어 트래픽이 중간에 끊김
- 즉 GitLab의 두 URL이 서로 다른 네트워크 망을 가리킴으로 인해 redirect loop 구조가 발생

### 조치

- `external_url`, `registry_external_url`을 IP가 아닌 도메인으로 설정
- 도메인은 공인 IP와 매핑하고, 사설 IP가 필요한 WAS 서버나 Jenkins 서버에서는 `hosts` 설정으로 IP를 사설 IP로 지정되도록 설정

### 결론

GitLab Container Registry는 push 과정에서 canonical URL을 기반으로 내부 API 호출 및 redirect가 발생하기 때문에 `external_url`과 `registry_external_url` 이 서로 다른 네트워크를 가리키면 push 실패가 발생할 수 있다.

이를 해결하려면 url을 도메인으로 설정해서 리다이렉트가 발생하지 않도록 해야 한다.
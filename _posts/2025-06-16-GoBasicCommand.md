---
title : Go 기본 명령어
date : 2025-06-16 09:00:00 +09:00
categories : [Go, Go Basic]
tags : [go, go build, go install, go get, go tidy]  #소문자만 가능
---

### 컴파일 명령어

- go build
    - main 패키지면 현재 디렉터리에 실행 파일 생성
        - main 패키지 없다면 실행 파일 생성 안 됨 (컴파일만 됨)
- go install
    - main 패키지면 실행 파일이 `$GOBIN` 또는 `$GOPATH/bin` 에 생성됨
        - main 패키지 없다면 실행 파일 없음, 캐시에만 저장
    - 아래 코드처럼 실행 가능한 go로 작성된 cli 툴들 가져오기 가능(`GOBIN`에 등록됐을 시)
        
        ```bash
        go install github.com/user/tool@latest
        ```
        
    - 환경 변수 설정
    
    ```bash
    # GOBIN 기본 경로 확인
    go env GOBIN
    # 없으면 GOPATH/bin이 기본
    go env GOPATH
    
    # 실행 파일 전역 사용을 위해 PATH 설정
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
    source ~/.bashrc
    ```

---

### 모듈 라이브러리 명령어

- go get
    - main 패키지가 없는 모듈 라이브러리(`go.mod`)에 추가하는 명령어
    
    ```bash
    go get github.com/user/package
    ```
    
- go mod tidy
    - go.mod에 필요 없는 모듈 제거
    - 누락된 의존성 자동 추가

---

### 코드 스타일 정리

- 들여쓰기, 공백, import 정렬 등을 표준화
    
    ```go
    go fmt           # 현재 디렉토리 전체 포맷
    go fmt file.go   # 특정 파일만 포맷
    go fmt ./...     # 하위 디렉토리까지 전체 포맷. bash 문법이 아닌 go 문법.
    ```
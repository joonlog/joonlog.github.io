---
title : Go 설치와 실행
date : 2025-05-26 09:00:00 +09:00
categories : [Go, Go Basic]
tags : [go, go install, go run] #소문자만 가능
---

공식 문서: https://go.dev/doc/install

- 설치
    
    ```bash
    wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
    tar -xvzf go1.24.3.linux-amd64.tar.gz
    sudo mv go /usr/local/
    ```
    
- `.bashrc`에 영구 적용
    
    ```bash
    go env GOROOT
    
    # 위 명령어에서 나온 값 변수로 저장
    # 기본값: /usr/local/go
    echo 'export GOROOT=$(go env GOROOT)' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOROOT/bin' >> ~/.bashrc
    source ~/.bashrc
    ```
    
- 설치 확인
    
    ```bash
    go version
    ```
    

---

- go 모듈 초기화
    - Go 프로젝트 시작 시 go.mod 파일을 만들고 의존성 관리 가능하게 함
    - 보통 하나의 `프로젝트(리포지토리)`당 한 번만 실행
    - 모듈 이름은 GitHub 경로처럼 짓는 것이 관례이며, 실제로 GitHub에 올릴 경우 import 경로로 사용됨
    
    ```bash
    mkdir go-basic && cd go-basic
    go mod init github.com/joonlog/go-basic
    ```
    
    - 이 명령은 현재 디렉토리를 `github.com/joonlog/go-basic`이라는 이름의 모듈로 선언함
        
        → GitHub에 올리지 않아도 작동하지만, **나중에 배포하려면 실제 리포지토리 경로와 맞춰야 함**
        
    - go.mod 파일 생성됨
    
    ```bash
    go-workspace/
    └── myproject/                ← 모듈 디렉터리
        ├── go.mod                ← 모듈 정보
        └── main.go
    ```
    

---

- main.go 작성
    
    > package main의 목적: go run / go build로 실행 가능한 바이너리 생성
    반드시 func main 존재 필요
    > 
    
    ```go
    package main
    
    import "fmt"
    
    func main() {
        fmt.Println("Hello, Go!")
    }
    ```
    

---

- 실행
    
    ```bash
    go run main.go
    ```
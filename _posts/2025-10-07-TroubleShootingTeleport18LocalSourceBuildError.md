---
title : Teleport 18버전 로컬 소스 빌드 오류 TroubleShooting
date : 2025-10-07 09:00:00 +09:00
categories : [Go, Teleport]
tags : [go, teleport 18, local source build] #소문자만 가능
---

> Teleport 18버전을 소스 빌드하면서 겪은 의존성 관련 문제를 정리
Node, Rust, Corepack 등의 버전이 중요
> 

### **1. `Corepack is not installed` 오류**

- 소스 코드 다운로드 후 make

```bash
wget https://github.com/gravitational/teleport/archive/refs/tags/v18.0.0.tar.gz
mkdir teleport && tar -xzf v18.0.0.tar.gz -C teleport --strip-components=1
cd teleport
rm -rf .github/dependabot.yml
rm -rf e

make full
```

- make full 시 오류 발생
    
    ```bash
    Error: Corepack is not installed, cannot enable pnpm. See the installation guide https://pnpm.io/installation#using-corepack
    make[2]: *** [Makefile:1779: ensure-js-package-manager] Error 1
    make[2]: Leaving directory '/root/teleport'
    make[1]: *** [Makefile:1803: ensure-js-deps] Error 2
    make[1]: Leaving directory '/root/teleport'
    make: *** [Makefile:1762: ensure-webassets] Error 2
    ```
    
    - Teleport 웹 UI는 React 기반
        - `pnpm` 패키지 매니저를  Corepack으로 관리
- `build.assets/versions.mk`에 로컬 설치 시 필요 버전 명시되어 있음
    
    ```bash
    GOLANG_VERSION ?= go1.24.4
    NODE_VERSION ?= 22.14.0
    RUST_VERSION ?= 1.81.0
    WASM_PACK_VERSION ?= 0.12.1
    ```
    
    - nodejs 22에선 아래 명령어로 해결
    
    ```bash
    corepack enable
    corepack prepare pnpm@10.12.4 --activate
    ```
    

### 2. Rust + wasm-pack 버전 에러

- Rust 1.72 미만 또는 wasm-pack 설치 누락 시 오류 발생
    
    ```bash
     ERR_PNPM_RECURSIVE_RUN_FIRST_FAIL  @gravitational/shared@1.0.0 build-wasm: `node ../../scripts/clean-up-ironrdp-artifacts.mjs && RUST_MIN_STACK=16777216 wasm-pack build ./libs/ironrdp --target web`
    Exit status 1
     ELIFECYCLE  Command failed with exit code 1.
    /root/teleport/web/packages/teleport:
     ERR_PNPM_RECURSIVE_RUN_FIRST_FAIL  @gravitational/teleport@1.0.0 build: `pnpm build-wasm && vite build`
    Exit status 1
     ELIFECYCLE  Command failed with exit code 1.
    ```
    
    - wasm-pack 설치 후에도 wasm-opt 버전이 낮으면 빌드 실패
        - wasm-opt 123 버전 필요
        
        ```bash
        curl -LO https://github.com/WebAssembly/binaryen/releases/download/version_123/binaryen-version_123-x86_64-linux.tar.gz
        tar -xzf binaryen-version_123-x86_64-linux.tar.gz
        mv binaryen-version_123/bin/* /usr/local/bin/
        ```
        
        - 이후 캐시된 wasm-pack 바이너리 덮어쓰기
        
        ```bash
        mv ~/.cache/.wasm-pack/wasm-opt-*/bin/wasm-opt \
           ~/.cache/.wasm-pack/wasm-opt-*/bin/wasm-opt.bak
        mv /usr/local/bin/wasm-opt ~/.cache/.wasm-pack/wasm-opt-*/bin/wasm-opt
        chmod +x ~/.cache/.wasm-pack/wasm-opt-*/bin/wasm-opt
        ```
        

### 3. libfido2 관련 경고

- libfido2-devel 패키지 미설치 오류
    - Teleport의 MFA 기능을 포함하려면 설치 필요
    
    ```bash
    Warning: Building tctl without libfido2. Install libfido2 to have access to MFA.
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags "     kustomize_disable_go_plugin_support" -o build/tctl  -ldflags '-w -s -X k8s.io/component-base/version.gitVersion=v1.33.2' -trimpath -buildmode=pie  ./tool/tctl
    Warning: Building tsh without libfido2. Install libfido2 to have access to MFA.
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags "     kustomize_disable_go_plugin_support" -o build/tsh  -ldflags '-w -s -X k8s.io/component-base/version.gitVersion=v1.33.2' -trimpath -buildmode=pie  ./tool/tsh
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags " kustomize_disable_go_plugin_support" -o build/tbot  -ldflags '-w -s -X k8s.io/component-base/version.gitVersion=v1.33.2' -trimpath   ./tool/tbot
    ```
    
    libfido2-devel 패키지 설치
    
    ```bash
    yum install libfido2-devel -y
    ```
    

### 4. 디렉토리 여유 공간 확보

- Teleport 빌드 시 go cache, cargo, rustup, node_modules 등의 캐시 파일들이 쌓이며 최소 10GB 이상의 여유 공간 필요
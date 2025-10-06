---
title : Teleport 18 로컬 환경 설치 및 바이너리 빌드
date : 2025-10-06 09:00:00 +09:00
categories : [Go, Teleport]
tags : [go, teleport, teleport 18]  #소문자만 가능
---

> Teleport 18 소스 코드를 로컬에서 빌드하기 위한 환경 세팅
> 
- https://github.com/gravitational/teleport

## 1. 기본 패키지 설치

```
dnf groupinstall -y "Development Tools"
```

---

## 2. Node.js 22.14.0

```
curl -fsSL https://rpm.nodesource.com/setup_22.x | bash -
dnf install -y nodejs

node -v
```

---

## 3. Corepack ≥ 0.31 및 pnpm ≥ 10.11.0

```
corepack enable
corepack prepare pnpm@10.12.4 --activate

corepack -v
```

---

## 4. Python 3.8

```
dnf install -y python38
ln -s /usr/bin/python3.8 /usr/bin/python
export PYTHON=$(which python3.8)

python --version
```

---

## 5. Rust ≥ 1.72 및 wasm-pack 0.12.1

```
curl https://sh.rustup.rs -sSf | sh -s -- -y
source $HOME/.cargo/env
cargo install wasm-pack --locked --version 0.12.1

rustc --version
```

---

## 6. (선택) wasm-opt 123 (Binaryen)

- wasm-opt로 인한 빌드 오류 발생 시 최신 버전으로 교체
    - 111 버전은 오류 발생

```
cd /usr/local/bin
curl -LO https://github.com/WebAssembly/binaryen/releases/download/version_123/binaryen-version_123-x86_64-linux.tar.gz
tar -xzf binaryen-version_123-x86_64-linux.tar.gz
mv binaryen-version_123/bin/* /usr/local/bin/
rm -rf binaryen-version_123*

wasm-opt --version
```

---

## 7. libfido2-devel

```
dnf install -y libfido2-devel
```

---

## 8. Teleport 소스코드 다운로드 및 빌드

```
wget https://github.com/gravitational/teleport/archive/refs/tags/v18.0.0.tar.gz
mkdir teleport && tar -xzf v18.0.0.tar.gz -C teleport --strip-components=1
cd teleport

make clean
make full
```
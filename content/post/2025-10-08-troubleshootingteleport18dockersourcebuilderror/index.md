---
title: "Teleport 18버전 도커 환경에서의 소스 빌드 오류 TroubleShooting"
date: 2025-10-08T09:00:00+09:00
categories: ["Go", "Teleport"]
tags: ["go", "teleport 18", "docker source build"]
---


> Teleport 18을 도커 환경에서 빌드하는 과정에서 발생한 오류 정리
> 

### 1. `UID 0 is not unique` 오류

```bash
wget https://github.com/gravitational/teleport/archive/refs/tags/v18.0.0.tar.gz
mkdir teleport && tar -xzf v18.0.0.tar.gz -C teleport --strip-components=1
cd teleport
rm -rf .github/dependabot.yml
rm -rf e

make -C build.assets build-binaries
```

- 에러
    
    ```bash
     => CANCELED [git2 1/2] RUN yum-builddep -y git                                    3.8s
     => ERROR [buildbox  1/27] RUN (groupadd ci --gid=0 -o && useradd ci --uid=0 --gi  3.7s
    ------
     > [buildbox  1/27] RUN (groupadd ci --gid=0 -o && useradd ci --uid=0 --gid=0 --create-home --shell=/bin/sh &&      mkdir -p -m0700 /var/lib/teleport && chown -R ci /var/lib/teleport):
    3.000 useradd: UID 0 is not unique
    ------
    Dockerfile-centos7:218
    --------------------
     217 |     ARG GID
     218 | >>> RUN (groupadd ci --gid=$GID -o && useradd ci --uid=$UID --gid=$GID --create-home --shell=/bin/sh && \
     219 | >>>      mkdir -p -m0700 /var/lib/teleport && chown -R ci /var/lib/teleport)
     220 |     
    --------------------
    ERROR: failed to solve: process "/bin/sh -c (groupadd ci --gid=$GID -o && useradd ci --uid=$UID --gid=$GID --create-home --shell=/bin/sh &&      mkdir -p -m0700 /var/lib/teleport && chown -R ci /var/lib/teleport)" did not complete successfully: exit code: 4
    exit status 1
    make: *** [Makefile:213: buildbox-centos7] Error 1
    make: Leaving directory '/data/teleport/build.assets'
    ```
    
    - root사용자로 작업하던 상황에서 Dockerfile 내 `ARG UID` 코드로 인해 uid, gid가 0으로 설정됨
        - 소스 코드 디렉토리 권한 부여 + UID, GID 지정해서 다시실행
        
        ```bash
        chown -R 1000:1000 ./*
        nohup make -C build.assets build-binaries UID=1000 GID=1000 > build.log 2>&1 &
        ```
        

### 2. CentOS7 환경 빌드 실패

- `make -C build.assets build-binaries` 실행 시 centos7에서 실행 ⇒ 버전이 낮아서 빌드 실패
    - `nohup make -C build.assets release-ng > build.log 2>&1`
    - release-ng를 사용한 make가 권장되지면 teleport-buildbox-thirtparty가 18버전 이미지가 없어서 에러 발생
        - https://github.com/gravitational/teleport/pkgs/container/teleport-buildbox-thirdparty
        
        ```bash
        failed to pull teleport-buildbox-thirdparty:teleport18
        ```
        
    - 로컬에서 thirdparty 이미지 빌드 필요
        - 이미지 빌드에 1시간 이상 걸리니 주의
        
        ```bash
        docker buildx build --platform=linux/amd64 --load --tag test-thirdparty:manual -f build.assets/buildbox/Dockerfile-thirdparty build.assets/buildbox
        
        docker tag test-thirdparty:manual <username>/teleport-buildbox-thirdparty:teleport18
        docker push <username>/teleport-buildbox-thirdparty:teleport18
        ```
        
        - thirdparty 이미지 빌드 후에는 아래 명령어로 바이너리 빌드
            
            ```bash
            nohup make -C build.assets release-ng UID=1000 GID=1000 BUILDBOX_THIRDPARTY=<username>/teleport-buildbox-thirdparty:teleport18 > build.log 2>&1 &
            ```
            

### **3. `error obtaining VCS status: exit status 128` 에러**

- 컨테이너 내부에 .git 디렉토리가 있을 경우 `go build`가  VCS 정보를 가져오다 실패
- https://www.linkedin.com/pulse/addressing-error-obtaining-vcs-status-issue-leo-liu-rhuye/
    
    ```bash
    error obtaining VCS status: exit status 128
            Use -buildvcs=false to disable VCS stamping.
    make[2]: *** [Makefile:388: build/teleport] Error 1
    make[2]: Leaving directory '/home/teleport'
    make[1]: *** [Makefile:354: all] Error 2
    make[1]: Leaving directory '/home/teleport'
    make: *** [Makefile:507: full] Error 2
    exit status 2
    make[1]: *** [Makefile:624: release-ng] Error 2
    make[1]: Leaving directory '/data/teleport/build.assets'
    ```
    
- 루트 Makefile 상단에 `GOFLAGS ?= -buildvcs=false` 를 통해 비활성화

### 4. 정상 동작

```bash
nohup make -C build.assets release-ng UID=1000 GID=1000 BUILDBOX_THIRDPARTY=<username>/teleport-buildbox-thirdparty:teleport18 > build.log 2>&1 &
```

- 빌드 완료
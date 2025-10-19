---
title : Pokerogue 게임을 Container 이미지로 빌드하기 위한 방법
date : 2025-10-19 09:00:00 +09:00
categories : [Container, Pokerogue]
tags : [pokerogue, rogueserver, docker, container, dockerfile, oauth]  #소문자만 가능
---

- Pokerogue 이미지 빌드
    
    > Pokerogue 게임을 컨테이너 이미지로 빌드하기 위한 방법
    > 
    
    ### 개요
    
    - Pokerogue 게임을 K8S 환경에서 동작하게 하기 위해서 약간의 소스코드 수정과 Dockerfile 작성이 필요
        - 공식 GitHub에는 로컬 기준으로 코드가 작성된 부분이 많음
    - Pokerogue는 Vite 기반 React 프론트엔드와 Go 기반 백엔드 그리고 Mariadb로 구성
        
        ![PokerogueContainerImageBuild01.png](/assets/img/container/PokerogueContainerImageBuild01.png)
    
    
    ## Frontend
    
    - 공식 GitHub: https://github.com/pagefaultgames/pokerogue
    - TypeScript로 작성
    - Vite 번들러를 사용하는 React SPA
    - Pokerogue의 UI, 정적파일, 로직 및 랜더링 제공
    - `/api/*`로 백엔드 API 호출
    - OAuth 외부 로그인 UI 제공
    
    ### Frontend 이미지 빌드
    
    1. **Dockerfile**
        
        ```docker
        ARG NODE_VERSION=22.14
        ARG OS=alpine
        
        # Build stage
        FROM node:${NODE_VERSION}-${OS} AS build
        
        RUN apk add --no-cache git rsync
        WORKDIR /app
        
        RUN corepack enable && corepack prepare pnpm@10.14.0 --activate
        
        COPY . .
        
        # Initialize git submodules
        RUN git submodule update --init --recursive
        
        # Install dependencies with cache
        RUN --mount=type=cache,target=/root/.pnpm-store \
            pnpm install --frozen-lockfile
        
        # Production build
        RUN pnpm run build
        
        # Production stage
        FROM nginx:alpine
        
        COPY nginx.conf /etc/nginx/conf.d/default.conf
        
        # Copy built files from build stage
        COPY --from=build /app/dist /usr/share/nginx/html
        
        EXPOSE 80
        CMD ["nginx", "-g", "daemon off;"]
        ```
        
        - 공식 GitHub에는 beta 브랜치에만 Dockerfile이 존재
            - development mode로 실행되게 작성되어 있어서 Dockerfile은 새로 전부 작성
        1. Build Stage와 Production Stage 분리
            - Build Stage
                - Node.js 정적 파일 빌드
                - Git submodule 초기화
                - pnpm dependencies 설치
                - pnpm 빌드로 `/app/dist` 생성
            - Production Stage
                - Nginx Alpine 이미지
                - 빌드된 정적 파일만 복사
                - Node.js 런타임 포함하지 않음
        2. Git Submodule 처리
            - Pokerogue는 assets와 locales를 Git submodule로 관리
            - 최종 이미지에는 빌드된 `dist` 포함되므로 빌드에 사용된 `.git`은 제외됨
        3. Nginx 사용
            - SPA 라우팅
                - https://pokerogue.com/ 이외에 /game, /home 등 제공하지 않는 경로로 접근 시 index.html을 제공하도록 설정
            - gzip 압축
                - JavaScript, CSS 파일들을 gzip으로 압축하여 전송
            - 보안 헤더 적용
            - 정적 파일 캐싱
                - 이미지, 폰트, JS, CSS 등을 1년간 브라우저 캐시
            
            ```bash
            server {
                listen 80;
                server_name localhost;
                root /usr/share/nginx/html;
                index index.html;
            
                # SPA routing - all routes go to index.html
                location / {
                    try_files $uri $uri/ /index.html;
                }
            
                gzip on;
                gzip_vary on;
                gzip_min_length 1024;
                gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json application/javascript;
            
                add_header X-Frame-Options "SAMEORIGIN" always;
                add_header X-Content-Type-Options "nosniff" always;
                add_header X-XSS-Protection "1; mode=block" always;
            
                location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
                    expires 1y;
                    add_header Cache-Control "public, immutable";
                }
            }
            ```
            
    2. **`.env` 파일 관리**
        - 환경변수들을 빌드 시점에 포함시켜야함
        - `.dockeringnore` 파일에 명시되어 있던 `.env` 주석처리
            
            ```bash
            # .dockerignore
            node_modules
            *.log
            *.md
            .gitignore
            Dockerfile
            #.env 
            ```
            
        - Vite의 `.env` 우선순위
            - Vite는 다음 순서로 `.env` 파일 로드
                - 나중에 로드된 파일이 이전 파일의 변수값을 덮어씀
                1. `.env` - 모든 환경에서 로드
                2. `.env.[mode]` - 특정 모드 (예: `.env.production`)
            - 따라서 pnpm run build 실행 시 production 모드로 실행되므로 최종적으로 `.env.production`에 있는 환경변수를 참조하게 됨
        - `.env` | `.env.production`
            - VITE_SERVER_URL: 실제 게임을 서비스할 도메인으로 변경
            - VITE_DISCORD_CLIENT_ID: Discord OAuth 로그인을 위한 본인 값 입력
            - VITE_GOOGLE_CLIENT_ID: Google OAuth 로그인을 위한 본인 값 입력
            
            ```bash
            VITE_BYPASS_LOGIN=0
            VITE_BYPASS_TUTORIAL=0
            VITE_SERVER_URL=https://pokerogue.<도메인>:<포트>/api
            VITE_DISCORD_CLIENT_ID=<discord client id>
            VITE_GOOGLE_CLIENT_ID=<google client id>
            VITE_I18N_DEBUG=0
            ```
            
    3. 이미지 빌드
        
        ```bash
        docker build -t ksi05298/pokerogue-web:latest .
        docker push ksi05298/pokerogue-web:latest
        ```
        
    
    ## Backend
    
    - 공식 GitHub: https://github.com/pagefaultgames/rogueserver
    - Go로 작성
    - 사용자 인증, 게임 데이터 저장, 조회, 세션 관리 등을 수행하는 RESTFUL API 서버
        - `/account/*`, `/gamedata/*`, `/auth/*` 엔드포인트 제공
    - MariaDB와 연동
    
    ### Backend 이미지 빌드
    
    1. Dockerfile
        - Dockerfile은 공식 GitHub에 있는 파일 그대로 사용
            - scratch 이미지를 사용해서 정적 바이너리만 포함
        
        ```bash
        # SPDX-FileCopyrightText: 2024-2025 Pagefault Games
        #
        # SPDX-License-Identifier: AGPL-3.0-or-later
        ARG GO_VERSION=1.22
        
        FROM docker.io/library/golang:${GO_VERSION} AS builder
        
        WORKDIR /src
        
        COPY ./go.mod /src/
        COPY ./go.sum /src/
        
        RUN go mod download && go mod verify
        
        COPY . /src/
        
        RUN CGO_ENABLED=0 \
            go build -tags=devsetup -o rogueserver
        
        RUN chmod +x /src/rogueserver
        
        # ---------------------------------------------
        
        FROM scratch
        
        WORKDIR /app
        
        COPY --from=builder /src/rogueserver .
        COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
        
        EXPOSE 8001
        
        ENTRYPOINT ["./rogueserver"]
        ```
        
    2. OAuth 쿠키 설정 수정
        - 공식 rogueserver 코드에서 구글/디스코드 로그인에 필요한 OAuth 쿠키 설정에 문제 있어 수정 필요
        - `api/endpoints.go` 파일의 658줄 `handleProviderCallback` 함수 내부
            - 수정 전
                
                ```bash
                		http.SetCookie(w, &http.Cookie{
                			Name:     "pokerogue_sessionId",
                			Value:    sessionToken,
                			Path:     "/",
                			Secure:   true,
                			SameSite: http.SameSiteStrictMode,
                			Domain:   "pokerogue.net",
                			Expires:  time.Now().Add(time.Hour * 24 * 30 * 3), // 3 months
                		})
                ```
                
            - 수정 후
                
                ```bash
                		http.SetCookie(w, &http.Cookie{
                			Name:     "pokerogue_sessionId",
                			Value:    sessionToken,
                			Path:     "/",
                			Secure:   true,
                			SameSite: http.SameSiteLaxMode,
                			Domain:   "",
                			Expires:  time.Now().Add(time.Hour * 24 * 30 * 3), // 3 months
                		})
                ```
                
            - 변경 사항
                - SameSite: SameSiteStrictMode → SameSiteLaxMode
                    - OAuth callback은 cross-site 요청이 될 수 밖에 없음 (google.com → 내 도메인)
                    - Strict로 설정하면 cross-site 쿠키 전송이 차단돼서 CSRF오류가 발생됨
                    - Lax로 변경하여 안전한 도메인에서의  GET 요청(oauth callback)에서 쿠키 전송 허용
                - Domain: “pokerogue.net” → “”
                    - pokerogue 공식 도메인에서 빈 문자열로 변경해서 현재 도메인이 자동 적용되도록 설정
    
    3. 이미지 빌드 및 푸시
    
    ```bash
    docker build -t ksi05298/rogueserver:latest .
    docker push ksi05298/rogueserver:latest
    ```
    
    ### 참고
    
    Vite 환경변수 가이드: 
    
    - https://vitejs.dev/guide/env-and-mode.html
    
    OAuth SameSite Cookie 설정: 
    
    - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite
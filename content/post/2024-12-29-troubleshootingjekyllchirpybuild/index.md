---
title: "GitHub 블로그(Jekyll Chirpy) 오류 TroubleShooting"
date: 2024-12-29T09:00:00+09:00
categories: ["GitHub Blog", "Jekyll"]
tags: ["git", "github", "github blog", "github action error", "jekyll", "jekyll chirpy", "troubleshooting"]
---


## Jekyll - Chirpy 테마

### 로컬 테스트 오류

> `Unable to monitor directories for changes because iNotify max watches exceeded.`
> 
> - Linux의 `inotify` 파일 시스템 감시자의 최대 수 제한을 초과했기 때문에 발생
> - 1코어 1기가로 실행 중이기 때문에 제한인 것으로 추정
> - inotify 값 변경해서 해결
>     
>     ```bash
>     cat /proc/sys/fs/inotify/max_user_watches
>     	
>     sudo sysctl fs.inotify.max_user_watches=524288
>     ```
>     

### GitHub Action 의존성 오류

- git push 후에 github action에서 발생한 오류

> 의존성 누락으로 인해 _sass/main.bundle.scss에서 @use 'vendors/bootstrap' 참조 불가
> 
> 
> ```bash
> Error: Can't find stylesheet to import.
>   ╷
> 1 │ @use 'vendors/bootstrap';
>   │ ^^^^^^^^^^^^^^^^^^^^^^^^
>   ╵
>   main.bundle.scss 1:1                                                                           @use
>   /home/runner/work/joonlog.github.io/joonlog.github.io/assets/css/jekyll-theme-chirpy.scss 1:1  root stylesheet 
>   Conversion error: Jekyll::Converters::Scss encountered an error while converting 'assets/css/jekyll-theme-chirpy.scss':
>                     Can't find stylesheet to import.
>                     ------------------------------------------------
>       Jekyll 4.3.4   Please append `--trace` to the `build` command 
>                      for any additional information or backtrace. 
>                     ------------------------------------------------
> ```
> 
> **참고: https://github.com/cotes2020/jekyll-theme-chirpy/discussions/1809**
> 
> - ci 파일에 npm 의존성 주입 추가로 해결 가능
> - 아래 코드를 `jekyll.yml`의 `Build with Jekyll` Step 위에 작성
> 
> ```bash
> name: npm build
> run: npm install && npm run build
> ```
> 

### GitHub Action 의존성 오류2

- git push 후에 github action에서 발생한 오류

> 7시간 전까지 정상 동작 했는데 GitHub Action 중 오류
> 
> 
> ```bash
> ...
> npm warn EBADENGINE Unsupported engine {
> npm warn EBADENGINE   package: '@semantic-release/npm@12.0.1',
> npm warn EBADENGINE   required: { node: '>=20.8.1' },
> npm warn EBADENGINE   current: { node: 'v18.20.5', npm: '10.8.2' }
> npm warn EBADENGINE }
> npm warn EBADENGINE Unsupported engine {
> npm warn EBADENGINE   package: '@semantic-release/release-notes-generator@14.0.1',
> npm warn EBADENGINE   required: { node: '>=20.8.1' },
> npm warn EBADENGINE   current: { node: 'v18.20.5', npm: '10.8.2' }
> npm warn EBADENGINE }
> ...
> [js] [!] Error: Cannot find module '/home/runner/work/joonlog.github.io/joonlog.github.io/node_modules/@jridgewell/gen-mapping/dist/gen-mapping.umd.js'
> [js]     at createEsmNotFoundErr (node:internal/modules/cjs/loader:1177:15)
> [js]     at finalizeEsmResolution (node:internal/modules/cjs/loader:1165:15)
> [js]     at resolveExports (node:internal/modules/cjs/loader:590:14)
> [js]     at Function.Module._findPath (node:internal/modules/cjs/loader:664:31)
> [js]     at Function.Module._resolveFilename (node:internal/modules/cjs/loader:1126:27)
> [js]     at Function.Module._load (node:internal/modules/cjs/loader:981:27)
> [js]     at Module.require (node:internal/modules/cjs/loader:1231:19)
> [js]     at require (node:internal/modules/helpers:177:18)
> [js]     at Object.<anonymous> (/home/runner/work/joonlog.github.io/joonlog.github.io/node_modules/@jridgewell/source-map/dist/source-map.cjs:6:18)
> [js]     at Module._compile (node:internal/modules/cjs/loader:1364:14)
> ```
> 
> 정상 동작 하던 커밋으로 롤백했는데도 같은 에러 발생
> ⇒ GitHub 쪽 버전 문제?
> 
> warn 발생하는 node 문제인가 싶어 jekyll.yml에 아래 코드 추가해서 node 22이상으로 사용 후 해결됨 - npm 설치 후 빌드 시 node 20 이상이 필요하게 변경된 것으로 추정
> ⇒ 버전 지정 안하면 깃허브는 node를 18버전으로 사용
> 
> ```bash
>       - name: Setup Node.js
>         uses: actions/setup-node@v3
>         with:
>           node-version: '22.x'  # Node 22 이상 사용
> ```
>
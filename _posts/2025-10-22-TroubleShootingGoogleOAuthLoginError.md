---
title : 깃허브 코드로 배포한 Pokerogue 앱에서 구글/디스코드 외부 로그인이 실패하는 이슈
date : 2025-10-22 09:00:00 +09:00
categories : [Kubernetes, Pokerogue]
tags : [Kubernetes, k8s, self managed k8s, pokerogue, rogueserver, oauth, samesitelax] #소문자만 가능
---

> 공식 rogueserver 코드의 쿠키 설정 버그로 인해 구글 OAuth 로그인이 작동하지 않는 이슈
> 

### 환경

- Vite 기반 React 프론트엔드와 Go 기반 백엔드 그리고 Mariadb로 구성
- K8S + HAProxy를 통한 외부 접근

### 문제 상황

- 일반 계정 로그인은 정상 작동
- 구글/디스코드 계정 연결 시도 시 에러:
    1. 일반 로그인 후 설정 메뉴에서 구글 연결 클릭
    2. 구글 OAuth 인증 완료
    3. 메인 화면으로 리다이렉트되지만 설정이 저장되지 않고 새로고침만 됨
    4. 로그아웃 후 구글로 직접 로그인 시도 시 로그인 실패

## 에러 분석

### 1. 백엔드 로그 분석

```bash
kubectl logs -n pokerogue pokerogue-server
```

```
2025/10/14 23:02:24 [DEBUG] Google token response status: 200
2025/10/14 23:02:24 [DEBUG] Successfully extracted userId: 102118711772610730702
2025/10/14 23:02:47 /account/info: missing token
```

- 구글 OAuth 토큰 교환은 성공 (userId 추출 완료)
- DB에 googleId도 정상 저장됨 확인
- 하지만 세션 토큰이 누락되어 로그인 상태 유지 실패

### 2. 브라우저 네트워크 탭 분석

```
Request: GET /api/auth/google/callback
Response: 303 See Other
Set-Cookie: pokerogue_sessionId=...; SameSite=Strict; Domain=pokerogue.net
Sec-Fetch-Site: cross-site
```

**문제점 발견:**

1. `Domain=pokerogue.net`
    - 공식 pokerogue URL로 하드코딩 되어 있어 현재 도메인과 불일치
2. `SameSite=Strict`
    - 구글에서 돌아오는 요청이 `cross-site`이므로 쿠키가 전송되지 않음
- OAuth 리다이렉트 흐름:
    - `pokerogue.<도메인>` → `accounts.google.com` → `pokerogue.<도메인>/callback`
        - `SameSite=Strict`는 같은 사이트 내부 요청에만 쿠키 전송
        - 구글에서 돌아오는 것은 cross-site 요청으로 간주되어 쿠키 차단

## 해결 방법

- `api/endpoints.go` 파일의 쿠키 설정 수정:

```go
http.SetCookie(w, &http.Cookie{
    Name:     "pokerogue_sessionId",
    Value:    sessionToken,
    Path:     "/",
    Secure:   true,
    SameSite: http.SameSiteLaxMode,    // Strict → Lax 변경
    Domain:   "",                       // "pokerogue.net" → "" 변경
    Expires:  time.Now().Add(time.Hour * 24 * 30 * 3),
})
```

- **변경 사항**
    1. `Domain:   ""`
        - 빈 문자열로 설정하여 현재 도메인 자동 적용
    2. `SameSite: http.SameSiteLaxMode`
        - `Strict` → `Lax`
        - `Lax`는 안전한 cross-site GET 요청(OAuth 콜백 포함)에서 쿠키 전송 허용
        - OAuth 표준에서 권장하는 설정
        - 여전히 CSRF 공격 방어 가능 (위험한 POST 요청은 차단)

### 결과

- 일반 로그인 → 구글 계정 연결 → 로그아웃 → 구글 직접 로그인 모두 정상 작동
- 세션 쿠키가 올바르게 설정되어 로그인 상태 유지
- 브라우저 네트워크 탭에서 `Set-Cookie: ...; SameSite=Lax` 확인

### 참고

OAuth를 사용하는 많은 프로젝트에서 동일한 이슈가 발생했으며, 모두 `SameSite=Lax`로 해결

- https://github.com/oauth2-proxy/oauth2-proxy/issues/1663
- https://github.com/grafana/grafana/pull/18332
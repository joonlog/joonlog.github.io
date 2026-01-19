---
title: "K6 스크립트 및 결과값 해석"
date: 2026-01-19T09:00:00+09:00
categories: ["Testing", "k6"]
tags: ["k6", "performance test", "grafana k6", "k6 script", "login test"]
---


> 테스트 스크립트 작성 방법 및 결과값 해석에 대한 정리
> 
- JS 파일로 시나리오를 작성 후 k6 명령어로 실행
    - `k6 run test.js`

### 로그인 없이 테스트

- 100명의 사용자가 1분동안 1초마다 재요청하는 테스트
- `vus`: 동시에 실행되는 가상 사용자 수
- `duration`: 테스트 총 수행 시간
- `sleep(1)`: 사용자 요청 간 간격을 모델링
- `tags.name`: 요청 유형별 메트릭 분리를 위한 식별자

```jsx
import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  vus: 100,          // 동시 사용자
  duration: '1m',  // 테스트 시간
};

export default function () {

  let homeRes = http.get(
    '<URL>',
    {
      tags: {
        name: 'home',
        page: 'home'
      }
    }
  );

  check(homeRes, {
    'home 200': r => r.status === 200,
  });

  sleep(1);
}
```

```jsx
k6 run test.js
```

### 세션 로그인 성능 테스트

- 10명의 사용자가 1분동안 1초마다 재요청하는 테스트
- 세션 기반 로그인 성능 테스트
    - JWT, SSO 기반 로그인과는 스크립트 각각 다름
- 로그인 / 로그인 후 동작으로 분리
    - 로그인 url에 POST로 계정정보를 전달
        - 계정정보는 환경변수 `LOGIN_ID`, `LOGIN_PW`에 사전 세팅

```jsx
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10,
  duration: '1m',
};

export default function () {

  // 로그인 요청
  let loginRes = http.post(
    '<로그인 URL>',
    {
      id: __ENV.LOGIN_ID,
      password: __ENV.LOGIN_PW,
    },
    {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      tags: { name: 'login' },
    }
  );

  check(loginRes, {
    'login status 200': r => r.status === 200,
  });

  // 로그인 이후 접근 페이지
  let afterLoginRes = http.get(
    '<URL>',
    {
      tags: { name: 'after_login' },
    }
  );

  check(afterLoginRes, {
    'after login 200': r => r.status === 200,
  });

  sleep(1);
}
```

```jsx
k6 run logintest.js
```

## 결과

- `http_req_duration`
    - avg: 평균 응답시간
    - p95: 상위 95% 요청의 응답시간
    - max: 최장 응답시간. 이상치(outlier) 존재 여부 판단
- `http_req_waiting`
    - 요청 전송 및 응답 수신을 제외한 서버 처리 지연
    - `http_req_duration`과 거의 동일한 경우 서버 처리 지연
- `http_req_failed`
    - HTTP 오류(4xx/5xx), 타임아웃, 네트워크 실패 포함
- `iteration_duration`
    - 여러 요청 + sleep을 포함한 시나리오 1회 수행 시간

### 정상 패턴

- avg, p95가 **완만하게 증가**
- 오류율 0 유지
- iteration_duration 선형 증가

### 비정상 패턴

- 특정 VU 구간에서 p95 급증
- http_req_failed 증가
- iteration_duration 비선형 증가

### K6 시각화 보고서 생성

```jsx
$env:K6_WEB_DASHBOARD="true"
$env:K6_WEB_DASHBOARD_EXPORT="report.html"
k6 run test.js
```

- 테스트 결과를 시각화된 html로 생성
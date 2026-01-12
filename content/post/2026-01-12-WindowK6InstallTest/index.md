---
title: "윈도우 K6 설치 및 테스트"
date: 2026-01-12T09:00:00+09:00
categories: ["Testing", "k6"]
tags: ["window", "k6", "performance test"]
---


### K6

- VU(가상 사용자) 기반의 HTTP 부하 생성 도구
- 스크립트에 정의된 시나리오를 수행

### 환경

- 실행 환경: Windows
- 테스트 대상: 홈페이지 사이트

### K6 설치

- winget을 사용해 설치

```powershell
winget install k6
k6 version
```

### K6 사용법

- JS 파일로 시나오를 작성 후 k6 명령어로 실행
- 128명의 사용자가 1분동안 1초의 대기시간 후 재요청하는 테스트

```jsx
import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  vus: 128,          // 동시 사용자
  duration: '1m',  // 테스트 시간
};

export default function () {

  let homeRes = http.get(
    'URL',
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

- `vus`: 동시에 실행되는 가상 사용자 수
- `duration`: 테스트 총 수행 시간
- `sleep(1)`: 사용자 요청 간 간격을 모델링
- `tags.name`: 요청 유형별 메트릭 분리를 위한 식별자

### 결과

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

### K6 대시보드

```jsx
$env:K6_WEB_DASHBOARD="true"
$env:K6_WEB_DASHBOARD_EXPORT="report.html"
k6 run test.js
```

- 테스트 결과를 시각화된 html로 생성
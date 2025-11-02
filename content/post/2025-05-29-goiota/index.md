---
title: "Go iota"
date: 2025-05-29T09:00:00+09:00
categories: ["Go", "Go Basic"]
tags: ["go", "iota", "const"]
---


> Go 특유의 상수용 자동 증가 매크로 키워드
> 
- `const` 블록에서만 사용 가능
- 첫 번째 상수는 0부터 시작
- 아래로 한 줄씩 내려갈 때마다 1씩 자동 증가

---

- 기본 예제

```go
const (
    A = iota  // 0
    B         // 1
    C         // 2
)
```

`iota`는 행 기준으로 증가 → 위 예제에서 `A=0, B=1, C=2`

---

- 중간에 값을 지정해도 다음 `iota`는 계속 증가

```go
const (
    A = iota       // 0
    B = 100        // 고정값
    C = iota       // 2
)
```

---

- 비트 연산자와 조합 (플래그)

```go
const (
    Read = 1 << iota  // 1 << 0 → 1
    Write             // 1 << 1 → 2
    Exec              // 1 << 2 → 4
)
```

→ `Read=1`, `Write=2`, `Exec=4` → 비트 플래그 용도로 자주 사용

---

- 여러 개 동시에 선언할 때도 동작

```go
const (
    x, y = iota, iota   // x=0, y=0
    a, b = iota, iota   // a=1, b=1
)
```
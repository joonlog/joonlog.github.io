---
title: "Go Blank Identifier"
date: 2025-05-30T09:00:00+09:00
categories: ["Go", "Go Basic"]
tags: ["go", "blank identifier"]
---


- Go에서 `_`는 사용하지 않는 값 또는 결과를 버리는 용도로 사용
- **읽기 전용**이며, 어떤 값도 저장하거나 참조 불가
- 실제 메모리에 변수로 존재하지 않음

---

### 주요 사용 예

- 다중 리턴값 중 일부 무시할 때

```go
val, _ := someFunc()  // 두 번째 리턴값 무시
```

- 반복문에서 인덱스나 값을 무시할 때

```go
for _, v := range list {
    fmt.Println(v)  // 인덱스는 필요 없을 때
}
```

- 인터페이스 구현 시, 사용하지 않는 메서드 인자 무시

```go
func handler(_ int) {
    // 인자를 사용하지 않음
}
```

- import된 패키지를 코드에서 직접 사용하지 않을 때

```go
import _ "net/http/pprof"  // side-effect용 import
```

- 변수 선언 후 사용하지 않아 에러가 날 때 회피용

```go
msg := "hello"
_ = msg
```
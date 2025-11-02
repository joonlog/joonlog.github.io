---
title: "Go의 기본 자료형과 단축 선언"
date: 2025-05-28T09:00:00+09:00
categories: ["Go", "Go Basic"]
tags: ["go", "data type", "fmt", "casting", "short declaration"]
---


- **정수형**
    - `int`, `int8`, `int16`, `int32`, `int64`
    - `uint` 계열은 양수 전용
    - 보통은 그냥 `int` 사용
- **실수형**
    - `float32`, `float64` (기본은 `float64`)
    - 소수점 숫자 저장용
- **문자형**
    - `rune`: 유니코드 문자 하나 (`int32`)
    - `'A'`, `'한'` , 이모티콘을 포함한 문자 저장 가능
- **문자열**
    - `string`: UTF-8 인코딩 문자열
    - 불변(immutable)
    - 슬라이싱, 인덱싱 가능
- **Boolean**
    - `bool`: `true` / `false`
    - 조건문 등에 사용
- **Byte**
    - byte는 uint8의 alias
    - 문자열 인코딩 관련 작업에서 자주 사용

```go
var i int = 10
var f float64 = 3.14
var ch rune = '가'
var s string = "hello"
var b bool = true
var b byte = 'A' // byte는 uint8
```

### 기본값

- Go에서는 변수 선언 시 초기값을 주지 않으면 자료형에 맞는 기본값(zero value)으로 초기화됨

```go
var i int      // 0
var s string   // ""
var b bool     // false
var f float64  // 0.0
```

---

### 출력 포맷 (`fmt.Printf`)

| 포맷 | 설명 | 예시 |
| --- | --- | --- |
| `%d` | 정수 | `fmt.Printf("%d", 10)` |
| `%f` | 실수 | `fmt.Printf("%.2f", 3.14)` |
| `%s` | 문자열 | `fmt.Printf("%s", "go")` |
| `%c` | 문자 | `fmt.Printf("%c", 'A')` |
| `%v` | 기본 포맷 | `fmt.Printf("%v", 변수)` |
| `%T` | 타입 | `fmt.Printf("%T", 변수)` |

---

### 형 변환

- Go는 자동 형 변환 없음 → 반드시 명시적 변환 필요
- 정수 ↔ 실수, 정수 ↔ uint 변환 시 주의

```go
var a int = 42
var b float64 = float64(a)
var c uint = uint(b)
```

- rune <===> string 인덱싱 시 주의
    - `string[index]`는 `byte` 리턴 → 유니코드 문자는 `rune`으로 변환해야 제대로 출력 가능
    
    ```go
    s := "가나다"
    fmt.Printf("%c\n", s[0])           // 깨짐
    fmt.Printf("%c\n", []rune(s)[0])   // 가
    ```
    

---

### 단축 선언 (`:=`)

- 타입 추론과 동시에 변수 선언
- 주로 함수 내부에서 사용

```go
name := "Go"     // string
score := 95      // int
rate := 4.5      // float64
valid := true    // bool
```
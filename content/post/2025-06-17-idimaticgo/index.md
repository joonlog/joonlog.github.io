---
title: "Go 코딩 규칙(Idiomatic Go)"
date: 2025-06-17T09:00:00+09:00
categories: ["Go", "Go Basic"]
tags: ["go", "idiomatic go"]
---


### 변수 및 함수 이름

- camelCase 사용
- 축약어는 중간이나 끝에 (예: userID, httpServer)
- 간결하고 의미 있는 이름 권장 (`err`, `ok`, `i`, `r` 등)

```go
var userID int
func getUserInfo() {}
```

---

### **타입 및 인터페이스**

- 타입 이름은 PascalCase
- 인터페이스는 동작 기반 이름 사용 (`Reader`, `Closer` 등)

```go
type FileReader interface {
    Read(p []byte) (n int, err error)
}

type User struct {
    Name string
}
```

---

### **에러 처리**

- `if err != nil` 패턴이 기본
- 에러 발생 즉시 처리하고 빠르게 반환
- 에러 메시지는 소문자로 시작하고 마침표는 생략

```go
f, err := os.Open("file.txt")
if err != nil {
    return err
}
```

---

### **패키지 이름 규칙**

- 항상 소문자
- 단수형
- 의미는 명확하되 짧게 (`strings`, `math`, `http` 등)
- 밑줄(`_`), 복수형, 접두사 등은 지양

---

### **주석 스타일**

- 공개 함수나 타입에는 주석 필수
- 함수명으로 시작하는 문장형 주석 사용

```go
// Add는 두 숫자의 합 반환
/*
Add는 두 숫자의 합 반환
*/
func Add(a, b int) int {
    return a + b
}
```
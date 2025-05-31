---
title : Go Public / Private 식별자
date : 2025-05-31 09:00:00 +09:00
categories : [Go, Go Basic]
tags : [go, public, ] #소문자만 가능
---

> Go 언어는 `public`, `private` 키워드 없이 **이름의 첫 글자 대소문자**만으로 식별자의 공개 여부를 결정
함수, 변수, 인터페이스 등 **모든 식별자**에 적용되는 규칙
> 

## 공개(Public) 식별자와 비공개(Private) 식별자

- **대문자로 시작하는 이름** → 해당 패키지 외부에서도 접근 가능 (공개)
- **소문자로 시작하는 이름** → 해당 패키지 내부에서만 사용 가능 (비공개)
    
    ```go
    // 공개 함수
    func Reverse(s string) string {
        return reverseTwo(s)
    }
    
    // 비공개 함수
    func reverseTwo(s string) string {
        // 문자열 뒤집기
    }
    ```
    

## Public 함수 접근 예시

- 패키지 구조
    
    ```
    GolangTraining/
      ├── 02_library/
      │   └── stringutil/
      │       └── reverse.go
      └── 01_helloWorld/
          └── hello.go
    ```
    
- `reverse.go` (패키지 내부 함수 정의)
    
    ```go
    package stringutil
    
    func Reverse(s string) string {
        return reverseTwo(s)
    }
    
    func reverseTwo(s string) string {
        r := []rune(s)
        for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
            r[i], r[j] = r[j], r[i]
        }
        return string(r)
    }
    ```
    
- `hello.go` (외부 패키지에서 import)
    
    ```go
    package main
    
    import (
        "fmt"
        "github.com/goestoeleven/GolangTraining/02_library/stringutil"
    )
    
    func main() {
        fmt.Println(stringutil.Reverse("Hello, Go!"))
    }
    ```
    
- 실행 결과
    
    ```
    !oG ,olleH
    ```
    

## Private 함수에 접근하려 할 경우

- `reverseTwo`는 `stringutil` 내부에서만 사용 가능하므로, 아래와 같이 사용하면 컴파일 에러 발생
    
    ```go
    fmt.Println(stringutil.reverseTwo("Hello"))
    ```
    
- 컴파일 에러
    
    ```
    undefined: stringutil.reverseTwo
    ```
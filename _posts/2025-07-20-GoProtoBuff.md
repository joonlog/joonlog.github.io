---
title : Go 언어에서의 Protocal Buffers
date : 2025-07-20 09:00:00 +09:00
categories : [Go, Go Ecosystem]
tags : [go, protobuf, protocal buffers]  #소문자만 가능
---

- 언어 중립적, 플랫폼 중립적, 확장 가능한 데이터 직렬화(serialization) 포맷
    - gRPC의 기본 포맷
- 구조화된 데이터를 효율적으로 저장하고 전송하기 위한 방식
- 데이터 타입이 명확히 정의되어 있음 (`int32`, `string` 등)
- `.proto` 파일에서 데이터 구조를 정의하고, `protoc` 컴파일러로 각 언어로 변환해 사용
    - `protoc --go_out=. person.proto`

### JSON vs protobuf

- json보다 훨씬 적은 용량으로 전송 가능(binary 포맷)

```json
{ "name": "Alice", "age": 30 }
```

```protobuf
message Person {
  string name = 1;
  int32 age = 2;
  string email = 3;
}
```

- 필드 번호 기반으로 동작해서, 일부 필드가 빠져도 괜찮음 (`optional`처럼 동작)
    - 위에처럼 protobuf가 정의되어 있을 때 1, 2, 3이 각 필드의 번호
    - 아래처럼 데이터를 보낼 때, 필드 3이 없어도 동작
        - 이런 동작 방식 덕분에 새 필드를 추가해도 이전 클라이언트는 무시하고 정상 작동
    
    ```protobuf
    Person {
      name: "Alice"
      age: 30
      // email은 없음
    }
    ```
    

### Go에서의 protobuf

- protobuf를 Go 코드로 변환할 때의 추가 설정 옵션

```protobuf
string InitScript = 16 [(gogoproto.jsontag) = "init_script,omitempty"];
```

- `gogoproto`: gogo/protobuf라는 성능 최적화된 Protobuf 구현체에서 제공하는 옵션
- `jsontag`: 필드가 Go 구조체로 변환될 때, JSON 직렬화 시 사용될 `struct tag`를 직접 지정하는 설정입니다.
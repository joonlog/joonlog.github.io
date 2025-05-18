---
title : Ansible 조건문 / 반복문 / 핸들러 / 오류 처리
date : 2025-05-18 09:00:00 +09:00
categories : [Linux, Ansible]
tags : [linux, ansible, loop, if, handler, block rescue] #소문자만 가능
---

## 반복문과 조건문

1. 반복문
    - **단순 반복문**: `with_items`
    - **중첩 반복문**: `with_nested`
    - **기타 반복 지시문**:
        - `with_file`
        - `with_fileglob`
        - `with_sequence`
        - `with_random_choice`
    
    > 반복문에서 사용되는 변수는 List 형식으로 정의
    > 
    > 
    > 각 항목은 **key-value 형식**으로 사용 가능
    > 
    
    ```yaml
    - name: 단순 반복 예시
      debug:
        msg: "{{ item }}"
      with_items:
        - value1
        - value2
    ```
    
2. 조건문
    - `when` 구문 사용
    - 조건식은 **Jinja2 문법**을 사용하지만, **중괄호({})는 쓰지 않음**
    
    ```yaml
    - name: 조건문 예시
      debug:
        msg: "x86_64 시스템입니다."
      when: ansible_machine == "x86_64"
    ```
    
    - **자주 사용하는 조건 예시**:
        - `ansible_machine == "x86_64"`
        - `max_memory == 512`
        - `min_memory != 512`
        - `min_memory < 120`
        - `min_memory is defined`
        - `min_memory is not defined`
        - `memory_available`
        - `not memory_available`
        - `"ny special user" in superusers`
3. 반복문+조건문
    
    ```yaml
    - name: 특정 조건에 맞는 마운트 포인트 처리
      debug:
        msg: "{{ item.mount }} 마운트에 여유 공간 있음"
      with_items: "{{ mount_list }}"
      when: item.mount == "/" and item.size_available > 300000000
    ```
    

## 핸들러 구현

- `notify`는 해당 task가 **changed 상태일 때만** 트리거됨
- `handlers`는 **모든 task가 끝난 후** 실행됨

```yaml
tasks:
  - name: Task1
    file:
      src: httpd.conf.template
      dest: /etc/httpd/conf/httpd.conf
    notify:
      - restart httpd
handlers:
  - name: restart httpd
    service:
      name: httpd
      state: restarted
```

## 작업 오류 처리


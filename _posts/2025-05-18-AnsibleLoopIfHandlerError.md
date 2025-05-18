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

1. 실패 무시하기

```yaml
- name: 실패해도 무시
  command: /bin/false
  ignore_errors: yes
```

2. 실패 조건 지정

```yaml
- name: 실패 조건 직접 지정
  shell: some_command
  register: result
  failed_when: "'error' in result.stderr"
```

3. Changed 상태 지정

```yaml
- name: 강제로 Changed 처리
  shell: echo "Force change"
  changed_when: true
```

4. 핸들러 강제 실행

```yaml
- name: 전체 플레이북에 핸들러 강제 실행 설정
  hosts: all
  force_handlers: yes
```

5. Ansible 블록 및 오류 처리 - block ~ rescue ~ always

```yaml
- name: block-rescue-always 예시
  hosts: all
  tasks:
    - name: 중요한 작업 실행
      block:
        - name: 작업 1
          shell: /bin/false

        - name: 작업 2
          shell: echo "이건 실행 안됨"
      rescue:
        - name: 에러 발생 시 대체 작업
          debug:
            msg: "문제가 발생했지만 복구 작업 실행"
      always:
        - name: 항상 실행되는 작업
          debug:
            msg: "이 작업은 성공/실패 여부와 무관하게 항상 실행됨"
```
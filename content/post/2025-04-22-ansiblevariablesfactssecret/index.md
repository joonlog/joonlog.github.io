---
title: "Ansible Variables / Secrets / Facts"
date: 2025-04-22T09:00:00+09:00
categories: ["Ansible", "Basic"]
tags: ["ansible", "variables", "facts", "secrets"]
---


## 1. Variable

1. 변수 이름 지정
    - 변수 이름은 문자, 숫자 ,밑줄 가능
    - 변수 이름은 문자로 시작해야 한다.
2. 변수 적용 볌위
    - 변수 적용 범위는 좁은 범위가 넓은 범위보다 우선한다.
3. 플레이북에서 변수 정의
    - 변수 정의
        - 플레이북내에 변수를 지정하는 경우
            
            ```yaml
            # $ cat site.yaml
            -------------
            - hosts: web1
              vars:
                user: user01
                home: /home/user01
            ```
            
        - 별도의 파일로 변수를 설정하고 플레이북에서 지정하는 경우
            
            ```yaml
            # $ cat vars/users.yaml
            -------------
            user: user01
            home: /home/user01
            ---------------
            # $ cat site.yaml
            - hosts: web1
              var_files:
                - vars/users.yaml
            ```
            
    - 변수 사용
        - 변수의 사용은 “{{ variable }}” 형식으로 사용한다.
4. 변수의 종류
    - 호스트 변수
    - 그룹 변수
    
    [참고] 변수의 우선 순위
    
    - 0) 명령어 라인에서 지정하는 변수(ex: ansible -e ….)
    - 1) 플레이북에서 정의하는 변수(play.yaml(vars, var_files))
    - 2) 호스트 변수(host_vars/*)
    - 3) 그룹 변수(group_vars/*)
5. 배열 변수
    
    ```yaml
    # $ cat vars.yaml
    ----------------
    users:
      user01:
        name: lee         # => users.user01.name
        shell: /bin/bash  # => users['user01']['name']
      user02:
        name: kim
        shell: /bin/tsch
    ```
    
6. register 구문
    - register 구문을 사용하여 명령 출력을 캡쳐할 수 있다.
    - register 구문은 일반적으로 디버깅 용도로 사용한다
        - ansible-playbook -v …
    
    ```yaml
    tasks:
    - yum:
        name: httpd
        state: present
      register: VAR
    - debug:
        var: VAR
    ```
    

## 2. Secret

1. ansible-vault
    - 암호화된 새파일 만들기
    $ ansible-vault create secret.yaml
    - 암호화된 파일 보기
    $ ansible-vault view secret.yaml
    - 암호화된 파일 편집하기
    $ ansible-vault edit secret.yaml
    - 평문 파일 ↔ 암호화
    $ ansible-vault encrypt secret.yaml
    $ ansible-vault decrypt secret.yaml
    - 암호화 키 변경하기
    $ ansible-vault rekey secret.yaml
2. ansible-vault와 playbook
$ ansible-playbook --help | grep password
$ ansible-playbook --ask-vault-password playbook.yaml
$ ansible-playbook --vault-password-file-vault-pass playbook.yaml

## 3. Facts

앤서블 호스트에 대해 자동으로 모아지는 변수 내용

1. Facts
    - ansible_hostname
    - andible_fqdn
    - ansible_default_ipv4.address
    - …
2. Facts 조회 방법
    
    setup 모듈을 사용한 ansible ad-hoc 명령어로 facts 조회
    
    $ ansible localhost -m setup -a ‘filter=ansible_*’
    $ ansible localhost -m setup -a ‘filter=ansible_interface’
    
    ⇒ Facts 사용
    ”{{ansible_hostname}}”, “{{ansible_default_ipv4.address}}”
    ”{{ansible_facts['hostname']}}”, “{{ansible_default['ipv4.address']}}”
    
3. Facts 끄는 방법
    
    ```yaml
    # $ vi playbook.yaml
    ----------------
    - hosts: large_farm
      gather_facts: no
      tasks:
    	  ...
    ```
    
4. 사용자 정의 Facts
    - /etc/ansible/facts.d/*.fact (INI 형식 or JSON 형식)
    
    ```yaml
    {
      "packages": {
        "web_package": "httpd",
        "db_packages": "mariadb-server"
      },
      "users": {
        "user1": "user01",
        "user2": "user02"
      }
    }
    ```
    
    - 사용: playbook에서 사용하는 예
    # ansible localhost -m setup -a ‘filter=ansible_local’ 
    ⇒ ansible_local[’local’][’packages’][’web_package’]
    ⇒ ansible_local.local.packages.web_package
5. 매직변수(Magic Variable)
    - hostvars
    - group_names
    - groups
    - inventory_hostname
    - 매직변수 조회 방법
        - debug 모듈을 사용한 ad-hoc 명령어로 조회
        
        $ ansible localhost -m debug -a ‘var=hostvars’
        $ ansible localhost -m debug -a ‘var=group_names’
        $ ansible localhost -m debug -a ‘var=inventory_hostname’
        

## 4. include vars, include tasks 관리

1. include Overview
    
    작업 포함: include_tasks, import_tasks
    변수 포함: include_vars
    
2. include_tasks
    
    ```yaml
    # $ cat tasks/main.yaml
    ------
    - name: Task1
      yum:
        ...
        
    # $ cat play.yaml
    ------
    tasks:
    - include_tasks: tasks/main.yaml
    ```
    
3. include_vars
    
    ```yaml
    # $ cat include_vars
    ------
    web_pkg: httpd
    
    # $ cat play.yaml
    ------
    tasks:
    - include_vars: vars/vars.yaml
    ```

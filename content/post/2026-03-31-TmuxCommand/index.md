---
title: "Tmux 명령어"
date: 2026-03-31T09:00:00+09:00
categories: ["Session", "Tmux"]
tags: ["session", "tmux"]
---


> SSH가 끊겨도 터미널 세션과 실행 중인 작업을 서버에서 계속 유지시켜주는 세션 관리자
> 

### tmux란

- 하나의 SSH/터미널 안에서 여러 터미널 세션을 관리하는 **터미널 멀티플렉서**
- Claude Code 같이 긴 시간 세션을 유지해야 할 경우 효과적
- SSH가 끊겨도 세션이 계속 살아있고, 다시 같은 세션에 재접속 가능

### 명령어

- 새 세션 생성
    
    ```jsx
    tmux new -s 세션이름
    ```
    
- 세션 목록 확인
    
    ```jsx
    tmux ls
    ```
    
- 기존 세션 접속
    
    ```jsx
    tmux attach -t 세션이름
    ```
    
- 세션에서 빠져나오기
    - 종료 아님
    
    ```jsx
    Ctrl+b d
    ```
    
- 세션 종료
    
    ```jsx
    tmux kill-session -t 세션이름
    ```
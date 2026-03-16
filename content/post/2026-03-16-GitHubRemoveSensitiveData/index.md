---
title: "GitHub 내의 민감한 정보 제거"
date: 2026-03-16T09:00:00+09:00
categories: ["SCM", "Git"]
tags: ["git", "github", "github sensitive data"]
---


> 공개된 GitHub에 민감한 정보가 올라갔을 경우 히스토리에서 민감 정보 제거하는 방법
> 
1. git-filter-repo 설치
    
    ```jsx
    curl https://raw.githubusercontent.com/newren/git-filter-repo/main/git-filter-repo -o /usr/local/bin/git-filter-repo
    chmod +x /usr/local/bin/git-filter-repo
    ```
    
2. 교체할 문자열 파일 작성
    
    ```jsx
    echo "삭제할문자열==>대체할문자열" > /tmp/passwords.txt
    ```
    
3. mirror 리포에서 작업
    
    ```jsx
    git clone --mirror <repo url>
    cd <repo>.git
    ```
    
4. 히스토리 재작성
    
    ```jsx
    git filter-repo --replace-text /tmp/passwords.txt --force
    ```
    
5. origin 재등록
    - `git filter-repo`는 기본적으로 `origin`을 제거하기 때문에 재등록
    
    ```jsx
    git remote add origin <GITHUB URL>
    ```
    
6. Force 푸시
    
    ```jsx
    git push --force --mirror origin
    ```
    

---

주의사항

- 히스토리를 건드렸기 때문에 모든 커밋 해시가 변경됨
- 팀 작업 중이면 모든 팀원이 git fetch --all 후 로컬 브랜치 재설정 필요
- GitHub에 캐시되었을 때나 PR, Fork에 남아 있을 경우 하기 url을 통해 수정 필요
    - `https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository`
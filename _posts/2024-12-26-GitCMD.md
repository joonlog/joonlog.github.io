---
title : GitHub 명령어 정리
date : 2024-12-26 09:00:00 +09:00
categories : [SCM, Git]
tags : [git, github, github command] #소문자만 가능
---

- 깃허브 리포지토리 clone

```bash
git clone <repo-url>
```

- 깃허브 Pull / Push

```bash
git pull origin main 
git add .
git commit -am "commit message"
git push origin main
```

- 깃허브 rebase

```bash
# git pull 할 때 rebase 방식을 사용하도록 설정
# 로컬 브랜치의 변경 사항을 원격 브랜치의 끝에 위치
# 기본 설정은 merge
git config pull.rebase true

# 다른 브랜치의 커밋을 main 브랜치로 이동
# 커밋 히스토리를 정리하거나 브랜치를 재구성하는 데 사용
git rebase main

# rebase 작성 취소
git rebase --abort
```

- 깃허브 리포지토리 관리

```bash
git remote add origin <repo-url>
git remote set-url origin <repo-url>
```

- 깃허브 브랜치 변경

```bash
git show-ref
# default branch가 master로 되어있을 경우
git branch -m main
```

- 커밋 히스토리

```bash
git log --online
```

- 최근 커밋 취소

```bash
# 코드까지 돌아감
git reset --hard HEAD~2

# 코드는 그대로 두고 커밋만 삭제
git reset --soft HEAD~1
```
---
title : 깃허브 잔디 누락
date : 2024-12-27 09:00:00 +09:00
categories : [Git, GitHub]
tags : [git, github, github contribution] #소문자만 가능
---

커밋을 해도 깃허브에 잔디가 보이지 않아서 찾아보니 잔디가 심어지는 것에 조건이 있었다.

1. github 이메일 계정과 로컬의 이메일 정보가 같아야 한다
2. branch는 **`main`** 혹은 **`gh-pages`** 둘 중 하나에서 커밋해야 한다
3. forked repo가 아니어야 한다

내 경우에는 git config에 이메일이 등록되어 있지 않아서 잔디가 심어지지 않았다.

```bash
git config --global --list
git config --global user.email “<email>”
```

이메일 추가 후 해결

### 참고

공식문서: https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-github-profile/managing-contribution-settings-on-your-profile/why-are-my-contributions-not-showing-up-on-my-profile
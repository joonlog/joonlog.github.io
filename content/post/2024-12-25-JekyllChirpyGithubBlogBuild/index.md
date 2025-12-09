---
title: "GitHub 블로그 구축 - Jekyll Chirpy"
date: 2024-12-25T09:00:00+09:00
categories: ["GitHub Blog", "Jekyll"]
tags: ["git", "github", "github blog", "jekyll", "chirpy"]
---


### Jekyll Chirpy 테마

https://github.com/cotes2020/jekyll-theme-chirpy

## Jekyll Chirpy 사용 이유

- 기술 블로그용으로 가장 적합
    - 다른 Jekyll 테마들은 스타트업 블로그 스타일이나 쇼핑몰에서 사용할 것 같은 많은데, 개인용 기술블로그 테마로는 가장 깔끔하다고 느꼈다
    - 카테고리, 타임라인 별로 작성한 글을 볼 수 있고, 블로그 내 검색 기능이 빨라서 좋았다
    - 게시글 내에 있는 항목 별 바로가기 기능과, 게시글마다 읽는데 얼마나 걸리는지 나와있는 것도 맘에 들었다

![GitHubBlogBuild1.png](GitHubBlogBuild1.png)
    
- 많은 사용자
    - 커스터마이징 할 때 오류가 발생해도 사용자가 많으면 고치기가 수월하다
    - 특히 한국인 사용자가 많은 테마로 더 편한 면이 있다

## Jekyll Chirpy 구조

### **_config.yml**

- Jekyll 블로그의 기본 설정 파일
- 사이트 제목, 설명, URL, 언어, 플러그인 등 전반적인 설정을 정의

### **_data**

- 언어 설정, 사이드바 구성, 카테고리 또는 태그와 같은 반복적인 데이터 관리에 사용

### **_includes**

- 재사용이 가능한 header, footer, banner와 같은 HTML을 저장

### **_layouts**

- 페이지 레이아웃 템플릿을 저장
- `default.html`, `post.html` 등의 파일이 포함되어 각 페이지의 구조를 정의

### **_posts**

- 블로그 글을 작성하는 공간
- 각 파일은 `YYYY-MM-DD-title.md` 형식으로 저장
- Markdown 형식으로 작성하며, YAML front matter로 메타데이터(제목, 날짜, 태그, 카테고리)를 지정

### **_site**

- Jekyll 빌드 결과물이 저장되는 디렉터리
- 사용자가 직접 수정하지 않으며, GitHub Actions 또는 로컬 빌드 시 생성

### **_tabs**

- 블로그의 탭 네비게이션을 구성
- `about.md`, `archive.md`와 같은 파일을 포함
- 탭 메뉴에 표시되는 페이지의 내용과 경로를 정의

### **assets**

- 정적 파일(css, 이미지, JS 파일 등)을 저장
- `css/`, `img/`, `js/` 등의 하위 디렉터리가 포함

### **_sass**

- Sass 파일을 저장
- Chirpy 테마의 CSS를 수정하려면 이 디렉터리를 편집

### **.github**

- GitHub Actions 설정이 위치
- `workflows/` 디렉터리에는 Chirpy의 CI/CD 파이프라인 설정이 저장됨

## **깃허브 블로그 구축 과정**

### 블로그 Repo 생성

1. 저장소 생성
    - 깃허브에서 `username.github.io` Repo 생성
2. GitHub Pages 활성화
    - 저장소 설정 > Pages 섹션에서 브랜치 선택 후 활성화
3. 블로그 업로드
    - 작성 후 저장소에 Push → GitHub Pages가 자동으로 GitHub Actions 사용해 배포

### Jekyll - Chirpy 테마 적용

- Chirpy 복사
    - chirpy를 단순 fork하면 이후에 깃허브 잔디 심기가 곤란하니까 clone 후 copy

```bash
git clone git@github.com:joonlog/joonlog.github.io.git
git clone https://github.com/cotes2020/jekyll-theme-chirpy.git chirpy
cp -r chirpy/* joonlog
cp -r chirpy/.* joonlog
```

- Git 설정

```bash
cd joonlog

git remote set-url origin git@github.com:joonlog/joonlog.github.io.git
git pull origin main
```

### 로컬 테스트(선택)

- 의존성 설치

```bash
sudo dnf install ruby ruby-devel gcc make -y
sudo dnf groupinstall "Development Tools"
sudo gem install bundler

bundle install
npm install

bundle exec jekyll serve
  http://localhost:4000
```

### 깃허브 블로그 배포

- Git Push

```bash
git add .
git commit -am "test"
git push origin main
```

- Repo - 저장소 설정 - 페이지 섹션 - 배포 형식을 GitHub Actions로 변경
    - GitHub Actions 말고 기존 껄로 배포하면 index.html만 나타남
    
![GitHubBlogBuild2.png](GitHubBlogBuild2.png)
    
    - Configure 후 npm 의존성을 “Build with Jekyll” 단계 전에 추가 후에 Commit
        
        ```bash
        name: npm build
        run: npm install && npm run build
        ```
        
    - ./github/workflows/starter/pages-deploy.yml 삭제
        - 삭제 안하면 pages-deploy.yml도 같이 실행되어서 오류 발생
        

여기까지 하고 브라우저에서 username.github.io로 접속하면 블로그에 테마가 적용된 게 보일 것이다.

## 글 포스팅

| 구분 | 새 페이지 추가 | 새 게시물 추가 |
| --- | --- | --- |
| **목적** | 고정 콘텐츠(정적 페이지) 작성 | 블로그 글(동적 콘텐츠) 작성 |
| **위치** | 프로젝트 루트 또는 서브 디렉토리 | `_posts/` 디렉토리 |
| **파일명 규칙** | 자유로운 이름 (`about.md`) | `YYYY-MM-DD-title.md` 형식 |
| **URL** | `/about`, `/contact` 등 | `/year/month/day/title` 형식 |
| **목록 표시** | 목록에 포함되지 않음 | 블로그 목록에 자동 포함 |

- 형식 : `yyyy-mm-dd-제목.md`
- 확장자는 `.md` 또는 `.markdown`
- 중간에 공백을 넣지 않음
- 작성한 파일은 _posts 디렉토리에 위치

### **Front Matter**

- 사이트의 게시물에 대해 제목 및 레이아웃과 같은 변수와 메타데이터를 설정하기 위해 Markdown 또는 HTML 파일의 맨 위에 YAML 전문을 추가
- 파일명: `YYYY-MM-DD-title.md`
    
    ```yaml
    ---
    title : Rocky8 MySQL Replication & HAproxy 연동
    date : 2024-12-10 09:00:00 +09:00
    categories : [Linux, DB]
    tags : [rocky8, mysql, haproxy, replication] #소문자만 가능
    ---
    ```
    

블로그에 포스팅이 되었는지 확인

## 참고

공식문서: https://docs.github.com/ko/pages

jekyll chirpy 테마: https://github.com/cotes2020/jekyll-theme-chirpy
---
title : GitHub Blog(Jekyll Chirpy) SEO 설정
date : 2025-04-25 09:00:00 +09:00
categories : [GitHub Blog, Jekyll]
tags : [git, github, github blog, jekyll, jekyll chirpy, seo, sitemap] #소문자만 가능
---

### SEO

SEO는 Search Engine Optimization, 즉 검색 엔진 최적화로 구글이나 네이버 같은 검색 엔진에서 내 블로그가 더 잘 노출되도록 하는 작업이다.

블로그를 운영하더라도, 검색에 노출되지 않으면 방문자가 내 블로그 도메인을 직접 입력해야만 접근 가능하다. 사실상 필수 설정.

### GitHub Pages + Jekyll 환경 SEO 설정

1. `jekyll-seo-tag` 플러그인 설치
    - 블로그에 필요한 SEO 메타 태그들을 자동으로 생성해주는 공식 플러그인
    - `Gemfile`에 추가
    
    ```bash
    gem 'jekyll-seo-tag'
    ```
    
    - `_config.yml`에 추가
    
    ```bash
    plugins:
      - jekyll-seo-tag
    ```
    
2. `jekyll-sitemap` 플러그인 설치
    - 블로그 글 전체를 포함한 sitemap.xml 파일을 자동을 생성해주는 플러그인
    - 구글, 네이버 같은 검색엔진에 내 사이트 구조를 전달하는 목적의 파일
    
    ```bash
    gem 'jekyll-sitemap'
    ```
    
    ```bash
    plugins:
      - jekyll-sitemap
    ```
    
    - 배포 후 `https://<블로그주소>/sitemap.xml` 로 접근 가능하면 성공
3. Google Search Console 등록
    
    1) Google Search Console 접속
    
    2) `URL 접두어` 방식 선택 → `https://<블로그주소>` 입력
    
    ![GitHubBlogSEO1.png](/assets/img/git/githubblog/GitHubBlogSEO1.png)

    3) 생성된 verifications 값 _config.yml에 입력
    
    ```bash
    webmaster_verifications:
      google: "<값 입력>"
    ```
    
    4) 배포 후 Google Search Console 내 소유권 확인
    
    ![GitHubBlogSEO2.png](/assets/img/git/githubblog/GitHubBlogSEO2.png)
    
    5) 2번에서 생성된 사이트맵 Google Search Console에 등록
    
    - 등록 후에 즉시 활성화가 되는 것은 아님
        - 하루 정도 대기 필요

![GitHubBlogSEO3.png](/assets/img/git/githubblog/GitHubBlogSEO3.png)

### 참고

SEO: https://standing-o.github.io/posts/jekyll-seo/
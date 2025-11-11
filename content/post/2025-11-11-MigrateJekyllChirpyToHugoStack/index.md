---
title: "Jekyll Chirpy 테마에서 Hugo Stack 테마로 블로그 이관하기"
date: 2025-11-11T09:00:00+09:00
categories: ["GitHub Blog", "Hugo"]
tags: ["github blog", "hugo", "stack", "github pages", "go"]
---


> Jekyll Chirpy 테마로 포스팅했던 글들을 Hugo Stack으로 이관
> 

> 자세한 Hugo Stack 구축 과정은 이 글 참고
> 

### Jekyll Chirpy에서 Hugo Stack으로 이관한 이유

1. `Jekyll`은 `Liquid`를 사용해서, 포스팅에 `{%`, `{{` 같은 글자가 있으면 Liquid 템플릿 엔진이 이걸 Liquid 글로 인식을 해서 저 문자 자체로 사용하려면 `{% raw %}` 태그로 감싸야 했음
2. `Jekyll`은 `Ruby` 언어를 베이스로 짜여져 있고, 계층적 카테고리 구성과 같이 블로그를 커스터마이징할 일이 종종 있는데, 코드를 수정하기에는 `Ruby`에 대해 아는 바가 없음
3. 최근에 `Go`를 배우기 시작한 상황에서 `Go`를 베이스로 한 `Hugo`로 블로그를 구성하면 좋겠다 싶었음
4. 현재 Jekyll보다 `Hugo`가 더 점유율이 높고, 업데이트가 자주 있음
5. 직접 구성해보니 구조 자체도 `Hugo`가 훨씬 간단
6. 빌드부터 배포까지 `Hugo`가 더 빠름
7. 기존 `Chirpy` 테마와 비슷하면서 사용자가 많은 테마인 `Stack`을 선택

### 블로그 이관 과정

1. 기존 Jekyll 소스코드 별도 브랜치로 백업
    - Jekyll 브랜치 생성 후 백업
2. Hugo Stack Starter로 뼈대만 있는 블로그 배포
3. 기존 블로그 글 이관
    - 107개의 블로그 글 이관
    - 전체 이관 전에 글 1개로 테스트

### 이관 시 포스팅 글 변경점

- Front Matter
    - 날짜 형식
        - 첫 공백을 T로 변경하고, timezone 앞 공백 제거
        
        ```bash
        # Jekyll
        date: 2024-12-09 18:08:00 +09:00
        
        # Hugo
        date: 2024-12-09T18:08:00+09:00
        ```
        
    - Categories/Tags 배열 따옴표 추가
        
        ```bash
        # Jekyll
        categories: [Linux, Middleware]
        tags: [rocky8, apache, tomcat]
        
        # Hugo
        categories: ["Linux", "Middleware"]
        tags: ["rocky8", "apache", "tomcat"]
        ```
        
    - 폴더 구조
        - 각 포스팅마다 <포스팅 이름>.md가 아닌, <포스팅 이름> 폴더를 만들고 그 아래 index.md가 있는 구조
        
        ```bash
        # Jekyll
        _posts/2024-12-09-Rocky8ApacheTomcat.md
        
        # Hugo
        content/post/2024-12-09-rocky8apachetomcat/
          ├── index.md
          ├── Rocky8ApacheTomcat1.png
          └── Rocky8ApacheTomcat2.png
        ```
        
    - 이미지 경로 변경
        - Jekyll에서는 assests/img 경로에 모든 사진을 넣었지만, Hugo에서는 각 포스팅 폴더 안에 각각의 이미지를 넣음
        
        ```bash
        # Jekyll
        /img/linux/Rocky8ApacheTomcat1.png
        
        # Hugo
        Rocky8ApacheTomcat1.png
        ```
        

### 이후 예정 작업

- Hugo Theme에는 계층적 카테고리를 지원하는 테마가 없다
- 이 기능을 사용하기 위해선 별도 커스텀 필요
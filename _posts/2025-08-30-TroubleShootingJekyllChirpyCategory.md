---
title : Jekyll Chirpy 테마 카테고리 에러 TroubleShooting(1)
date : 2025-08-30 09:00:00 +09:00
categories : [GitHub Blog, Jekyll]
tags : [jekyll, chirpy, chirpy category error, troubleshooting] #소문자만 가능
---

### 문제 확인

깃허브 블로그 템플릿으로 Jekyll Chirpy 테마를 사용 중이었는데, 쓰다 보니 카테고리가 내가 설정한 것과 다르게 동작하는 것을 발견했다.

- EX 1) AWS를 상위 카테고리로 설정한 글을 작성했는데 실제 블로그에 AWS 카테고리가가 없는 경우
- EX 2) Kubernetes는 상위 카테고리로만 설정했었는데, Kubernetes 하위 카테고리로 Kubernetes가 생기고 이전 글들이 똑같이 복사돼서 들어가 있는 경우
- EX 3) 하위 카테고리로 AWS를 정했는데 AWS 상위 카테고리에 글이 들어가 있는 경우

### 원인

- Jekyll Chirpy의 `_layouts/categories.html` 템플릿에서 `site.categories`를 사용하는 방식의 문제
- `site.categories`의 동작 방식
    - `site.categories`는 모든 포스트의 모든 카테고리를 수집
    - `categories : [Certification, AWS]`인 글이 있으면 `site.categories['AWS']`에도 포함됨
    - 결과적으로 같은 글이 여러 상위 카테고리에 중복 표시
        
        ```html
        &#123;% for category in sort_categories %&#125;
         &#123;% assign category_name = category | first %&#125;
         &#123;% assign posts_of_category = category | last %&#125;
         &#123;% if category_name == first_post.categories[0] %&#125;
        ```
        

### 해결 방법

1. `_layouts/categories.html` 상위 카테고리 로직
    - `site.categories` 대신 `site.posts` 직접 순회
    - 첫 번째 카테고리(`post.categories[0]`)만 상위 카테고리로 수집
    - 해당 상위 카테고리를 가진 글들만 대상으로 처리
    
    ```html
    &#123;% assign group_index = 0 %&#125;
    
    &#123;% comment %&#125; Collect all primary categories from posts &#123;% endcomment %&#125;
    &#123;% assign primary_categories = '' | split: '' %&#125;
    &#123;% for post in site.posts %&#125;
      &#123;% assign primary_cat = post.categories[0] %&#125;
      &#123;% if primary_cat %&#125;
        &#123;% unless primary_categories contains primary_cat %&#125;
          &#123;% assign primary_categories = primary_categories | push: primary_cat %&#125;
        &#123;% endunless %&#125;
      &#123;% endif %&#125;
    &#123;% endfor %&#125;
    &#123;% assign primary_categories = primary_categories | sort %&#125;
    
    &#123;% for category_name in primary_categories %&#125;
      &#123;% comment %&#125; Get posts that have this as primary category &#123;% endcomment %&#125;
      &#123;% assign posts_of_category = '' | split: '' %&#125;
      &#123;% for post in site.posts %&#125;
        &#123;% if post.categories[0] == category_name %&#125;
          &#123;% assign posts_of_category = posts_of_category | push: post %&#125;
        &#123;% endif %&#125;
      &#123;% endfor %&#125;
    ```
    

2. `_layouts/categories.html` 하위 카테고리 로직

- 해당 상위 카테고리를 첫 번째로 가진 글들의 두 번째 카테고리만 수집
- 중복 제거 로직 유지
    
    ```html
    &#123;% comment %&#125; Collect subcategories &#123;% endcomment %&#125;
    &#123;% assign sub_categories = '' | split: '' %&#125;
    &#123;% for post in posts_of_category %&#125;
      &#123;% assign second_category = post.categories[1] %&#125;
      &#123;% if second_category %&#125;
        &#123;% unless sub_categories contains second_category %&#125;
          &#123;% assign sub_categories = sub_categories | push: second_category %&#125;
        &#123;% endunless %&#125;
      &#123;% endif %&#125;
    &#123;% endfor %&#125;
    ```
    

3. `_layouts/categories.html` 포스트 수 계산 로직

- `site.categories[category_name]` 대신 실제 해당 카테고리 글 개수 계산
- 하위 카테고리도 정확한 개수 표시
    
    ```html
    # Primary category
    &#123;% assign top_posts_size = posts_of_category | size %&#125;
    
    # Subcategory
    &#123;% comment %&#125; Count posts for this subcategory under current primary category &#123;% endcomment %&#125;
    &#123;% assign posts_size = 0 %&#125;
    &#123;% for post in posts_of_category %&#125;
      &#123;% if post.categories[1] == sub_category %&#125;
        &#123;% assign posts_size = posts_size | plus: 1 %&#125;
      &#123;% endif %&#125;
    &#123;% endfor %&#125;
    ```
    

4. 빌드 에러 발생 시

> Liquid syntax error (line 150): 'endif' is not a valid delimiter for for tags. use endfor
> 
- 불필요한 `&#123;% endif %&#125;` 제거 (for 루프에는 `&#123;% endfor %&#125;`만 필요)
    
    ```html
    &#123;% assign group_index = group_index | plus: 1 %&#125;
    &#123;% endfor %&#125;
    ```
    

### 해결

- 각 카테고리가 첫 번째 상위 카테고리 기준으로만 분류하도록 수정
- 중복 포스팅 문제 해결!
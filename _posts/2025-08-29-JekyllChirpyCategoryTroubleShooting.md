---
title : Jekyll Chirpy 테마 카테고리 중복/누락 TroubleShooting
date : 2025-08-29 09:00:00 +09:00
categories : [Blog, GitHub Blog]
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
        {% for category in sort_categories %}
         {% assign category_name = category | first %}
         {% assign posts_of_category = category | last %}
         {% if category_name == first_post.categories[0] %}
        ```
        

### 해결 방법

1. `_layouts/categories.html` 상위 카테고리 로직
    - `site.categories` 대신 `site.posts` 직접 순회
    - 첫 번째 카테고리(`post.categories[0]`)만 상위 카테고리로 수집
    - 해당 상위 카테고리를 가진 글들만 대상으로 처리
    
    ```html
    {% assign group_index = 0 %}
    
    {% comment %} Collect all primary categories from posts {% endcomment %}
    {% assign primary_categories = '' | split: '' %}
    {% for post in site.posts %}
      {% assign primary_cat = post.categories[0] %}
      {% if primary_cat %}
        {% unless primary_categories contains primary_cat %}
          {% assign primary_categories = primary_categories | push: primary_cat %}
        {% endunless %}
      {% endif %}
    {% endfor %}
    {% assign primary_categories = primary_categories | sort %}
    
    {% for category_name in primary_categories %}
      {% comment %} Get posts that have this as primary category {% endcomment %}
      {% assign posts_of_category = '' | split: '' %}
      {% for post in site.posts %}
        {% if post.categories[0] == category_name %}
          {% assign posts_of_category = posts_of_category | push: post %}
        {% endif %}
      {% endfor %}
    ```
    

2. `_layouts/categories.html` 하위 카테고리 로직

- 해당 상위 카테고리를 첫 번째로 가진 글들의 두 번째 카테고리만 수집
- 중복 제거 로직 유지
    
    ```html
    {% comment %} Collect subcategories {% endcomment %}
    {% assign sub_categories = '' | split: '' %}
    {% for post in posts_of_category %}
      {% assign second_category = post.categories[1] %}
      {% if second_category %}
        {% unless sub_categories contains second_category %}
          {% assign sub_categories = sub_categories | push: second_category %}
        {% endunless %}
      {% endif %}
    {% endfor %}
    ```
    

3. `_layouts/categories.html` 포스트 수 계산 로직

- `site.categories[category_name]` 대신 실제 해당 카테고리 글 개수 계산
- 하위 카테고리도 정확한 개수 표시
    
    ```html
    # 상위 카테고리
    {% assign top_posts_size = posts_of_category | size %}
    
    # 하위 카테고리
    {% comment %} Count posts for this subcategory under current primary category {% endcomment %}
    {% assign posts_size = 0 %}
    {% for post in posts_of_category %}
      {% if post.categories[1] == sub_category %}
        {% assign posts_size = posts_size | plus: 1 %}
      {% endif %}
    {% endfor %}
    ```
    

4. 빌드 에러 발생 시

> Liquid syntax error (line 150): 'endif' is not a valid delimiter for for tags. use endfor
> 
- 불필요한 `{% endif %}` 제거 (for 루프에는 `{% endfor %}`만 필요)
    
    ```html
    {% assign group_index = group_index | plus: 1 %}
    {% endfor %} 
    ```
    

### 해결

- 각 카테고리가 첫 번째 상위 카테고리 기준으로만 분류하도록 수정
- 중복 포스팅 문제 해결!
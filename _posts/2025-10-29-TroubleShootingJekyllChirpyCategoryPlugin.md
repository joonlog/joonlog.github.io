---
title : Jekyll Chirpy 테마 카테고리 에러 TroubleShooting(2)
date : 2025-10-29 18:00:00 +09:00
categories : [GitHub Blog, Jekyll]
tags : [jekyll, chirpy, chirpy category error, jekyll plugin, ruby, troubleshooting] #소문자만 가능
---

### 문제 재발견

- 이전에 [Jekyll Chirpy 테마 카테고리 중복/누락 TroubleShooting](https://joonlog.github.io/posts/troubleshooting-jekyll-chirpy-category/) 글에서 `_layouts/categories.html` Liquid 템플릿을 수정해서 카테고리 문제를 해결했다고 생각지만, 계속해서 버그가 확인됨
    - 카테고리 목록 페이지(`/categories/`)는 제대로 표시됨
    - 하지만 개별 카테고리 페이지(`/categories/kubernetes/`)에 들어가면 메인 카테고리와 서브 카테고리가 전부 섞여서 표시됨
    - 예시:
        - `Certification > Kubernetes` (서브카테고리): 1개 포스트
        - `Kubernetes > Architecture` (메인카테고리): 34개 포스트
        - `/categories/kubernetes/` 접속 시: 35개 포스트 전부 표시

### 근본 원인

- 이전 해결책은 `_layouts/categories.html` (목록 페이지)만 고친 것
- 진짜 문제는 Jekyll Archives 플러그인
1. Jekyll Archives 플러그인의 동작 방식
    - `categories: [Certification, Kubernetes]` 형태의 포스트가 있으면 `site.categories['Certification']`에도 포함되고 `site.categories['Kubernetes']`에도 포함됨
    - `/categories/kubernetes/` 와 같은 1단계 경로로만 페이지를 생성하고 `/categories/certification/kubernetes/` 같이 2단계 이상으로는 페이지를 생성하지 않음
2. `_layouts/category.html`의 한계

    ```
    {% raw %}{% for post in page.posts %}{% endraw %}
    ```

    - `page.posts`는 Jekyll Archives가 이미 생성한 데이터
    - Liquid 템플릿에서는 이미 섞인 데이터를 받음
    - 페이지 생성 자체는 Ruby 플러그인이 하기 때문에 Liquid로는 근본 해결 불가

### 해결 방법

- 커스텀 Ruby 플러그인을 작성해서 Jekyll Archives를 대체

### 1. Jekyll Archives의 categories 비활성화

- `_config.yml`
    
    ```yaml
    # 변경 전
    jekyll-archives:
      enabled: [categories, tags]
      layouts:
        category: category
        tag: tag
      permalinks:
        tag: /tags/:name/
        category: /categories/:name/
    
    # 변경 후
    jekyll-archives:
      enabled: [tags]
      layouts:
        tag: tag
      permalinks:
        tag: /tags/:name/
    ```
    

### 2. 카테고리 플러그인 작성

- `_plugins/hierarchical-categories.rb`
    
    ```ruby
    # frozen_string_literal: true
    
    module Jekyll
      # Generates hierarchical category pages
      # Primary categories: /categories/primary/
      # Secondary categories: /categories/primary/secondary/
      class HierarchicalCategoryPage < Page
        def initialize(site, base, primary, secondary = nil)
          @site = site
          @base = base
          @dir = if secondary
                   File.join('categories', primary.downcase.gsub(' ', '-'))
                 else
                   'categories'
                 end
          @name = if secondary
                    "#{secondary.downcase.gsub(' ', '-')}.html"
                  else
                    "#{primary.downcase.gsub(' ', '-')}.html"
                  end
    
          self.process(@name)
          self.read_yaml(File.join(base, '_layouts'), 'category.html')
    
          # Set category name for display
          self.data['title'] = secondary || primary
          self.data['category_name'] = secondary || primary
          self.data['primary_category'] = primary
          self.data['secondary_category'] = secondary
    
          # Filter posts based on category hierarchy
          if secondary
            # Secondary category: posts where categories[0] == primary AND categories[1] == secondary
            self.data['posts'] = site.posts.docs.select do |post|
              post.data['categories'] &&
              post.data['categories'][0] == primary &&
              post.data['categories'][1] == secondary
            end
          else
            # Primary category: posts where categories[0] == primary
            self.data['posts'] = site.posts.docs.select do |post|
              post.data['categories'] &&
              post.data['categories'][0] == primary
            end
          end
        end
      end
    
      class HierarchicalCategoryPageGenerator < Generator
        safe true
        priority :low
    
        def generate(site)
          return unless site.layouts.key? 'category'
    
          # Collect all primary and secondary categories
          primary_categories = {}
    
          site.posts.docs.each do |post|
            next unless post.data['categories']
    
            categories = post.data['categories']
            primary = categories[0]
            secondary = categories[1] if categories.length > 1
    
            next unless primary
    
            primary_categories[primary] ||= Set.new
            primary_categories[primary] << secondary if secondary
          end
    
          # Generate pages for each primary category
          primary_categories.each do |primary, secondaries|
            # Generate primary category page
            site.pages << HierarchicalCategoryPage.new(site, site.source, primary)
    
            # Generate secondary category pages
            secondaries.each do |secondary|
              next if secondary.nil?
              site.pages << HierarchicalCategoryPage.new(site, site.source, primary, secondary)
            end
          end
        end
      end
    end
    ```
    
- 핵심 로직:
    - 메인 카테고리: `post.categories[0]`만 매칭
    - 서브 카테고리: `post.categories[0] == primary AND post.categories[1] == secondary` 동시 매칭
    - URL 구조:
        - 메인 카테고리: `/categories/kubernetes/`
        - 서브 카테고리: `/categories/certification/kubernetes/`

### 3. 카테고리 목록 페이지 링크 수정

- `_layouts/categories.html`
    - 서브카테고리 URL을 `/categories/sub/`에서 `/categories/primary/sub/`로 수정
    - URL 마지막 `/` 제거 (Jekyll이 `.html` 파일을 디렉토리로 인식하는 문제 방지)
- 56줄 메인 카테고리 링크

    ```
    # 변경 전
    {% raw %}{% capture _category_url %}/categories/{{ category_name | slugify | url_encode }}/{% endcapture %}{% endraw %}

    # 변경 후
    {% raw %}{% capture _category_url %}/categories/{{ category_name | slugify | url_encode }}{% endcapture %}{% endraw %}
    ```
    
- 121줄 서브 카테고리 링크

    ```
    # 변경 전
    {% raw %}{% capture _sub_ctg_url %}/categories/{{ sub_category | slugify | url_encode }}/{% endcapture %}{% endraw %}

    # 변경 후
    {% raw %}{% capture _sub_ctg_url %}/categories/{{ category_name | slugify | url_encode }}/{{ sub_category | slugify | url_encode }}{% endcapture %}{% endraw %}
    ```
    

### 결과

- 메인 카테고리 페이지: `/categories/kubernetes/`
    - 메인이 카테고리가 kubernetes인 모든 포스트 표시
- 서브 카테고리 페이지: `/categories/certification/kubernetes/`
    - 메인 카테고리가 Certification, 서브 카테고리가 Kubernetes인 포스트만 표시
- 카테고리 목록 페이지에서 클릭 시 올바른 URL로 이동
- 카테고리 버그 해결!
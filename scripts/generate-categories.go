package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CategoryPair struct {
	Primary   string
	Secondary string
}

func main() {
	// 카테고리 맵: primary -> []secondary
	categories := make(map[string]map[string]bool)

	// content/post 디렉토리 스캔
	err := filepath.Walk("content/post", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// index.md 파일만 처리
		if !info.IsDir() && info.Name() == "index.md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", path, err)
				return nil
			}

			// front matter에서 categories 추출
			cats := extractCategories(string(content))
			if len(cats) >= 2 {
				primary := cats[0]
				secondary := cats[1]

				if categories[primary] == nil {
					categories[primary] = make(map[string]bool)
				}
				categories[primary][secondary] = true
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	// 기존 categories 디렉토리 삭제 및 재생성
	os.RemoveAll("content/page/categories")
	os.MkdirAll("content/page/categories", 0755)

	// 카테고리 페이지 생성
	for primary, secondaries := range categories {
		// Primary 카테고리 디렉토리 생성
		primaryDir := filepath.Join("content/page/categories", slugify(primary))
		os.MkdirAll(primaryDir, 0755)

		// Primary 카테고리 _index.md 생성
		primaryIndexPath := filepath.Join(primaryDir, "_index.md")
		primaryContent := fmt.Sprintf(`---
title: "%s"
description: "%s 카테고리의 모든 포스트"
primary_category: "%s"
layout: "category-primary"
---
`, primary, primary, primary)
		ioutil.WriteFile(primaryIndexPath, []byte(primaryContent), 0644)
		fmt.Printf("Created: %s\n", primaryIndexPath)

		// Secondary 카테고리 페이지 생성
		for secondary := range secondaries {
			secondaryPath := filepath.Join(primaryDir, slugify(secondary)+".md")
			secondaryContent := fmt.Sprintf(`---
title: "%s"
description: "%s > %s 카테고리의 포스트"
primary_category: "%s"
secondary_category: "%s"
layout: "category-secondary"
---
`, secondary, primary, secondary, primary, secondary)
			ioutil.WriteFile(secondaryPath, []byte(secondaryContent), 0644)
			fmt.Printf("Created: %s\n", secondaryPath)
		}
	}

	fmt.Printf("\n✓ Generated %d primary categories\n", len(categories))
}

// extractCategories는 front matter에서 categories를 추출
func extractCategories(content string) []string {
	// categories: ["Primary", "Secondary"] 형식 파싱
	re := regexp.MustCompile(`categories:\s*\[([^\]]+)\]`)
	matches := re.FindStringSubmatch(content)

	if len(matches) < 2 {
		return nil
	}

	// ["Primary", "Secondary"] -> [Primary, Secondary]
	categoriesStr := matches[1]
	categoriesStr = strings.ReplaceAll(categoriesStr, `"`, "")
	categoriesStr = strings.ReplaceAll(categoriesStr, `'`, "")

	parts := strings.Split(categoriesStr, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// slugify는 문자열을 URL에 적합한 형태로 변환
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// 한글은 그대로 유지, 특수문자만 제거
	re := regexp.MustCompile(`[^\w\-가-힣]+`)
	s = re.ReplaceAllString(s, "")
	return s
}

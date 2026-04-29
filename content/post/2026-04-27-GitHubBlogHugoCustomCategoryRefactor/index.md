---
title: "Hugo Stack 계층형 카테고리 리팩토링"
date: 2026-04-27T09:00:00+09:00
categories: ["GitHub Blog", "Hugo"]
tags: ["github blog", "hugo", "stack", "go", "hierarchical category", "refactoring"]
---


> 계층형 카테고리 커스텀 이후 발견한 버그 및 불필요한 코드 정리
>

이전 포스트에서 구현한 계층형 카테고리 시스템에서 3가지 문제를 발견

1. 포스트 `categories` frontmatter의 오타로 인해 카테고리가 잘못 생성되고 있었음
2. `generate-categories.go`에서 Go의 `slugify()`와 Hugo의 `urlize` 동작 불일치로 인한 잠재적 404 위험
3. 테마 사이드바 위젯 오버라이드 파일이 테마 기본값과 완전히 동일한데도 계속 유지되어, 테마 자동 업데이트 내용이 반영 안 되는 문제

---

## 1. 포스트 카테고리 오타 수정

- frontmatter의 `categories` 배열에 오타가 있으면 카테고리 페이지가 잘못된 이름으로 생성됨
- 발견된 오타 5건 수정

| 파일 | 변경 |
|------|------|
| `2025-11-04-AnsiblePlaybookDockerInstall` | `"Ansilbe"` → `"Ansible"` |
| `2025-11-05-AnsiblePlaybookcmdlogrotateopenfilesar` | `"Ansilbe"` → `"Ansible"` |
| `2026-02-04-MiniPCPromoxInstall` | `"Promox"` → `"Proxmox"` |
| `2026-02-23-PromoxNetworkVMSettings` | `"Promox"` → `"Proxmox"` |
| `2025-09-16-troubleshootinghaproxyrestartfault` | `"LoadBalancer"` → `"Load Balancer"` |

---

## 2. `generate-categories.go` 개선

### 2.1 ioutil 제거

- `io/ioutil` 패키지는 Go 1.16부터 deprecated
- `os.ReadFile` / `os.WriteFile`으로 직접 대체

```go
// 변경 전
import "io/ioutil"
content, err := ioutil.ReadFile(path)
ioutil.WriteFile(path, data, 0644)

// 변경 후
content, err := os.ReadFile(path)
os.WriteFile(path, data, 0644)
```

### 2.2 url 필드 추가

- **문제:** Go의 `slugify()` 결과와 Hugo의 `urlize` 함수 결과가 불일치하는 경우 발생 가능
    - 예: `"Load Balancer"` → `slugify()` → `"load-balancer"` vs Hugo `urlize` → `"load-balancer"` (일치하는 경우도 있지만, 특수문자·한글 혼용 시 불일치 가능)
    - 레이아웃 템플릿에서 `urlize`로 생성한 링크와 실제 생성된 페이지 URL이 달라지면 404
- **해결:** 생성되는 `_index.md`와 `.md` 파일에 `url:` 필드를 명시적으로 추가
    - Hugo가 `url:` 필드를 최우선으로 사용하므로 slugify 결과와 완전히 동일해짐

```go
// Primary _index.md
primarySlug := slugify(primary)
primaryContent := fmt.Sprintf(`---
title: "%s"
description: "%s 카테고리의 모든 포스트"
primary_category: "%s"
layout: "category-primary"
url: "/categories/%s/"
---
`, primary, primary, primary, primarySlug)

// Secondary .md
secondarySlug := slugify(secondary)
secondaryContent := fmt.Sprintf(`---
title: "%s"
description: "%s > %s 카테고리의 포스트"
primary_category: "%s"
secondary_category: "%s"
layout: "category-secondary"
url: "/categories/%s/%s/"
---
`, secondary, primary, secondary, primary, secondary, primarySlug, secondarySlug)
```

### 생성 결과 예시

```
content/categories/linux/_index.md
---
title: "Linux"
url: "/categories/linux/"
...
---

content/categories/linux/load-balancer.md
---
title: "Load Balancer"
url: "/categories/linux/load-balancer/"
...
---
```

---

## 3. 사이드바 위젯 오버라이드 삭제

- `layouts/partials/widget/categories.html`
- 테마 Stack v3의 기본 위젯 파일을 오버라이드하고 있었음
- 내용 비교 결과 테마 기본값과 **완전히 동일**
    - 유지할 이유가 없는 상태
- 매일 실행되는 테마 자동 업데이트 워크플로우가 있는데, 이 파일이 존재하면 테마에서 해당 위젯이 변경되어도 반영이 안 됨
    - 불필요한 오버라이드 파일은 삭제하는 것이 테마 독립성 유지에 적합
- 해당 파일 삭제

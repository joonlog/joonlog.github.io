# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Korean-language personal tech blog (joonlog.github.io) built with Hugo and the Hugo Theme Stack (v3). Hosted on GitHub Pages with automated deployment.

## Commands

```bash
# Local development server (with drafts)
hugo server -D

# Production build
hugo --minify --gc

# Generate category pages (must run before building when post categories change)
go run scripts/generate-categories.go

# Update theme to latest version
hugo mod get -u github.com/CaiJimmy/hugo-theme-stack/v3
hugo mod tidy
```

The dev server runs on port 1313.

## Architecture

### Category System

The blog uses a two-tier category hierarchy (`Primary > Secondary`), which is **not a built-in Hugo feature** — it's implemented via a custom pipeline:

1. **`scripts/generate-categories.go`** — scans all `content/post/*/index.md` files, reads the `categories: ["Primary", "Secondary"]` front matter, and regenerates the entire `content/categories/` directory. This script must be re-run whenever post categories change. The `content/categories/` directory is gitignored and generated at build time.

2. **Generated files** create two page types:
   - `content/categories/{primary-slug}/_index.md` — uses `layout: "category-primary"`
   - `content/categories/{primary-slug}/{secondary-slug}.md` — uses `layout: "category-secondary"`

3. **Custom layouts** in `layouts/_default/` render these pages:
   - `category-primary.html` — groups and displays all posts by secondary category
   - `category-secondary.html` — displays posts for a specific primary+secondary pair

### Content Structure

Posts follow the Hugo leaf bundle format:
```
content/post/{YYYY-MM-DD}-{PostTitle}/
├── index.md        # Post content with front matter
└── *.png/jpg       # Images referenced in the post
```

Front matter must include `categories: ["Primary", "Secondary"]` for category pages to work. Tags use a flat array: `tags: ["tag1", "tag2"]`.

### Configuration

All Hugo config is split across `config/_default/*.toml` files (Hugo's standard multi-file config pattern). The main files are:
- `config.toml` — base URL, language (ko-kr), pagination
- `params.toml` — theme settings (sidebar, comments via Giscus, color scheme)
- `markup.toml` — Goldmark renderer settings including LaTeX passthrough and code highlighting
- `module.toml` — Hugo module import for the theme

### CI/CD

Two GitHub Actions workflows:
- **`deploy.yml`** — triggered on push to `main`; runs `go run scripts/generate-categories.go`, then `hugo --minify --gc`, deploys to GitHub Pages
- **`update-theme.yml`** — daily cron at 00:00 UTC; auto-updates the Hugo Theme Stack module and commits the result

### Custom Additions Over Base Theme

- `layouts/partials/google_analytics.html` — Google Analytics integration
- `layouts/partials/widget/categories.html` — sidebar category widget
- `layouts/categories/list.html` — category listing page
- `assets/scss/custom.scss` — reserved for custom styles (currently empty)
- `static/robots.txt` and `static/googleb31e573b7a3eb431.html` — SEO/Search Console files

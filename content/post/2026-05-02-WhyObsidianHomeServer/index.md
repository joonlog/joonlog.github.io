---
title: "Notion에서 Obsidian으로 넘어온 이유"
date: 2026-05-02T09:00:00+09:00
categories: ["Home Server", "Obsidian"]
tags: ["home server", "obsidian", "notion", "syncthing", "ai"]
---


> Notion을 쓰다가 Obsidian으로 넘어온 이유
>
> 홈서버가 이미 있었기 때문에 가능한 구조

## 0. 배경

노션을 꽤 오래 써왔다. 메모, 자료 정리, 파일 첨부까지 하나의 앱에서 다 됐기 때문에 불편함이 없었다.

바꾸게 된 계기는 AI를 더 적극적으로 활용하고 싶어서였다. 내가 쌓아온 데이터에 AI가 언제든지 접근할 수 있었으면 했다.

## 1. Notion에서 AI 접근이 어려운 이유

Notion에 AI가 접근하려면 Notion API를 써야 한다.

- 페이지마다 API 호출이 따로 필요
- Rate limit: 초당 3건
- 수백 페이지를 실시간으로 뒤지는 건 현실적으로 불가능

Notion MCP도 결국 Notion API 위에 얹는 것이라 근본적인 제약은 동일하다.

원하는 건 단순했다. AI가 내 데이터에 언제든지, 비용 없이, 즉시 접근할 수 있는 구조.

## 2. Obsidian으로 넘어온 이유

Obsidian의 vault는 그냥 로컬 폴더다. 안에 있는 건 전부 마크다운 파일이다.

AI 입장에서는:
- 파일 직접 읽기
- 전체 검색
- API 없음, Rate limit 없음, 비용 없음

파일이 로컬에 있는 한 AI는 즉시, 제한 없이 접근할 수 있다.

부가적으로 Notion 무료 플랜의 파일 첨부 5MB 제한도 없어지고, 오프라인에서도 쓸 수 있다.

## 3. 홈서버가 있었기 때문에 가능했다

Obsidian vault가 로컬 파일이라는 건 장점이지만, 그러려면 파일이 항상 어딘가에 있어야 한다.

- AI가 언제든지 접근하려면 파일이 항상 살아있어야 함
- PC나 모바일은 꺼질 수 있음

홈서버에 mgt LXC가 이미 구축돼 있었다. 항상 켜져 있고, 내 네트워크 안에 있고, Claude Code가 직접 파일에 접근할 수 있다. 따로 뭔가를 더 구축할 필요 없이 vault를 여기 두면 됐다.

홈서버가 없었다면 이 구조를 생각하지 않았을 것 같다.

## 4. 구성한 구조

```
외부 PC / 모바일
  └─ Obsidian + Syncthing
       ↕ 동기화
mgt LXC (홈서버, 항상 켜짐)
  ├─ /root/obsidian (vault 폴더)
  ├─ Syncthing (동기화 허브)
  └─ Claude Code (vault 직접 접근)
```

Syncthing이 모든 기기와 mgt LXC를 동기화한다. PC나 모바일에서 메모를 쓰면 즉시 LXC에 반영되고, AI는 LXC의 파일을 직접 읽는다.

동기화 도구로 Syncthing을 선택한 이유는 단순하다.
- 무료, 오픈소스
- E2E 암호화
- 홈서버가 허브 역할을 하기 때문에 외부 릴레이 서버 없이 직접 연결 가능
- 기기 추가는 Device ID 교환으로 간단히 확장

실제 Syncthing 설정 방법은 다음 글에서 다룬다.

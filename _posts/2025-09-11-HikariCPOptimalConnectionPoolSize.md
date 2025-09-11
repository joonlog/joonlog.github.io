---
title : HikariCP 개발자가 제안한 커넥션 풀 사이즈 공식
date : 2025-09-11 09:00:00 +09:00
categories : [Linux, DB]
tags : [linux, spring, hikaricp, mysql, pool size] #소문자만 가능
---

> https://joonlog.github.io/posts/TroubleShootingHikariCPConnectionError/에서 이슈가 됐던 설정 중 최적 커넥션 풀 설정 방법에 대한 글 발견
> 

출처: https://www.threads.com/@codevillains/post/DFj2-wQSREJ?xmt=AQF0jSXFspf2rIJ-PoZij1A2Xy3Ak_g0vjc0SNxC2kYoOA

원본: https://github.com/brettwooldridge/HikariCP

1. 서버 스펙으로 계산
    - Pool Size = (( CPU 코어 수 * 2) + 스토리지 디스크 (또는 SSD) 갯수)
    - 코어가 8개고 스토리지 디스크가 2개라면 ((8*2)+2) = 18 가 최적의 커넥션 풀 사이즈
2. TPS로 계산
    - (TPS × 평균 응답 시간) / (1 - DB 사용률)
    - 500TPS에 평균 응답시간이 0.05초 일 경우 서버 자원 사용률이 70퍼라고 가정하면 500 * 0.05 / 1-0.7 이 되고, 83.3 개가 최적의 커넥션 풀 사이즈
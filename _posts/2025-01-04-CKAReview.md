---
title : CKA 후기
date : 2025-01-04 09:00:00 +09:00
categories : [Certification, Kubernetes]
tags : [kubernetes, k8s, cka] #소문자만 가능
---

![CKAReview1.png](/assets/img/certification/CKAReview1.png)

자격증들은 보통 깊은 기술을 요구하기 보단 넓고 얕게 아는 걸 요구하기 때문에 자격증 취득은 그 기술에 대해 어느 정도 알아보는 데 도움이 된다고 생각한다. 관련해서 쿠버네티스도 CKA 자격증 취득을 쿠버네티스 공부의 첫 목표로 잡았다. 

CKA는 시험에 불합격 하더라도 retake 기회를 1회 제공하고, 자격증 유효 기간은 내가 합격 했을 땐 3년이었지만 지금은 2년으로 알고 있다. 

### 할인 코드

Linux Foundation은 40% 할인 쿠폰을 1년에 두세번 정도 배포한다. 40%가 아니더라도 30%정도는 상시 할인하는 느낌이니 쿠폰은 반드시 적용해야 한다. 

### 공부 방법 및 기간

CKA를 공부했다면 모두가 아는 mumshad 선생님의 udemy 강의(https://www.udemy.com/course/certified-kubernetes-administrator-with-practice-tests/?couponCode=KEEPLEARNING)를 통해 기초를 잡았다. 외에도 따배쿠 유튜브도 많이 보는 것 같다. 

CKA를 신청하면 Linux Foundation에서 killer.sh 문제를 풀어 볼 수 있게 제공했던 걸로 기억하는데, 다른 블로그에서 하나같이 시험 난이도에 비해 killer.sh가 훨씬 어렵다고 해서 한번도 풀어보지 않았다. 대신 https://www.itexams.com/exam/CKA 여기 있는 문제들과 다른 블로그들을 참고해서 연습용 문제들을 골랐다.

공부 기간은 3주동안 매일 4시간 정도였는데, 그때는 리눅스를 약간 다룰 수 있는 정도였고 K8S는 하나도 몰랐으니 K8S를 다뤄 봤다면 더 짧아질 수 있을 것이다.

### 연습 환경

문제 풀이 실습들은 https://killercoda.com/cka에서 하는 걸 적극 추천한다. CKA는 쿠버네티스 구축에 대한 역량은 고려하지 않기 때문에 VM같은 걸로 K8S 구축하는데 에너지를 쏟기 보다 killercoda를 사용하는 게 좋다. 

이 사이트는 K8S를 테스트 하는데 최적화되어 있는데, 무료는 1개 세션에서 클러스터 1개와 노드 1개를 1시간 동안 사용할 수 있지만 이정도면 충분하다. 1시간이 지나더라도 그 세션만 초기화 되고 다시 사용할 수 있기 때문에 한문제 풀고 초기화하고 또 풀면 된다. 

### 시험 환경

CKA는 정해진 시험장이 없기 때문에 보통 노트북을 가지고 시험을 본다. 노트북은 캠과 오디오가 동작해야 한다. 시험 신청을 하면 안내 메일이 오는데, 여기에 시험 환경을 테스트 할 수 있는 링크가 있다. 그 링크에서 캠/오디오 등등을 체크 해야 시험 당일에 곤란하지 않을 것이다. 또 감독관 쪽에서 내 노트북을 원격으로 감시해야 하는데, 혹시 몰라서 윈도우 방화벽도 내렸다. 

시험 시작 30분 전부터 시험 환경인 PSI 브라우저를 설치할 수 있다. 접속하면 감독관이 본인 확인과 방 안 검사를 요구하는데, 이 시간 안에 감독관한테 ok가 나와야 이후 시험에 지장이 가지 않는다. 본인 확인은 **여권**을 지참해야 한다.

정해진 시험장이 없는 만큼 주위 환경도 깨끗해야 한다. 노트북을 제외한 어떤 것도 시험장에 있으면 안되고, 벽에 안내문 같은 종이 붙어 있다면 어떻게든 가려야 한다. 나는 집에서는 이 환경을 구성하기가 어려울 것 같아, 공유 오피스에서 최대한 조용한 구석으로 자리를 잡았다.

### 시험 팁

1. 모든 문제는 그 문제만의 k8s context가 명시되어 있다. 따라서 문제 풀기 전에 반드시 context 체크부터 해야 한다. 
2. PSI 브라우저에서 https://kubernetes.io/docs/home/ 공식 문서에 접근이 가능하다. 즉 오픈북 시험이란 것. 대부분의 문제에 대한 구축 과정을 공식 문서에서 확인이 가능하니 그 문서를 빨리 찾을 수 있게 연습하는 것도 좋을 것이다. 특히 K8S Upgrade나 ETCD Backup 문제 같은 경우 공식 문서를 참조하는 걸 추천한다.

### 시험 변경점

문제들은 처음 나왔을 때부터 지금까지 변경이 거의 없는 것으로 알고 있는데, 올해 1월 15부터 변경이 있다고 한다. 

https://www.youtube.com/watch?v=fvvgM3QmKGo 

https://www.inflearn.com/notices/1375488

실무에서 사용하는 기술들을 많이 반영한 것 같은데, 대표적으로

1. Helm 차트
2. Kustomize
3. CRD

등등이 추가되고 기존 문제도 강화된다고 한다.
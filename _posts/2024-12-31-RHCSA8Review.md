---
title : RHCSA 8 후기
date : 2024-12-31 09:00:00 +09:00
categories : [Certification, RedHat]
tags : [redhat, linux, rhcsa8, rhel8] #소문자만 가능
---

![RHCSA8Review1.png](/assets/img/certification/RHCSA8Review1.png)

작년에 본 시험에 대해 이제야 후기를 남긴다

당시 나는 이제 막 인프라로 진로를 튼 상황이었기 때문에 첫 목표를 RHCSA 취득으로 정했다. 

내 리눅스 실력은 기본 명령어만 약간 아는 정도였고 서비스 구축은 할 줄 몰랐다. 

udemy 강의(https://www.udemy.com/course/unofficial-linux-redhat-certified-system-administrator-rhcsa-8/?couponCode=KEEPLEARNING)를 들으면서 기본 지식을 쌓았고, 문제들은 다른 블로그들을 전부 뒤져봤다.

영어로 되어 있지만 요즘은 크롬 확장 프로그램 trancy 사용하면 한글로 더 편하게 볼 수 있을 것이다.

RHCSA가 절대적인 난이도로 보면 취득한 자격증 중에 가장 난이도가 낮지만, 난 이 시험이 가장 어려웠다. 1번 탈락하여 재시험 보기도 했고.

원래는 재시험 기회가 없었지만, 2023년부터 정책이 바뀐건지 retake 기회가 1번 주어졌고, 여기서 합격할 수 있었다.

자격증 유효 기간은 3년이다.

### 할인 쿠폰

시험 비용이 꽤 비싼 편인데, 할인 쿠폰도 많이 없고 할인율도 15%로 크진 않다.

이 쿠폰은 RedHat 시험에 합격하면 나오는 쿠폰으로 1달? 3달? 동안 3명까지 쓸 수 있다. 그래서 종종 사람들이 커뮤니티에 공유하기도 하는데, 나는 redhat reddit에서 주워서 사용했다. 

redhat discount reddit으로 검색하면 쿠폰 공유하는 게시물을 찾을 수 있을 것이다.

### 시험 환경

**여권 준비 필수**

시험 장소는 에티버스러닝이었는데, 칸막이가 쳐진 작은 방에 2명이 시험 볼 수 있게 되어 있었다. 입장해서 준비된 노트북에 로그인을 하면 감독관이 채팅으로 말을 걸어온다. 준비된 웹캠에 여권을 보여주고 웹캠을 가지고 방 안을 꼼꼼히 살펴보도록 요구한다. 

이후 시험을 시작하고 VM으로 된 노드 2개를 가지고 작업하면 된다.

### 유의사항

모든 문제 풀이들은 재부팅 후에도 동작해야 한다. 

예를 들어 httpd 서비스를 기동했다면, systemctl start httpd만 해선 안되고 반드시 enable 까지 해줘야 한다. 그래서 start와 enable을 같이 해주는 systemctl enable --now httpd를 사용하면 좋다.

서비스 이외에도 디스크 마운트, podman 문제도 마찬가지다. 따라서 문제를 다 풀었다면 노드들을 재부팅하고 잘 동작하는지 확인해야 한다.

### 출제 문제

RHCSA : [https://www.redhat.com/ko/services/training/ex200-red-hat-certified-system-administrator-rhcsa-exam?section=개요](https://www.redhat.com/ko/services/training/ex200-red-hat-certified-system-administrator-rhcsa-exam?section=%EA%B0%9C%EC%9A%94)

작년에는 RHCSA 8버전으로 시험을 봤었고 지금은 아마 9버전으로만 볼 수 것이다. 그때가 8에서 9로 넘어가는 과도기였는데 9버전 후기들로 파악할 수 있는 문제가 많지 않아 8버전으로 시험을 봤었다. 문제가 크게 바뀌진 않은 것 같지만 2~3문제정도에서 사용하는 서비스가 아예 달랐던 것으로 기억한다.

구버전 문제기도 하고, 문제 풀이들은 다른 블로그에도 많으니 어떤 문제가 나왔는지 정도만 공유하겠다.

루트 비밀번호 초기화

- 부팅 프로세스 rd.break로 인터럽트
- /sysroot 마운트 후 /sysroot 루트로 설정 후 비번 재설정

yum repository 구성

- 문제에서 주어지는 baseurl을 가지고 repo파일 생성
- dnf repolist all로 설정 됐는지 확인

ip, gateway, dns등 네트워크 인터페이스 구성

- nmtui 사용
- NetworkManager 서비스 재시작
- nmcli con show, nmcli dev status 같은 명령어로 잘 설정 됐는지 확인

hostname 설정

cronjob 생성

user/group 관리 명령어

- user에 group 추가, user nologin 설정

파일/디렉토리 권한 명령어

- 파일/디렉토리 user:group 권한 설정
- setgid 같은 파일/디렉토리 특수권한 설정

chrony 설정

- timedatectl 사용
- chrony.conf 설정

grep 후 리다이렉션 명령어

selinux 켜져있고 82번 포트로 httpd 동작중일 때 오류 해결

- semanage port 명령어로 http프로토콜 82포트 허용
- 방화벽에도 82포트 등록
- httpd, firewall 재시작

디스크 생성 후 마운트

- fdisk로 디스크 붙이고 vg, lv, ext4 초기화, 마운트 작업

swap 파티션 만들고, 재부팅하면 자동으로 on 되게 설정

- fdisk, fstab, swapon 작업

~~vdo 구축~~

- RHCSA 9 에는 안나오는 걸로 알고 있다.

~~autofs 구축~~

- RHCSA 9 에는 안나오는 걸로 알고 있다.

podman 루트리스  방식으로 호스트의 디렉터리와 연결된 컨테이너를 재부팅 후 자동으로 실행되게 설정

쉘 스크립트 파일 작성

### 참고 블로그:

https://not-to-be-reset.tistory.com/491

https://blog.naver.com/PostView.naver?blogId=asd7005201&logNo=222412331243&parentCategoryNo=&categoryNo=15&viewDate=&isShowPopularPosts=true&from=search

https://whitestudy.tistory.com/83

https://blog.naver.com/moonlitare/222418066071
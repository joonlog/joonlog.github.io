---
title : RHCSA 8 후기
date : 2024-12-31 09:00:00 +09:00
categories : [Certification, RHCSA 8]
tags : [redhat, linux, rhcsa8, rhel8] #소문자만 가능
---

![image.png](https://prod-files-secure.s3.us-west-2.amazonaws.com/8db24459-7346-4d70-afa5-8fa6c279412c/12122efc-1a0e-4ddd-8e93-70c07b825b00/image.png)

작년에 본 시험에 대해 이제야 후기를 남긴다

당시 나는 이제 막 인프라로 진로를 튼 상황이었기 때문에 첫 목표를 RHCSA 취득으로 정했다. 

내 리눅스 실력은 기본 명령어만 약간 아는 정도였고 서비스 구축은 할 줄 몰랐다. 

udemy 강의를 들으면서 기본 지식을 쌓았고, 문제들은 다른 블로그들을 전부 뒤져봤다.

강의는 https://www.udemy.com/course/unofficial-linux-redhat-certified-system-administrator-rhcsa-8/?couponCode=KEEPLEARNING 이걸 들었다.

영어로 되어 있지만 요즘은 크롬 확장 프로그램 trancy 사용하면 한글로 더 편하게 볼 수 있을 것이다.

RHCSA가 절대적인 난이도로 보면 취득한 자격증 중에 가장 난이도가 낮지만, 난 이 시험이 가장 어려웠다. 1번 탈락하여 재시험 보기도 했고.

원래는 재시험 기회가 없었지만, 2023년부터 정책이 바뀐건지 retake 기회가 1번 주어졌고, 여기서 합격할 수 있었다.

### 필수 준비물: 여권

### 시험 범위

RHCSA : ![RHCSA8Review1.png](/assets/img/linux/RHCSA8Review1.png)

작년에는 RHCSA 8버전으로 시험을 봤었고, 지금은 9버전일 것이다. 

그래도 크게 바뀌진 않았을 테니 생각나는 문제들을 적자면

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
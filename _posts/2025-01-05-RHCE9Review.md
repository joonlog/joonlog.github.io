---
title : RHCE9 후기
date : 2025-01-05 09:00:00 +09:00
categories : [Certification, RedHat]
tags : [redhat, linux, rhce, rhce9, ansible] #소문자만 가능
---

![RHCE9Review1.png](/assets/img/certification/RHCE9Review1.png)

올해 초에 들은 클라우드 교육 과정 중 Ansible이 있었는데, 기왕 RHCSA를 가지고 있는 거 RHCE까지 따기로 했다. 아마 RHCSA를 취득해놓은 상태가 아니라면 굳이 RHCE까지는 따려고 하지 않았을 것 같다. RHCSA와 마찬가지로 1번의 retake 기회가 있고, 자격증 유효 기간은 3년이다.

### 할인 쿠폰

시험 비용이 꽤 비싼 편인데, 할인 쿠폰도 많이 없고 할인율도 15%로 크진 않다.

이 쿠폰은 RedHat 시험에 합격하면 나오는 쿠폰으로 1달? 3달? 동안 3명까지 쓸 수 있다. 그래서 종종 사람들이 커뮤니티에 공유하기도 하는데, 나는 redhat reddit에서 주워서 사용했다. 

redhat discount reddit으로 검색하면 쿠폰 공유하는 게시물을 찾을 수 있을 것이다.

### 공부 방법 및 기간

나에게는 RHCE가 RHCSA보다 쉽게 느껴졌다. 문제 자체도 더 심화된 리눅스 문제가 나오는 게 아닌, RHCSA 문제에 Ansible 문법만 더한 정도였다. 만약 RHCSA를 가지고 있고, 어떤 교육이든 강의든 Ansible을 배워 봤다면 쉽게 취득할 수 있을 거라고 생각한다. 하루 4시간씩 2주 정도 공부했던 것 같다.

### 시험 환경

**여권 준비 필수**

시험 장소는 에티버스러닝이었는데, 칸막이가 쳐진 작은 방에 2명이 시험 볼 수 있게 되어 있었다. 입장해서 준비된 노트북에 로그인을 하면 감독관이 채팅으로 말을 걸어온다. 감독관 요구에 따라 웹캠에 여권을 보여주고 방 안을 검사 맡으면 된다.

시험 시간은 4시간으로 매우 긴 편이다. ansible로 여러 노드를 제어해야 되게 때문이다. 연습을 충분히 했다면, 꼼꼼하게 풀고 에러도 거의 없을 때를 기준으로 2시간 반 정도면 전부 풀 수 있다.

### 시험 팁

1. **모든 문제 풀이들은 재부팅 후에도 동작해야 한다.** 
    
    채점할 때 모든 노드를 재부팅 후 채점하기 때문에 service가 enable 되어 있는지, 디스크가 마운트되어 있는지, 파일들이 정상적으로 생성되었는지 잘 체크해야 한다.
    
2. **doc를 보면서 풀어라.**
    
    Ansible은 설치 시에 docs가 내장되어 있다. 문제 풀기 전 내가 사용할 모듈의 docs를 확인하는 습관을 들인다면 오류를 줄일 수 있을 것 이다. 나는 모든 문제를 doc를 열고 보면서 풀었다.
    
    - `ansible-doc <모듈명>`으로 문서를 열고 `/EXAMPLES`를 통해 예시 검색
3. **alias 지정이 시간 단축이 도움이 된다.**
    
    자주 입력해야 하는 `ansible`, `ansible-playbook`, `ansible-playbook --syntax-check`, `ansible-doc` 등의 명령어들을  `~/.bashrc`에 `alias ans='ansible'`, `alias anp=’ansible-playbook’`과 같이 지정해 놓으면 좋다.
    
4. **ansible-navigator 사용 가능하다.**
    
    RHCE9 버전부터 사용 가능한 걸로 알고 있다. 내 경우엔 ansible-navigator로 배우지 않아 navigator가 더 어색해서 쓰지 않았지만, 처음부터 navigator을 쓰는 걸로 배웠다면 tui 환경이니 더 직관적으로 쉬울 것 같다.
    
5. **변수는 inventory에 직접 지정하는 게 더 빠를 때가 있다.**
    
    문제에서 ip와 같은 각 노드마다 고정된 변수를 가지고 ansible 작업을 해야할 때, 팩트와 매직 변수에서 필요한 변수를 찾기 어렵다면, 노가다식으로 inventory 각 노드에 변수를 직접 지정하는게 빠르다. 
    
    어떻게든 문제를 푸는 게 더 중요하다.
    
6. 팩트, 매직 변수에서 필요한 정보를 어떻게 찾는지 연습해야 한다.
    
    팩트 검색 예시: `ansible localhost -m setup -a 'filter=ansible_devices'`
    
    매직변수 검색 예시: `ansible localhost -m debug -a 'var=hostvars["ansible1.example.com"]'`
    

### 출제 문제

1. Ansible 설치 및 구성
2. Repository 구성
3. 컬렉션 설치
4. 패키지 설치
5. rhel-system-roles 사용
6. Role 생성 및 사용
7. Ansible Galaxy로 Role 설치
8. Ansible Galaxy로 Role 사용
9. 파티션 생성 및 사용
    - when 조건 사용
10. hosts 파일 생성
    - 매직변수 사용
11. 파일 배포
    - inventory 변수 사용
12. 웹 컨텐츠 디렉터리 배포
    - 파일 생성 후 setype이 httpd_sys_content_t와 같이 웹서버 디렉터리 맞는지 확인
13. 하드웨어 리포트 배포
    - 팩트 사용
14. ansible vault 생성
15. 사용자 계정 생성
    - 위에서 생성한 vault 사용
    - 패스워드에 password_hash 걸었는지 체크
16. vault 키 변경
17. cronjob 생성
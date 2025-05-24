---
title : NCP Backup
date : 2025-05-24 09:00:00 +09:00
categories : [NCP, Storage]
tags : [ncp, storage, backup] #소문자만 가능
---

- NCP의 Backup 서비스는 서버 내 파일/DB를 백업할 수 있는 상품
- Backup 설정은 5단계로 구성
    1. **Resource**: 백업 대상 서버 등록
    2. **Storage**: 백업 데이터 저장소
    3. **Policy**: 보관 정책 정의
    4. **Job**: 백업 작업 정의
    5. **Schedule**: 백업 실행 일정

---

### Resource

![NCPBackup1.png](/assets/img/ncp/storage/NCPBackup1.png)

- 존, 서버, 에이전트 유형(Data, DB) 선택
    - 아이디/비밀번호는 서버 root/root비밀번호 입력
    - Resource 설정 시 위 정보가 정상적이라면 자동으로 Agent 설치

> 주의) DB 백업을 위한 에이전트 추가 이후 `DB 인스턴스 추가`에 DB 버전 제한이 있으니 확인 필수
> 
- `DB 인스턴스 추가`에 실패해도 별도 에러 메세지 없어서 콘솔 상 실패 이유 확인 불가
    - `/var/log/commvault/Log_Files/MySqlBrowseAgent.log`
    - DB 인스턴스 추가 시 로그는 위 로그에 기록됨
    - https://guide.ncloud-docs.com/docs/backup-spec
    
    ![NCPBackup2.png](/assets/img/ncp/storage/NCPBackup2.png)
    

---

### Storage

- DB 백업용 스토리지
    - 스토리지에 별도 접근 불가
- Policy 생성 전 Storage 생성 필수
    - Storage 생성하지 안고 Policy 생성 시 별다른 에러 메세지 없이 실패

---

### Policy

![NCPBackup3.png](/assets/img/ncp/storage/NCPBackup3.png)

- 백업을 보관할 Storage, `보관 기간` 설정

---

### Job

![NCPBackup4.png](/assets/img/ncp/storage/NCPBackup4.png)

- Data 백업 시에 백업 대상 경로 선택

![NCPBackup5.png](/assets/img/ncp/storage/NCPBackup5.png)

- DB 백업 시에 백업 DBMS, 개별 백업 DB 선택
    - DB 백업은 mysqldump를 원격으로 하는 형태

---

### Schedule

![NCPBackup6.png](/assets/img/ncp/storage/NCPBackup6.png)

- 백업 주기, 백업 시간 설정
- 증분 백업과 전체 백업이 겹치면 오류
    - 따라서 하루 전체 백업, 6일 증분 백업으로 보통 정책 7개 이상 생성
    - 주로 ncp cli로 크론탭에 스크립트 만들어 사용
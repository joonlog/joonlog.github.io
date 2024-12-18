---
title : Rocky8 Firewall
date : 2024-12-18 09:00:00 +09:00
categories : [Linux, Firewall]
tags : [rocky8, firewall, firewalld] #소문자만 가능
---

### 목표

- 특정 포트를 포워딩하여 외부 클라이언트가 내부 리소스에 접근 가능하도록 구성
- 특정 웹사이트에 대한 외부 접속을 방화벽 규칙을 통해 차단

### 환경 설정

- KT Cloud
    - Firewall 2대
        - 공인 IP 포트포워딩: 2222→22, 8080→80
        - 공인 IP 포트포워딩: 2223→22, 8081→80
- OS: Rocky Linux 8.1

## 개념 및 설정

### FirewallD

- firewalld는 내부적으로 iptables를 사용하여 규칙을 관리
- iptables 명령으로 추가한 규칙은 firewalld에서 생성한 규칙보다 우선 적용됨
- `AllowZoneDrifting`은 하나의 NI가 여러 Zone에 속한 경우 Zone 간 규칙을 드리프트
    - Zone 경계 간 규칙을 느슨하게 만들어 의도치 않은 트래픽 허용 가능성
    - 이후 버전에서는 제거될 예정
- L4 방화벽이기 때문에 NAT를 통해 ip가 변환되는 공인 ip들은 제어하기 어려움
    - firewalld의 한계

## 작업 과정

## 1. Firewall 포트포워딩

### 1) FW1 방화벽 설정

- 사용자 생성

```bash
useradd -m -G wheel fw1
echo  "fw1" | passwd --stdin fw1
```

- 패키지 설치

```bash
sudo dnf update --exclude=kernel* -y

echo "<FW1-IP> FW1" | sudo tee /var/www/html/index.html
sudo systemctl enable --now httpd
```

- 방화벽 설정

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --reload
```

### 2) FW2 방화벽 설정

- 사용자 생성

```bash
useradd -m -G wheel fw2
echo "fw2" | passwd --stdin fw2
```

- 패키지 설치

```bash
sudo dnf update --exclude=kernel* -y

echo "<FW2-IP> FW2 " | sudo tee /var/www/html/index.html
sudo systemctl enable --now httpd
```

- 방화벽 설정

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --reload
```

### 3) 포트포워딩

- 방화벽 설정

```bash
sudo firewall-cmd --permanent --add-forward-port=port=80:proto=tcp:toport=80:toaddr=<FW2-IP>
sudo firewall-cmd --permanent --add-masquerade
sudo firewall-cmd --reload
```

- 포트포워딩 테스트

```bash
sudo firewall-cmd --list-all
curl <공인 IP>:8080
curl <공인 IP>:8081
```

- 두 curl 모두 FW2 서버가 응답

### 4) 웹사이트 외부 접속 차단 - rich rule

### 1. FW2에서 목적지가 FW2 ip인 모든 트래픽 거부

- FW2로의 모든 트래픽 거부

```bash
# FW2 
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" destination address="<FW2-IP>" reject'
sudo firewall-cmd --reload

# 위에서 포트포워딩 했기 때문에 아래 명령어 둘 모두 차단됨
curl <Domain>:8080
curl <Domain>:8081
```

- FW2 ip로의 접속 차단됨

### 2. 출발지가 FW2 사설 ip인 모든 트래픽 거부

- FW1에서 FW2 사설 ip로 포트포워딩된 공인 ip로 요청한다면 차단 안됨
- 트러블슈팅 참고

```bash
# FW1 포트포워딩 제거
sudo firewall-cmd --permanent --remove-forward-port=port=80:proto=tcp:toport=80:toaddr=<FW2-IP>
sudo firewall-cmd --reload

# FW2 
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="<FW1-IP>" reject’
sudo firewall-cmd --reload

# FW1
curl <Domain>:8080
```

- FW1에서 FW2 사설 ip로의 접속 차단됨

### 3. FW2에서 목적지가 FW1인 트래픽 거부

- 트러블슈팅 참고
    - iptable 사용

```bash
# FW2 
sudo iptables -I OUTPUT -s <FW2-IP> -d <FW1-IP> -j REJECT

# FW2
# 포트포어딩후 NAT으로 ip가 달라지기 때문에 차단 불가
curl <Domain:8080
# 차단 성공
curl <FW1-IP>
```

- 트러블슈팅 참고
    - firewalld 사용

```bash
sudo vi /etc/firewalld/firewalld.conf
AllowZoneDrifting=no

sudo systemctl restart firewalld
sudo firewall-cmd --reload

sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="<FW2-IP>" destination address="<FW1-IP>" reject'
sudo firewall-cmd --reload

# 차단 성공
curl <FW1-IP>
```

---

## 로그 확인 명령어

```bash
# 패킷 캡쳐
sudo tcpdump -i eth0 host <FW1-IP> -n

# Connection Tracking 초기화
sudo conntrack -F

# firewall-cmd --reload 적용시 로그 발생 가능성
sudo tail -f /var/log/messages
```

## 트러블슈팅

> KT VM에 포트포워딩중인 공인 ip를 `sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" destination address="공인ip" reject'`로 차단 시 여전히 `curl <Domain>:8080` 이 동작
> 
- 공인 ip로 요청이 오면 사설 ip로 변환된 후 vm에 전달됨
- **firewalld**는 VM 내부 트래픽에만 영향을 미침
    - 따라서 공인 ip를 차단하지 못함

> FW2 에서 `sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="<FW1-IP>" reject’` 로 FW1 ip차단 시 FW1에서 `curl <FW2-IP>` 은 차단 되지만 `curl <Domain>:8081` / `curl <공인 IP>:8081` 은 차단되지 않음
> 
- 공인ip:8081은 포트포워딩 되어있기 때문에 kt cloud가 요청이 FW2 에 도달하기 전에 NAT 수행
- NAT는 출발지 ip를 공인ip가 아닌 kt cloud에서 자체적으로 사설 ip를 임의로 설정
    - firewalld 만으로는 FW2 에 매핑된 공인 ip를 통한 트래픽은 차단하기 어려움
        - `sudo tcpdump -i eth0 port 80 -n` 명령어로 NAT변환된 ip를 추가 차단하거나
        - zone을 drop으로 변경해서 화이트리스트 방식으로 변경하거나
        - kt가 nat하는 ip의 대역을 전체 차단하거나

> FW2 에서 `sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" destination address="<FW1-IP>" reject’` 설정 시 `curl <FW1-IP>` 차단 불가
> 
- 아웃바운드 트래픽 차단 불가
- iptable 명령어 `sudo iptables -I OUTPUT -d <FW1-IP> -j ACCEPT` 는 정상 작동
- /var/log/messages 로그 확인해보니 --reload 시에 `WARNING: AllowZoneDrifting is enabled. This is considered an insecure configuration option. It will be removed in a future release. Please consider disabling it now.` 로그 발생.
    - AllowZoneDrifting 설정 때문에 발생하는 문제로 판단
    - firewalld.conf 에 AllowZoneDrifting no로 설정 후 해결

## 참고

firewalld 사용법: 
https://www.freekb.net/Article?id=2141
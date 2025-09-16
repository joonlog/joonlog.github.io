---
title : HAproxy 동시 재기동 오류 TroubleShooting
date : 2025-09-16 09:00:00 +09:00
categories : [Linux, LoadBalancer]
tags : [linux, troubleshooting, haproxy, loadbalancer] #소문자만 가능
---

- 두 서버의 HAProxy를 동시에 재기동하는 과정에서 VIP FAULT 현상이 발생
    - keepalived/corosync/pacemaker 같은 클러스터 프로세스가 VIP 소유권을 제어하는 앱에서 자주 있는 일
    - VIP(Virtual IP) 소유권 충돌이 발생하여 이로 인해 VIP가 정상적으로 할당되지 않고 모두 FAULT 상태로 전환 이후 BACKUP 상태로 변경됨
    - 이후 순차적으로 keepalived를 재기동해서 정상적으로 Master 상태로 전환 확인
        - `systemctl restart keepalived`
- 위 현상 방지하기 위해 한 HAproZxy 재기동 후, 재기동하지 않은 서버가 MASTER 승격되었는지, VIP 활성화 되었는지 확인 필수
    
    ```bash
    # HAProxy 재기동 시 로그
    
    Reloading HAProxy Load Balancer...
    Stopping proxy haproxy-ssl...
    Stopping proxy haproxy...
    
    Keepalived_vrrp[####]: (VI_1) received an invalid passwd!
    Keepalived_vrrp[####]: (VI_1) Dropping received VRRP packet...
    Keepalived_vrrp[####]: VRRP_Instance(VI_1) Entering FAULT STATE
    Keepalived_vrrp[####]: VRRP_Script(chk_haproxy) failed
    Keepalived_vrrp[####]: VRRP_Instance(VI_1) Now in FAULT state
    Keepalived_vrrp[####]: VRRP_Instance(VI_1) Entering BACKUP STATE
    
    # haproxy_test.sh 실행 결과
    /keepalived/haproxy_test.sh exited with status 1
    /keepalived/haproxy_test.sh exited with status 0
    
    # BACKUP 전환 확인
    Keepalived_vrrp[####]: VRRP_Instance(VI_1) Entering BACKUP STATE
    ```
---
title: "Cacti Conainer 구축"
date: 2025-07-30T09:00:00+09:00
categories: ["Container", "Cacti"]
tags: ["docker", "container", "cacti", "cacti container"]
---


- Containerd 사용

### Git 소스 코드

```bash
dnf install -y git-all
git clone https://github.com/scline/docker-cacti.git
```

### Containerd 환경 구축

- https://containerd.io/downloads/ 여기서 다운로드도 가능
- https://medium.com/@DannielWhatever/using-containerd-without-docker-9d08332781b4
    
    ```bash
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo dnf remove $pkg -y; done
    dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    dnf install -y containerd.io
    containerd config default | sudo tee /etc/containerd/config.toml
    
    ### containerd root 디렉토리 외에 data 디렉토리에서 동작하도록 설정
    sed -i 's|^root = .*|root = "/data/container/containerd"|' /etc/containerd/config.toml
    sed -i 's|^state = .*|state = "/data/container/containerd-state"|' /etc/containerd/config.toml
    mkdir -p /data/container/containerd
    mkdir -p /data/container/containerd-state
    mkdir -p /data/container/nerdctl
    echo 'export NERDCTL_DATA_ROOT=/data/container/nerdctl' >> ~/.bashrc
    source ~/.bashrc
    
    # nerdctl 설치
    wget https://github.com/containerd/nerdctl/releases/download/v0.17.0/nerdctl-0.17.0-linux-amd64.tar.gz
    tar Cxzvvf /usr/bin nerdctl-0.17.0-linux-amd64.tar.gz
    echo "source <(nerdctl completion bash)" >> ~/.bashrc
    
    # CNI 설치
    wget https://github.com/containernetworking/plugins/releases/download/v1.0.1/cni-plugins-linux-amd64-v1.0.1.tgz
    mkdir -p /opt/cni/bin
    tar Cxzvvf /opt/cni/bin cni-plugins-linux-amd64-v1.0.1.tgz
    
    # buildkit 설치
    wget https://github.com/moby/buildkit/releases/download/v0.10.0-rc1/buildkit-v0.10.0-rc1.linux-amd64.tar.gz 
    tar Cxzvvf /usr/bin buildkit-v0.10.0-rc1.linux-amd64.tar.gz 
    nohup /usr/bin/builkitd < /dev/null > /var/log/buildkitd 2>&1 &
    
    systemctl enable --now containerd
    ```

### Cacti Container 구조
- cacti master와 master용 DB
- cacti poller와 poller용 DB
    
- docker-compose-master.yaml
    
    ```bash
    version: '3.5'
    
    services:
      cacti-master:
        image: "smcline06/cacti"
        container_name: cacti_master
        hostname: cactimaster
        depends_on:
          - db-master
        ports:
          - "80:80"
          - "443:443"
        environment:
          - DB_NAME=cacti_master
          - DB_USER=cactiuser
          - DB_PASS=cactipassword
          - DB_HOST=db-master
          - DB_PORT=3306
          - DB_ROOT_PASS=rootpassword
          - INITIALIZE_DB=1
          - TZ=Asia/Seoul
        volumes:
          - cacti-master-data:/cacti
          - cacti-shared-rra:/cacti/rra
          - cacti-master-spine:/spine
          - cacti-master-backups:/backups
        networks:
          - cacti-net
    
      db-master:
        image: "mariadb:10.3"
        container_name: cacti_master_db
        hostname: db-master
        ports:
          - "3306:3306"
        command:
          - mysqld
          - --character-set-server=utf8mb4
          - --collation-server=utf8mb4_unicode_ci
          - --max_connections=200
          - --max_heap_table_size=128M
          - --max_allowed_packet=32M
          - --tmp_table_size=128M
          - --join_buffer_size=128M
          - --innodb_buffer_pool_size=1G
          - --innodb_doublewrite=ON
          - --innodb_flush_log_at_timeout=3
          - --innodb_read_io_threads=32
          - --innodb_write_io_threads=16
          - --innodb_buffer_pool_instances=9
          - --innodb_file_format=Barracuda
          - --innodb_large_prefix=1
          - --innodb_io_capacity=5000
          - --innodb_io_capacity_max=10000
        environment:
          - MYSQL_ROOT_PASSWORD=rootpassword
          - TZ=Asia/Seoul
        volumes:
          - cacti-db-master:/var/lib/mysql
        networks:
          - cacti-net
    
    volumes:
      cacti-db-master:
      cacti-master-data:
      cacti-shared-rra:
      cacti-master-spine:
      cacti-master-backups:
    
    networks:
      cacti-net:
        name: cacti-net
        external: true
    ```
    
- docker-compose-poller.yaml
    
    ```bash
    version: '3.5'
    
    services:
      cacti-poller:
        image: "smcline06/cacti"
        container_name: cacti_poller
        hostname: cactipoller
        depends_on:
          - db-poller
        ports:
          - "8080:80"
          - "8443:443"
        environment:
          - DB_NAME=cacti_poller
          - DB_USER=cactipolleruser
          - DB_PASS=cactipollerpassword
          - DB_HOST=db-poller
          - DB_PORT=3306
          - RDB_NAME=cacti_master
          - RDB_USER=cactiuser
          - RDB_PASS=cactipassword
          - RDB_HOST=db-master
          - RDB_PORT=3306
          - DB_ROOT_PASS=rootpassword
          - REMOTE_POLLER=1
          - INITIALIZE_DB=1
          - TZ=Asia/Seoul
        volumes:
          - cacti-poller-data:/cacti
          - cacti-shared-rra:/cacti/rra
          - cacti-poller-spine:/spine
          - cacti-poller-backups:/backups
    
        networks:
          - cacti-net
    
      db-poller:
        image: "mariadb:10.3"
        container_name: cacti_poller_db
        hostname: db-poller
        command:
          - mysqld
          - --character-set-server=utf8mb4
          - --collation-server=utf8mb4_unicode_ci
          - --max_connections=200
          - --max_heap_table_size=128M
          - --max_allowed_packet=32M
          - --tmp_table_size=128M
          - --join_buffer_size=128M
          - --innodb_buffer_pool_size=1G
          - --innodb_doublewrite=ON
          - --innodb_flush_log_at_timeout=3
          - --innodb_read_io_threads=32
          - --innodb_write_io_threads=16
          - --innodb_buffer_pool_instances=9
          - --innodb_file_format=Barracuda
          - --innodb_large_prefix=1
          - --innodb_io_capacity=5000
          - --innodb_io_capacity_max=10000
        environment:
          - MYSQL_ROOT_PASSWORD=rootpassword
          - TZ=Asia/Seoul
        volumes:
          - cacti-db-poller:/var/lib/mysql
        networks:
          - cacti-net
    
    volumes:
      cacti-db-poller:
      cacti-poller-data:
      cacti-shared-rra:
        name: containerd-cacti_cacti-shared-rra
        external: true
      cacti-poller-spine:
      cacti-poller-backups:
    
    networks:
      cacti-net:
        name: cacti-net
        external: true
    ```
    
- cacti-net CNI
    - ip 지정해서 쓰기 위한 사전 세팅
    
    ```bash
    cat <<EOF > /etc/cni/net.d/cacti-net.conflist
    {
      "cniVersion": "0.4.0",
      "name": "cacti-net",
      "nerdctlID": 2,
      "nerdctlLabels": {
        "com.docker.compose.network": "cacti-net",
        "com.docker.compose.project": "cacti"
      },
      "plugins": [
        {
          "type": "bridge",
          "bridge": "cactibr0",
          "isGateway": true,
          "ipMasq": true,
          "ipam": {
            "type": "host-local",
            "ranges": [
              [
                {
                  "subnet": "10.99.0.0/24",
                  "gateway": "10.99.0.1"
                }
              ]
            ],
            "routes": [
              { "dst": "0.0.0.0/0" }
            ]
          }
        },
        {
          "type": "portmap",
          "capabilities": { "portMappings": true }
        },
        {
          "type": "firewall"
        },
        {
          "type": "tuning"
        }
      ]
    }
    
    EOF
    ```
    
    - cacti 컨테이너 실행
    
    ```bash
    nerdctl compose -f docker-compose-master.yaml up -d
    nerdctl compose -f docker-compose-poller.yaml up -d
    
    # netstat으로 port 확인 안됨
    iptables -t nat -L -n -v
    ```
    
    ### SNMP 체크
    
    ```bash
    snmpwalk -v2c -c public <ip>:<port> system
    ```
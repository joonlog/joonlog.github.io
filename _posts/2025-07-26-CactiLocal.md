---
title : Rocky8.10 Cacti 설치
date : 2025-07-26 09:00:00 +09:00
categories : [Linux, Monitoring]
tags : [linux, rocky 8.10, cacti]  #소문자만 가능
---

참고: https://ko.linux-console.net/?p=2502

1. Apache 웹서버
    
    ```bash
    sudo dnf install httpd -y
    sudo systemctl enable --now httpd
    ```
    
2. MariaDB
    
    ```bash
    sudo dnf install -y mariadb-server mariadb
    sudo systemctl enable --now mariadb
    ```
    
3. PHP
    
    ```bash
    sudo dnf install dnf-utils http://rpms.remirepo.net/enterprise/remi-release-8.rpmmi 
    sudo dnf module reset php
    sudo dnf module enable php:remi-7.4
    sudo dnf install @php
    sudo dnf install -y php php-{mysqlnd,curl,gd,intl,pear,recode,ldap,xmlrpc,snmp,mbstring,gettext,gmp,json,xml,common}
    sudo systemctl enable --now php-fpm
    
    sudo vim /etc/opt/remi/php74/php.ini
    
    date.timezone = Africa/Nairobi
    memory_limit = 512M
    max_execution_style = 60
    ```
    
4. SNMP & RRDtool
    
    ```bash
    sudo dnf install -y net-snmp net-snmp-utils net-snmp-libs rrdtool
    sudo systemctl enable --now snmpd
    ```
    
5. Cacti DB 생성
    - cacti 사용자 생성 및 cacti DB ALL 권한 부여
    - cacti 사용자에 timezone 조회 권한 부여
    
    ```sql
    mysql -u root -p
    
    MariaDB [(none)]> CREATE DATABASE cactidb;
    MariaDB [(none)]> GRANT ALL ON cactidb.* TO cacti_user@localhost IDENTIFIED  BY 'passwd123';
    MariaDB [(none)]> FLUSH PRIVILEGES;
    MariaDB [(none)]> EXIT;
    
    mysql -u root -p mysql < /usr/share/mariadb/mysql_test_data_timezone.sql
    
    MariaDB [(none)]> GRANT SELECT ON mysql.time_zone_name TO cacti_user@localhost;
    MariaDB [(none)]> FLUSH PRIVILEGES;
    MariaDB [(none)]> EXIT;
    ```
    
    - 선택 - DB 파라미터 튜닝
    
    ```bash
    sudo vi /etc/my.cnf.d/mariadb-server.cnf
    
    collation-server=utf8mb4_unicode_ci
    character-set-server=utf8mb4
    max_heap_table_size=32M
    tmp_table_size=32M
    join_buffer_size=64M
    25% Of Total System Memory
    innodb_buffer_pool_size=1GB
    pool_size/128 for less than 1GB of memory
    innodb_buffer_pool_instances=10
    innodb_flush_log_at_timeout=3
    innodb_read_io_threads=32
    innodb_write_io_threads=16
    innodb_io_capacity=5000
    innodb_file_format=Barracuda
    innodb_large_prefix=1
    innodb_io_capacity_max=10000
    ```
    
6. Cacti 설치 및 테이블 생성
    
    ```bash
    sudo dnf install epel-release -y
    sudo dnf install cacti -y
    rpm -qi cacti
      cacti 1.12.17
    rpm -ql cacti | grep cacti.sql
      /usr/share/doc/cacti/cacti.sql
    mysql -u root -p cactidb < /usr/share/doc/cacti/cacti.sql
    ```
    
7. Cacti MariaDB 연동
    
    ```bash
    sudo vim /usr/share/cacti/include/config.php
    
    $database_type     = 'mysql';
    $database_default  = 'cactidb';
    $database_hostname = 'localhost';
    $database_username = 'cacti_user';
    $database_password = 'Passwd123';
    $database_port     = '3306';
    $database_retries  = 5;
    $database_ssl      = false;
    $database_ssl_key  = '';
    ```
    
8. Cacti 풀링 설정
    
    ```bash
    sudo vim /etc/cron.d/cacti
    
    */5 * * * *   apache /usr/bin/php /usr/share/cacti/poller.php > /dev/null 2>&1
    ```
    
9. Cacti 웹페이지 설정 
    
    ```bash
    sudo vim /etc/httpd/conf.d/cacti.conf
    
    Alias /cacti   /usr/share/cacti
    
    <Directory /usr/share/cacti/>
        <IfModule mod_authz_core.c>
            # httpd 2.4
            # Require all granted
            Require ip 192.168.0.0/24
        </IfModule>
        <IfModule !mod_authz_core.c>
            # httpd 2.2
            Order deny,allow
            Deny from all
            # Allow from all
            Allow from 192.168.0.0/24
        </IfModule>
    </Directory>
    
    <Directory /usr/share/cacti/install>
    </Directory>
    
    sudo systemctl restart httpd
    sudo systemctl restart php-fpm
    sudo firewall-cmd --permanent --add-service=http
    sudo firewall-cmd --reload
    ```
    
10. Cacti 웹페이지 접근 
    
    ```bash
    http://server-ip/cacti
     admin/admin
    ```
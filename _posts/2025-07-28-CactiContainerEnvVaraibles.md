---
title : Cacti Conainer 환경변수
date : 2025-07-28 09:00:00 +09:00
categories : [Docker, Cacti]
tags : [docker, cacti, cacti container, cacti env]  #소문자만 가능
---

### Cacti Master

| 환경 변수 | 기능 |
| --- | --- |
| DB_NAME | MySQL DB 이름. Cacti, Spine Poller 공유 |
| DB_USER | MySQL DB 사용자. Cacti, Spine Poller 공유 |
| DB_PASS | MySQL DB 비밀번호. Cacti, Spine Poller 공유 |
| DB_HOST | DB IP/FQDN/호스트/컨테이너이름 |
| DB_PORT | MySQL DB 포트. 기본값 3306 |
| DB_ROOT_PASS | `INITIALIZE_DB`가 1로 설정된 경우 필요 |
| INITIALIZE_DB | `0: false` / `true: 1`. true면 컨테이너는 대상 데이터베이스에 대한 `DB_ROOT_PASS`를 요구. DB 초기화 설정 |
| TZ | 타임존 설정 |
| BACKUP_RETENTION | 보관할 백업 파일의 수 |
| REMOTE_POLLER | `0: false` / `true: 1` |
| PHP_MEMORY_LIMIT | PHP 메모리 제한 조정, 기본값은 128M |
| PHP_MAX_EXECUTION_TIME | PHP 최대 실행 시간 조정, 기본값은 30초 |
| PHP_SNMP | 0으로 설정하면 PHP-SNMP 제거. 기본값은 php-snmp 활성화. 일부 스크립트나 snmpv3가 제대로 동작하기 위해 필요 |

### **Cacti Pollers**

| 환경 변수 | 기능 |
| --- | --- |
| DB_NAME | MySQL DB 이름. Cacti, Spine Poller 공유 |
| DB_USER | MySQL DB 사용자. Cacti, Spine Poller 공유 |
| DB_PASS | MySQL DB 비밀번호. Cacti, Spine Poller 공유 |
| DB_HOST | MySQL DB IP/FQDN/호스트/컨테이너이름 |
| DB_PORT | MySQL DB 포트. 기본값 3306 |
| CACTI_URL_PATH | 기본값은 cacti. ex) http://<ip>/cacti. 루트 URL 변경만 허용되므로 “cacti” 만 변경 가능. “ccc/bb” 같은 것 불가능 |
| DB_ROOT_PASS | `INITIALIZE_DB`가 1로 설정된 경우 필요 |
| INITIALIZE_DB | `0: false` / `true: 1`. true면 컨테이너는 대상 데이터베이스에 대한 `DB_ROOT_PASS`를 요구. DB 초기화 설정 |
| TZ | 타임존 설정 |
| BACKUP_RETENTION | 보관할 백업 파일의 수 |
| REMOTE_POLLER | `0: false` / `true: 1` |
| RDB_NAME | Cacti Master MySQL DB 이름. Cacti, Spine Poller 공유 |
| RDB_USER | Cacti Master MySQL DB 사용자 |
| RDB_PASS | Cacti Master MySQL DB 사용자의 비밀번호 |
| RDB_HOST | Cacti Master MySQL DB가 사용하는 IP/FQDN/호스트/컨테이너이름 |
| RDB_PORT | Cacti Master MySQL DB 사용 포트 |

### DB 파라미터

| **MySQL 변수** | **추천 값** | **설명** |
| --- | --- | --- |
| **Version** | >= 5.6 | MySQL 5.6+ 및 MariaDB 10.0+은 훌륭한 릴리스 버전으로, 안정성과 성능이 뛰어납니다. 특히 최신 릴리스는 네트워킹 문제를 해결하여 Spine의 신뢰성을 개선합니다. |
| **collation_server** | utf8mb4_unicode_ci | 영어 이외의 언어를 사용할 경우, 일부 문자가 1바이트 이상을 차지하므로 `utf8mb4_unicode_ci` 정렬 타입을 사용하는 것이 중요합니다. |
| **character_set_client** | utf8mb4 | 영어 이외의 언어를 사용할 경우, 일부 문자가 1바이트 이상을 차지하므로 `utf8mb4` 문자셋을 사용하는 것이 중요합니다. |
| **max_connections** | >= 100 | Spine 데이터 수집기와 로그인 수에 따라 많은 MySQL 연결이 필요할 수 있습니다. Spine의 연결 수 계산법은 다음과 같습니다: `total_connections = total_processes * (total_threads + script_servers + 1)`. 사용자 연결을 위한 여유 공간을 남겨야 합니다. |
| **max_heap_table_size** | >= 10% RAM | Cacti 성능 향상 기능을 사용하여 메모리 스토리지 엔진을 선택할 경우, 시스템 메모리 테이블 공간이 부족하기 전에 Performance Booster 버퍼를 플러시해야 합니다. 이 값은 시스템 메모리의 10%로 설정하는 것이 권장되지만, SSD 디스크를 사용하거나 작은 시스템에서는 무시할 수도 있습니다. |
| **max_allowed_packet** | >= 16777216 | 원격 폴링 기능을 사용할 때, 많은 데이터가 메인 서버에서 원격 폴러로 동기화됩니다. 따라서 이 값을 16M 이상으로 설정해야 합니다. |
| **tmp_table_size** | >= 64M | 서브쿼리를 실행할 때, 임시 테이블 크기를 크게 설정하면 메모리 내에 임시 테이블을 유지할 수 있습니다. |
| **join_buffer_size** | >= 64M | 조인 연산이 이 크기보다 작으면 메모리에 저장되어 임시 파일로 쓰여지지 않습니다. |
| **innodb_file_per_table** | ON | InnoDB 스토리지를 사용할 때, 테이블 공간을 분리하여 관리가 더 간편해집니다. 현재 OFF로 설정되어 있다면 이 기능을 활성화하고 InnoDB 테이블에 대해 `ALTER` 명령을 실행하여 마이그레이션할 수 있습니다. |
| **innodb_file_format** | Barracuda | `innodb_file_per_table`을 사용할 때, 파일 형식을 Barracuda로 설정해야 합니다. 이는 특정 Cacti 테이블에 필요한 긴 인덱스를 허용합니다. |
| **innodb_large_prefix** | 1 | 테이블에 매우 큰 인덱스가 있는 경우, Barracuda 파일 형식과 함께 이 값을 1로 설정해야 합니다. 그렇지 않으면 일부 플러그인이 테이블을 제대로 생성하지 못할 수 있습니다. |
| **innodb_buffer_pool_size** | >= 25% RAM | InnoDB는 가능한 한 많은 테이블과 인덱스를 시스템 메모리에 보관합니다. 이 값을 충분히 크게 설정하여 테이블과 인덱스를 메모리에 저장하세요. `/var/lib/mysql/cacti` 디렉토리의 크기를 확인하면 적절한 값을 설정하는 데 도움이 됩니다. |
| **innodb_doublewrite** | ON | ZFS 또는 FusionI/O와 같은 내부 저널링 기능을 가진 파일 시스템을 사용하는 경우가 아니라면 이 값을 ON으로 유지하세요. 하지만 시스템 안정성이 매우 높고 백업이 잘 되어 있다면 OFF로 설정하여 데이터베이스 성능을 약 50% 개선할 수 있습니다. |
| **innodb_lock_wait_timeout** | >= 50 | 비정상적인 쿼리로 인해 데이터베이스가 다른 사용자에게 응답하지 못하는 것을 방지합니다. |
| **innodb_flush_log_at_timeout** | >= 3 | 높은 I/O 시스템에서 이 값을 1초 이상으로 설정하면 디스크 I/O가 더 순차적으로 작동할 수 있습니다. |
| **innodb_read_io_threads** | >= 32 | 최신 SSD 스토리지에서는 높은 I/O 특성을 가진 애플리케이션에 다중 읽기 I/O 스레드가 유리합니다. |
| **innodb_write_io_threads** | >= 16 | 최신 SSD 스토리지에서는 높은 I/O 특성을 가진 애플리케이션에 다중 쓰기 I/O 스레드가 유리합니다. |
| **innodb_buffer_pool_instances** | >= 9 | InnoDB는 `innodb_buffer_pool`을 메모리 영역으로 분리하여 성능을 향상합니다. 1GB 이하에서는 `pool size / 128MB`를 사용하고, 최대값은 64로 설정합니다. |
| **innodb_io_capacity** | 5000 | SSD 디스크에서는 5000을 사용하고, 물리적 디스크에서는 활성 드라이브 수에 200을 곱한 값을 사용합니다. NVMe 또는 PCIe 플래시의 경우 최대 100,000까지 설정 가능합니다. |
| **innodb_io_capacity_max** | 10000 | SSD 디스크에서는 10,000을 사용하고, 물리적 디스크에서는 활성 드라이브 수에 2000을 곱한 값을 사용합니다. NVMe 또는 PCIe 플래시의 경우 최대 200,000까지 설정 가능합니다. |
| **memory_limit** | >= 800M | 메모리 제한은 최소 800MB로 설정하세요. |
| **max_execution_time** | >= 60 | PHP 스크립트 실행 시간을 최소 60초로 설정하세요. |
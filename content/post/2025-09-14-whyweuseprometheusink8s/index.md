---
title: "Kubernetes 환경에서 Prometheus를 많이 사용하는 이유"
date: 2025-09-12T09:00:00+09:00
categories: ["Kubernetes", "Monitoring"]
tags: ["kubernetes", "k8s", "prometheus", "metric monitoring", "kubelet", "cadvisor"]
---


> 흔히 Kubernetes 환경에서 매트릭 모니터링으로 Prometheus를 많이 사용한다는데, 왜 많이 쓰는지, 정확한 이유를 알아보고자 작성한 글입니다.
> 

### K8S에서 Prometheus를 쓰는 이유

Prometheus 구축을 kube-prometheus-stack Helm Chart로만 구축해봤다 보니 처음에는 Helm에 그 이유가 있을 줄 알았는데, 알고 보니 Prometheus 자체가 K8S와의 호환성이 좋았다. Helm은 Prometheus를 구축하기 편하게 Chart 형태로만 패키징 해놓은 것. 이 Helm Chart의 베이스 엔진인 Prometheus-Operator도 이유 중 하나로, CRD(podmonitor, servicemonitor)로 promtheus를 관리할 수 있다.

결론적으로, 크게 아래 두 이유 때문에 K8S에서 Prometheus를 사용한다고 볼 수 있다.

1. K8S 환경에서의 호환성
2. K8S 환경에서의 관리 용이성

### K8S 환경에서의 Prometheus 호환성

> Prometheus는 K8S 서비스 디스커버리를 지원
> 
- 내장된 서비스 디스커버리로 인해 K8S API(서비스/파드/노드/인그레스 등)가 스케일인/아웃, 신규 네임스페이스 추가되는 등등의 상황에서도 동적으로 탐색 가능
    - 서비스 디스커버리: 엔드포인트를 미리 하드코딩하지 않고, 런타임에 자동으로 찾아 쓰는 패턴
- K8S에서는 kubelet 프로세스가 cAdvisor를 통해서 지표를 수집하는데, 이 데이터를 Prometheus 서버가 읽어서 신규 서비스/파드 추가 시 모니터링 대상으로 자동으로 편입 가능
    - prometheus.yml에서 kubelet의 엔트포인트인 /metrics, /metrics/cadvisor를 스크랩하도록 설정할 경우에 Prometheus가 자동으로 읽음
    - 바닐라 prometheus를 사용할 경우의 prometheus.yml 예시
    
    ```bash
    scrape_configs:
      - job_name: kubelet
        scheme: https
        metrics_path: /metrics
        kubernetes_sd_configs:
          - role: node
        bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
        tls_config:
          ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        relabel_configs:
          - source_labels: [__meta_kubernetes_node_address_InternalIP]
            target_label: __address__
            replacement: $1:10250
    
      - job_name: kubelet-cadvisor
        scheme: https
        metrics_path: /metrics/cadvisor
        kubernetes_sd_configs:
          - role: node
        bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
        tls_config:
          ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        relabel_configs:
          - source_labels: [__meta_kubernetes_node_address_InternalIP]
            target_label: __address__
            replacement: $1:10250
    ```
    

### K8S 환경에서의 Prometheus 관리 용이성

- Prometheus Operator:
https://github.com/prometheus-operator/prometheus-operator
- Prometheus Operator는 Prometheus를 K8S에서 사용하기 편하게 패키징한 툴
    - CRD 기반의 선언형으로 Prometheus 관리 가능
        - CRD 예: ServiceMonitor, PodMonitor, PrometheusRule 등의 커스텀 리소스
        - CRD에서의 Label Selector, Namespace Selector 설정으로 신규 파드 추가 시 자동 모니터링 편입
        - kubelet을 ServiceMonitor로 생성할 경우 kubelet의 cadvisor가 수집한 지표를 이 kubelet CRD로 인해 Prometheus 서버가 읽고, 신규 서비스/파드 추가 시 모니터링 대상으로 자동으로 편입 가능
        - Prometheus CRD로 Alerting 룰 분리
- kube-prometheus-stack
    - prometheus operator부터 grafana까지 모니터링에 필요한 올인원 Helm 패키지
    - 구조
        
        ```bash
        kube-prometheus-stack
        ├── Prometheus Operator (핵심 엔진)
        ├── Prometheus 인스턴스 
        ├── Grafana (시각화 대시보드)
        ├── AlertManager (알람 관리)
        ├── Node Exporter (Node 매트릭 수집)
        ├── kube-state-metrics (K8S 리소스 매트릭 수집)
        └── 기본 ServiceMonitors/PrometheusRules (정책)
        ```
        
    - 기본 생성 servicemonitor 목록
        - 아래 리스트가 전부 생성되므로 사실상 kube-prometheus-stack을 설치하는 것만으로 기본 매트릭 모니터링은 끝
        
        ```bash
        prometheus-grafana
        prometheus-kube-prometheus-alertmanager
        prometheus-kube-prometheus-apiserver
        prometheus-kube-prometheus-coredns
        prometheus-kube-prometheus-kube-controller-manager
        prometheus-kube-prometheus-kube-etcd
        prometheus-kube-prometheus-kube-proxy
        prometheus-kube-prometheus-kube-scheduler
        prometheus-kube-prometheus-kubelet
        prometheus-kube-prometheus-operator
        prometheus-kube-prometheus-prometheus
        prometheus-kube-state-metrics
        prometheus-prometheus-node-exporter
        ```
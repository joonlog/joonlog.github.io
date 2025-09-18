---
title : 자체 관리형 Kubernetes에서의 Prometheus와 Grafana 구축
date : 2025-09-12 09:00:00 +09:00
categories : [Kubernetes, Monitoring]
tags : [kubernetes, k8s, self managed k8s, prometheus, grafana, metric monitoring]  #소문자만 가능
---

> 매니지드 K8S가 아닌 자체 관리형 K8S에서의 Prometheus + Grafana 구축 방법
> 

> UI 접근을 위한 외부 통신으로 MetalLB+HAproxy+Nginx Controller가 구성 되었다는 전제 하에 구축
> 
> 
> > CSP 환경이라면 Ingress 설정에서 AWS ALB Controller와 같은 로드밸런서를 사용해 훨씬 간편하게 구축 가능
> > 
- Prometheus 공식문서:
https://prometheus.io/docs/prometheus/latest/installation/
- Grafana 공식문서:
https://grafana.com/docs/grafana/latest/setup-grafana/installation/

### 개요

- kube-prometheus-stack helm chart 구조
    
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
    

Prometheus Operator란

- CRD 기반 Prometheus 관리 도구
- Prometheus, ServiceMonitor, PrometheusRule 등의 Custom Resource 제공
- 선언적 방식으로 Prometheus 설정 관리

### Prometheus + Grafana 설치

- 깃허브:
https://github.com/prometheus-community/helm-charts
1. helm 리포지토리 추가
    
    ```bash
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo update
    ```
    
2. kube-prometheus-stack 설치
    - prometheus + grafana 합쳐진 패키지
    - values.yaml 수정
    
    ```bash
    helm show values prometheus-community/kube-prometheus-stack > values.yaml
    ```
    
    - grafana pvc 설정
    
    ```bash
      persistence:
        enabled: true
        type: sts
        storageClassName: "longhorn"
        accessModes:
          - ReadWriteOnce
        size: 5Gi
    ```
    
    - prometheus pvc 설정
        - selector: {} 삭제
            - longhorn은 selector 지원하지 않음
    
    ```bash
      storageSpec:
        volumeClaimTemplate:
          spec:
            storageClassName: longhorn
            accessModes: ["ReadWriteOnce"]
            resources:
              requests:
                storage: 40Gi
    ```
    
    - alertmanager pvc 설정
        - selector: {} 삭제
            - longhorn은 selector 지원하지 않음
    
    ```bash
      storage:
        volumeClaimTemplate:
          spec:
            storageClassName: longhorn
            accessModes: ["ReadWriteOnce"]
            resources:
              requests:
                storage: 2Gi
    ```
    
    - kube-prometheus-stack 설치
    
    ```bash
    helm install prometheus prometheus-community/kube-prometheus-stack -f values.yaml --namespace monitoring --create-namespace
    ```
    

### Grafana UI 접근

1. ingress 생성
    - 8082포트를 사용하는 grafana service로 연결
    
    ```bash
    cat <<EOF > grafana-ingress.yaml
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: grafana-ingress
      namespace: monitoring
      annotations:
        nginx.ingress.kubernetes.io/backend-protocol: "HTTP"
        nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    spec:
      ingressClassName: nginx
      rules:
      - host: "grafana.local"
        http:
          paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: prometheus-grafana
                port:
                  number: 80
    EOF
    
    kubectl -n monitoring apply -f grafana-ingress.yaml
    ```
    
2. HAproxy 설정
- ingress를 조회하여 MetalLB로부터 부여받은 IP 확인
    - 현재는 172.27.1.100
    
    ```bash
    # kubectl get ingress -n monitoring 
    NAME              CLASS   HOSTS           ADDRESS        PORTS   AGE
    grafana-ingress   nginx   grafana.local   172.27.1.100   80      9s
    ```
    
- 확인한 IP를 HAproxy 설정에 추가
    - HAproxy 서버 공인 IP로 접근 시 Grafana의 Ingress로 통신되도록 설정
    - `http-request set-header Host`
        - 클라이언트에서 오는 모든 HTTP 요청의 Host 헤더를 grafana.local로 변경
    - `http-request del-header X-Forwarded-Host`
        - 이전 프록시에서 설정한 Host 헤더 제거
    - `http-request del-header X-Forwarded-Proto`
        - 이전 프록시에서 설정한 프로토콜 정보 제거
    
    ```bash
    tee -a /etc/haproxy/haproxy.cfg > /dev/null <<EOF
    frontend metallb_frontend_grafana
        bind *:8082
        mode http
        option forwardfor
        http-request set-header Host grafana.local
        http-request set-header X-Forwarded-Host %[req.hdr(host)]
        http-request set-header X-Forwarded-Proto http
        http-request set-header X-Forwarded-Port %[dst_port]
        default_backend metallb_backend_grafana

    backend metallb_backend_grafana
        server grafana 172.27.1.100:80
    EOF
    ```
    
3. 접근 성공!
    
    ![Prometheus1.png](/assets/img/kubernetes/Prometheus1.png)
    
    - Grafana 관리자 비밀번호 확인
        
        ```bash
        kubectl --namespace monitoring get secrets prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 -d ; echo
        ```
        
    
    ### 대시보드 그래프 origin not allowed 에러 발생 시
    
    - IP주소로의 직접 접근을 Grafana의 CSRF 정책이 요청 차단
        - 도메인을 통해 접근해서 CSRF 정책 통과
            - ingrass, haproxy에 설정되어있던 `grafana.local`을 새로 설정한 `도메인`으로 변경
            - values.yaml 수정
            
            ```bash
              grafana:
                ...
                adminPassword: prom-operator
                grafana.ini:
                  server:
                    domain: <도메인>:48082
                    root_url: http://<도메인>:48082
                  security:
                    csrf_trusted_origins: "<도메인>:48082"
                    csrf_cookies_secure: false
                    cookie_secure: false
                    cookie_samesite: lax
                  ...
            ```
            
            - 설정 반영
                
                ```bash
                
                helm uninstall prometheus --namespace monitoring
                kubectl --namespace monitoring delete pvc --all
                helm install prometheus prometheus-community/kube-prometheus-stack -f values.yaml --namespace monitoring
                ```
                

### Prometheus UI 접근

1. ingress 생성
    - 8080포트를 사용하는 prometheus service로 연결
    
    ```bash
    cat <<EOF > prometheus-ingress.yaml
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: prometheus-ingress
      namespace: monitoring
      annotations:
        nginx.ingress.kubernetes.io/backend-protocol: "HTTP"
        nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    spec:
      ingressClassName: nginx
      rules:
      - host: "prometheus.local"
        http:
          paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: prometheus-kube-prometheus-prometheus
                port:
                  number: 9090
    EOF
    
    kubectl -n monitoring apply -f prometheus-ingress.yaml
    ```
    
2. HAproxy 설정
- ingress를 조회하여 MetalLB로부터 부여받은 IP 확인
    - 현재는 172.27.1.100
    
    ```bash
    # kubectl get ingress -n monitoring 
    NAME                 CLASS   HOSTS              ADDRESS        PORTS   AGE
    prometheus-ingress   nginx   prometheus.local   172.27.1.100   80      9s
    ```
    
- 확인한 IP를 HAproxy 설정에 추가
    - HAproxy 서버 공인 IP로 접근 시 Prometheus의 Ingress로 통신되도록 설정
    - `http-request set-header Host`
        - 클라이언트에서 오는 모든 HTTP 요청의 Host 헤더를 prometheus.local로 변경
    - `http-request del-header X-Forwarded-Host`
        - 이전 프록시에서 설정한 Host 헤더 제거
    - `http-request del-header X-Forwarded-Proto`
        - 이전 프록시에서 설정한 프로토콜 정보 제거
    
    ```bash
    tee -a /etc/haproxy/haproxy.cfg > /dev/null <<EOF
    frontend metallb_frontend_prometheus
        bind *:8083
        mode http
        option forwardfor
        http-request set-header Host prometheus.local
        http-request set-header X-Forwarded-Host %[req.hdr(host)]
        http-request set-header X-Forwarded-Proto http
        http-request set-header X-Forwarded-Port %[dst_port]
        default_backend metallb_backend_prometheus

    backend metallb_backend_prometheus
        server prometheus 172.27.1.100:80
    EOF
    ```
    
3. 접근 성공!
    
    ![Grafana1.png](/assets/img/kubernetes/Grafana1.png)

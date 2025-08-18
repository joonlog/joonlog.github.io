---
title : Kubernetes 아키텍쳐
date : 2025-08-04 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, kubernetes architecture]  #소문자만 가능
---

## Master Node

### kube-apiserver

- Kubernetes 클러스터의 프론트엔드입니다.
- REST API를 통해 클러스터 상태를 읽고 조작하는 모든 요청을 처리합니다.

### kubectl

- `kubectl`은 **kube-apiserver**에 요청을 보냅니다. (예: `kubectl get nodes`)
- 또는, `kubectl`을 사용하지 않고도 직접 kube-apiserver에 요청을 보낼 수 있습니다.

### **controller-manager**

- 클러스터의 상태를 지속적으로 모니터링하고, 원하는 상태를 유지합니다.
- 다양한 컨트롤러들(잡, 노드, 서비스 계정, 엔드포인트)을 관리하여 리소스의 상태를 조정합니다.

### **kube-scheduler**

- **스케줄러**는 새로 생성된 Pod를 워커 노드에 배치하는 역할을 합니다.
- Pod의 리소스 요구사항과 각 노드의 리소스 가용성을 고려해 최적의 노드에 배치합니다.

### **ETCD**

- Kubernetes의 클러스터 데이터를 저장하는 분산 key-value 데이터베이스입니다.
- 클러스터의 현재 상태, 구성 정보, 시크릿, ConfigMap 등을 저장하는 데 사용됩니다.

---

## Worker Node

### **kubelet**

- 각 워커 노드에서 실행되며, API 서버와 통신하면서 Pod의 라이프사이클을 관리합니다.
- kube-scheduler가 할당한 작업을 워커 노드에서 실행합니다.

### **kube-proxy**

- 네트워크 트래픽을 라우팅하고 로드밸런싱하는 데 사용되는 기본 네트워크 프록시입니다.
- Pod과 서비스 간의 통신 규칙을 설정합니다.

### container runtime

- 실제 컨테이너를 실행하는 소프트웨어입니다.
- Containerd, CRI-O 등이 Kubernetes에서 사용될 수 있으며, 컨테이너 이미지를 다운로드하고 실행, 관리하는 역할을 합니다.

---

## Addons

### **CoreDNS**

- 클러스터 내에서 서비스 디스커버리와 DNS 네임 리졸버를 제공합니다
- 이를 통해 서비스 이름을 클러스터 내부의 IP로 변환할 수 있습니다.

### **CNI**

- 네트워크를 관리하기 위한 표준 인터페이스입니다.
- 각 Pod에 IP 주소를 할당해서, 다른 Pod와의 통신을 가능하게 만듭니다.
- 대표적인 플러그인으로는 Calico, Flannel 등이 있습니다.
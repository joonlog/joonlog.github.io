---
title : Kubernetes Service Overview
date : 2025-08-08 09:00:00 +09:00
categories : [Kubernetes, Architecture]
tags : [kubernetes, k8s, service, clusterip, nodeport, loadbalancer]  #소문자만 가능
---

- Pod 그룹에 대해 단일 진입점(VIP) 제공
- 외부에서 클러스터 내에 파드 그룹에 접근할 수 있도록 기능 제공

## 서비스 종류

- ClusterIP ⇒ (외부 → 내부) VIP
- NodePort ⇒ (외부 → 내부) VIP + NodePort
- LoadBalancer ⇒ (외부 → 내부) VIP + NodePort + LB’s IP
- ExternalName ⇒ (내부 → 외부) 도메인 매핑

## ClusterIP

- Pod그룹에 대해 단일 진입점(VIP)을 제공
- VIP는 사설 IP로 제공됨
- ~/root/kubespray/inventory/mycluster/group_vars/k8s_cluster
- kube_service_addresses: 10.233.0.0/18
- kube_pods_subnet: 10.233.64.0/

## NodePort

- Pod그룹에 대해 단일 진입점(VIP)을 제공
- 노드의 IP:Port(30000 ~ 32727)로 요청하면 Pod 그룹(같은 Label)에 대한 부하 분산을 제공

## LoadBalancer

- Pod그룹에 대해 단일 진입점(VIP)을 제공
- 노드의 IP:Port(30000 ~ 32727)로 요청하면 Pod 그룹(같은 Label)에 대한 부하 분산을 제공
- LB(VM(Haproxy), Container(Metalib))로 요청하면 Pod 그룹(같은 Label)에 대한 부하 분산을 제공

## ExternalName - DNS CNAME 역할

- (내부 → 외부/내부) DNS CNAME 등록 작업(서비스 이름 ← 매핑 → DNS 이름)
(pod 예) <podname>.default.pod.cluster.local ← 매핑 → google.com
(service 예) <svcname>.default.svc.cluster.local ← 매핑 → google.com

## 헤드리스 서비스

- 단일 진입점이 필요 없는 경우에 사용
- selector가 필요 없는 경우 ex) DB

Headless → DNS A record

ExternalName → DNS CNAME record

## 네트워크

- kube-proxy
    - worker node안에 pod에 연결할 때 사용하는 네트워크를 책임진다.
    - masquarading
    - loadbalancing(포트포워딩)
- kube-proxy 모드
    - iptables 모드 ⇒ iptables CMD
    - ipvs 모드 ⇒ ipvsadm CMD
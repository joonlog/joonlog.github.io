---
title: "Kubernetes Helm 설치"
date: 2025-08-03T09:00:00+09:00
categories: ["Kubernetes", "Plugins"]
tags: ["kubernetes", "k8s", "helm"]
---


- Kubernetes에서 사용하는 여러 플러그인 및 툴의 yaml 파일들을 패키징하여 쉽게 설치할 수 있는 툴
    - values.yaml을 통해 툴들의 옵션 조정 가능

```bash
cd
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh --version v3.17.4
```

- 3.18.5 버전에서 helm 설치 시에 values.schema.json이 원격 $ref(HTTPS)로 참조하는 JSON Schema를 불러오면 검증 단계에서 실패하는 이슈가 있어서 3.17.4 버전 명시
    - https://github.com/helm/helm/issues/31136?utm_source=chatgpt.com
    
    ```bash
    [root@master1 ~]# helm install my-release oci://ghcr.io/nginx/charts/nginx-ingress --version 2.2.1
    Pulled: ghcr.io/nginx/charts/nginx-ingress:2.2.1
    Digest: sha256:00f7a7017799eafc7f109cca10204a9dc2bc2d2f0feafb8eff328d13e2cedd12
    Error: INSTALLATION FAILED: values don't meet the specifications of the schema(s) in the following chart(s):
    nginx-ingress:
    failing loading "https://raw.githubusercontent.com/nginxinc/kubernetes-json-schema/master/v1.33.1/_definitions.json": invalid file url: https://raw.githubusercontent.com/nginxinc/kubernetes-json-schema/master/v1.33.1/_definitions.json
    ```
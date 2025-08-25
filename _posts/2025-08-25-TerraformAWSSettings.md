---
title : Terraform으로 AWS 리소스를 컨트롤하기 위한 세팅
date : 2025-08-25 09:00:00 +09:00
categories : [Terraform, AWS]
tags : [terraform, aws, awscli]  #소문자만 가능
---

- Terrafrom으로 AWS를 컨트롤하기 위한 사전 작업

### Terraform 설치

```bash
sudo yum install -y yum-utils shadow-utils
sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo
sudo yum -y install terraform
```

### awscli 설치

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

### Terraform, AWS 환경 설정

- terraform 자동완성, alias 설정
- aws 자동완성 설정
- 사전 발급한 aws access key 설정

```bash
cat << EOF >> ~/.bashrc 
complete -C '/usr/local/bin/aws_completer' aws
complete -C /usr/bin/terraform terraform
complete -C /usr/bin/terraform tf

alias tf='terraform'
alias tfap='tf apply --auto-approve'

export AWS_<ACCESS_KEY_ID>=<ACCESS_KEY_ID>
export AWS_<SECRET_ACCESS_KEY>=<SECRET_ACCESS_KEY>
export AWS_DEFAULT_REGION=ap-northeast-2
EOF
. ~/.bashrc
```
FROM golang:latest

# AWS SDK
RUN go get github.com/aws/aws-sdk-go
RUN go get github.com/aws/aws-lambda-go/lambda

# Go dependency management
ADD https://github.com/golang/dep/releases/download/v0.5.3/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Go Coverage
RUN go get github.com/axw/gocov/gocov
RUN go get github.com/axw/gocov/...
RUN go get -u gopkg.in/matm/v1/gocov-html

# Various packages
RUN apt-get update && apt-get install -y git ssh tar gzip zip unzip ca-certificates apt-utils golang-glide gnupg jq 

# Node and Yarn web package manager
RUN curl -sL https://deb.nodesource.com/setup_11.x  | bash -
RUN apt-get -y install nodejs
RUN npm install -g yarn

# Terraform
RUN curl -L -s https://releases.hashicorp.com/terraform/0.12.5/terraform_0.12.5_linux_amd64.zip -o /go/bin/terraform.zip
RUN unzip /go/bin/terraform.zip -d /go/bin
RUN chmod +x /go/bin/terraform

# Gomplate
RUN curl -L -s https://github.com/hairyhenderson/gomplate/releases/download/v3.5.0/gomplate_linux-amd64 -o /go/bin/gomplate
RUN chmod +x /go/bin/gomplate

# Python
RUN apt-get update && apt-get install -y python3-pip python3-dev && cd /usr/local/bin && ln -s /usr/bin/python3 python && pip3 install --upgrade pip

# AWS CLI
RUN pip3 install awscli --upgrade --user
ENV PATH="~/.local/bin:${PATH}"

# Resource Core Files
ADD client /root
ADD test /root
RUN chmod +x /root/*.sh


ENTRYPOINT [ "/bin/bash" ]
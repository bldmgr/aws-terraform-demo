&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;

# Ec2 Worker

This terraform project creates an Ec2 worker to run bash commands from a remote docker container. 

# Quick start

* Create a new terraform state S3 bucket using the [instructions below](#create-terraform-state-s3-bucket)
* Edit the ./test/terraform.tf to update the bucket = "<S3_BUCKET_NAME>" value 
* Optional - Create a new AMI image using the [instructions below](#create-custom-ami-using-packer)
* Test using the instruction below

```shell script
docker run -it $(docker build -q .)
cd ~
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=us-east-1
terraform init
export aws_ami_id=ami-0ec6d71f90a4daad0
export git_access_token=
terraform apply -auto-approve -var "git_access_token=$git_access_token" -var "ami_id=$aws_ami_id"
export LOCAL_IP=$(curl http://ipv4.icanhazip.com)
export INSTANCE_IP=$(terraform output -json | jq -r '.instance_ip.value' )
export SECURITY_GROUP=$(terraform output -json | jq -r '.security_group.value'  )
export INSTANCE_ID=$(terraform output -json | jq -r '.instance_id.value'  )
aws ec2 authorize-security-group-ingress --group-name $SECURITY_GROUP --protocol tcp --port 8081 --cidr $LOCAL_IP/32
go run client.go $INSTANCE_IP:8081 'ls -la'
```



## Create terraform state s3 bucket

```shell script
docker run -it $(docker build -q .)
cd ~
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=us-east-1
aws s3api create-bucket --bucket <S3_BUCKET_NAME> --region us-east-1
exit
```

## Create custom AMI using Packer

To create a new AMI out of the latest Amazon Linux AMI, and install the server.go application within the new created AMI.
* Edit the terraform.tf file to update the bucket = "<S3_BUCKET_NAME>" 
* Capture the AMI created 

```shell script
docker run -it $(docker build -q .)
cd ~
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_SOURCE_AMI=ami-0565af6e282977273
export AWS_REGION=us-east-1
packer build packer_template.json
==> ...................
==> amazon-ebs: Waiting for AMI to become ready...
==> amazon-ebs: Terminating the source AWS instance...
==> amazon-ebs: Cleaning up any extra volumes...
==> amazon-ebs: No volumes to clean up, skipping
==> amazon-ebs: Deleting temporary security group...
==> amazon-ebs: Deleting temporary keypair...
Build 'amazon-ebs' finished.

==> Builds finished. The artifacts of successful builds are:
--> amazon-ebs: AMIs were created:
us-east-1: ami-0b425c9074e5ac992
```


# Architecture

## Environment Setup

### JQ

[./jq](https://stedolan.github.io/jq/) is a lightweight and flexible command-line JSON processor.

```shell script
brew install jq
```


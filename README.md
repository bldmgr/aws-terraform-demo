&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;

# Ec2 Worker

This go project creates an Ec2 worker server to run bash commands from a remote docker container. 

# Quick start


## Create terraform state s3 bucket

```shell script
docker run -it $(docker build -q .)
cd ~
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=us-east-1
aws s3api create-bucket --bucket <S3_BUCKET_NAME> --region us-east-1
```

## Create custom AMI using Packer

To create a new AMI out of the latest Amazon Linux AMI, and install the server.go application within the new created AMI.
* Edit the terraform.tf file to update the bucket = "<S3_BUCKET_NAME>" 

```shell script
cd ec2-server
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


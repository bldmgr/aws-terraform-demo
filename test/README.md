cd ~
export AWS_ACCESS_KEY_ID=$aws_access_key
export AWS_SECRET_ACCESS_KEY=$aws_secret_access_key
export AWS_DEFAULT_REGION=$aws_default_region

aws s3api create-bucket --bucket $aws_s3_bucket_name --region us-east-1


terraform init


export aws_ami_id=ami-0ecb70e2947337d36
terraform apply -auto-approve -var "git_access_token=$store_git_access_token" -var "ami_id=$aws_ami_id"

export LOCAL_IP=$(curl http://ipv4.icanhazip.com)
export INSTANCE_IP=$(terraform output -json | jq -r '.instance_ip.value' )
export SECURITY_GROUP=$(terraform output -json | jq -r '.security_group.value'  )
export INSTANCE_ID=$(terraform output -json | jq -r '.instance_id.value'  )
aws ec2 authorize-security-group-ingress --group-name $SECURITY_GROUP --protocol tcp --port 8081 --cidr $LOCAL_IP/32
go run client.go $INSTANCE_IP:8081 'ls'

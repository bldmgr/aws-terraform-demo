terraform {
  backend "s3" {
    bucket = "tfstate-ec2-worker-test"
    key    = "main/terraform.tfstate"
    region = "us-east-1"
  }
}

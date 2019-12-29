terraform {
  backend "s3" {
    bucket = "tfstate-my-bucket"
    key    = "main/terraform.tfstate"
    region = "us-east-1"
  }
}

variable "git_access_token" {}
variable "ami_id" {}

provider "aws" {
  region = "us-east-1"
}

data "http" "my_ip" {
  url = "http://ipv4.icanhazip.com"
}

resource "aws_security_group" "ec2_worker_security_group" {
  name = "ec2-worker-securitygroup"
  description = "Allow tcp requests to a command server"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["${chomp(data.http.my_ip.body)}/32"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "tls_private_key" "ec2_worker_private_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "ec2_worker_key_pair" {
  key_name = "ec2-worker-keypair"
  public_key = "${tls_private_key.ec2_worker_private_key.public_key_openssh}"
}

resource "aws_instance" "go_docker_cmdserver_instance" {
  ami = "${var.ami_id}"
  instance_type = "t2.medium"
  tags = { Name = "concourse-ec2-worker" }
  security_groups = [ "${aws_security_group.ec2_worker_security_group.name}" ]
  key_name = "${aws_key_pair.ec2_worker_key_pair.key_name}"

  connection {
    host = self.public_ip
    type = "ssh"
    user = "ubuntu"
    private_key = "${tls_private_key.ec2_worker_private_key.private_key_pem}"
  }

  provisioner "remote-exec" {
    inline = [
      "export GIT_ACCESS_TOKEN=${var.git_access_token}",
      "nohup /usr/local/go/bin/go run /home/ubuntu/server.go &",
      "sleep 1",
    ]
  }
}

output "instance_ip" {
  value = "${aws_instance.go_docker_cmdserver_instance.public_ip}"
}

output "instance_id" {
  value = "${aws_instance.go_docker_cmdserver_instance.id}"
}

output "security_group" {
  value = "${aws_security_group.ec2_worker_security_group.name}"
}
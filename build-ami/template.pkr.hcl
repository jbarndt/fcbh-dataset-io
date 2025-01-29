variable "region" {
  type    = string
  default = "us-west-2"
}

variable "aws_profile" {
  type    = string
  default = "dbp-dev"
}
variable "vpc_id" {
  type    = string
  default = "vpc-0b6a6785e74d18db3"
}

variable "subnet_id" {
  type    = string
  default = "subnet-086d451813d884dd0"
}
variable "ami_name_prefix" {
  type    = string
  default = "deep-learning-golden-ami"
}

variable "instance_type" {
  type    = string
  default = "g6e.xlarge"
}

variable "source_ami" {
  type    = string
  default = "ami-013e597d66c833276"
}
variable "ssh_username" {
  type    = string
  default = "ubuntu"
}

packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "my-ami" {
  ami_name      = "${var.ami_name_prefix}-{{timestamp}}" 
  instance_type = var.instance_type 
  region        = var.region
  source_ami    = var.source_ami  
  ssh_username  = var.ssh_username

  ami_block_device_mappings {
    device_name = "/dev/sda1"
    volume_size = 128
    volume_type = "gp3"
    delete_on_termination = true
  }
  profile = var.aws_profile

  vpc_id           = var.vpc_id
  subnet_id        = var.subnet_id

  ami_description = "Custom deep-learning AMI built with Packer"
  tags = {
    Name = "MyCustomAMI"
  }
}

build {
  sources = [
    "source.amazon-ebs.my-ami"
  ]

  provisioner "shell" {
    script = "provision-ami.sh"
  }
}
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

# Note: The AMI is based on ubuntu 22.04
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

  ami_description = "Custom deep-learning AMI (ubuntu 22.04) built with Packer"
  tags = {
    Name = "Deep Learning"
  }
}

build {
  sources = [
    "source.amazon-ebs.my-ami"
  ]

  provisioner "file" {
    source      = "cloudwatch-nvidia.json"  
    destination = "/tmp/cloudwatch-nvidia.json"
  }
  provisioner "file" {
    source      = "env_vars_dev.txt"  
    destination = "/tmp/env_vars.txt"
  }
  provisioner "shell" {
      inline = [
        "echo 'Setting environment variables from file:'",
        // Read the file and set each variable
        "while IFS='=' read -r var val; do",
        "  echo \"Setting $var to $val\"",
        "  echo \"export $var='$val'\" >> /home/ubuntu/.bashrc",
        "done < /tmp/env_vars.txt",
      ]
  }

  # temporary, to verify file was copied with root permissions
  provisioner "shell" {
    inline = [
      "sudo mv /tmp/cloudwatch-nvidia.json  /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json",
      "sudo chmod 644 /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json",
      "echo 'export GOPATH=/home/ubuntu/go' >> /home/ubuntu/.bashrc",
      "echo 'export PATH=$PATH:$GOPATH/bin' >> /home/ubuntu/.bashrc"      
    ]
  }
  
  provisioner "shell" {
    script = "provision-ami.sh"
  }
}
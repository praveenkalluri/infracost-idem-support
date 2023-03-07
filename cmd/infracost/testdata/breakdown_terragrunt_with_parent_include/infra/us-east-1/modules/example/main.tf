variable "instance_type" {
  description = "The EC2 instance type for the web app"
  type        = string
}

variable "root_block_device_volume_size" {
  description = "The size of the root block device volume for the web app EC2 instance"
  type        = number
}

variable "block_device_volume_size" {
  description = "The size of the block device volume for the web app EC2 instance"
  type        = number
}

variable "block_device_iops" {
  description = "The number of IOPS for the block device for the web app EC2 instance"
  type        = number
}

resource "aws_instance" "web_app" {
  ami           = "ami-674cbc1e"
  instance_type = var.instance_type

  root_block_device {
    volume_size = var.root_block_device_volume_size
  }

  ebs_block_device {
    device_name = "my_data"
    volume_type = "io1"
    volume_size = var.block_device_volume_size
    iops        = var.block_device_iops
  }
}

variable "access_key" {
  description = "AWS access key"
}

variable "secret_key" {
  description = "AWS secret access key"
}

provider "aws" {
  access_key  = "${var.access_key}"
  secret_key  = "${var.secret_key}"
  region      = "${var.region}"
}

resource "aws_vpc" "primary" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_vpc" "app" {
  cidr_block = "10.1.0.0/16"
}

resource "exec", "peer_vpcs" {
  command = "commands/peer-vpc ${aws_vpc.primary.id} ${aws_vpc.app.id}"
  depends_on = ["aws_vpc.primary.id","aws_vpc.app"]
}

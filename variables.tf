variable "region" {
  default = "us-east-1"
}

variable "dynamodb_table_name" {
  default = "ShortUrls"
}

variable "aws_access_key" {
  default = "aws_access_key"
}

variable "aws_secret_key" {
  default = "aws_secret_key"
}

variable "dynamodb_endpoint" {
  default = "https://localhost.localstack.cloud:4566"
}

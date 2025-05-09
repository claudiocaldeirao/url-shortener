terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }

  required_version = ">= 1.0"
}

provider "aws" {
  access_key                  = var.aws_access_key
  secret_key                  = var.aws_secret_key
  region                      = var.region
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    lambda     = "http://localhost:4566"
    iam        = "http://localhost:4566"
    dynamodb   = "http://localhost:4566"
    apigateway = "http://localhost:4566"
  }
}

locals {
  binary_name  = "bootstrap"
  src_path     = "${path.module}/app/cmd/main.go"
  binary_path  = "${path.module}/build/bootstrap"
  archive_path = "${path.module}/build/lambda.zip"
}

// build the binary for the lambda function in a specified path
resource "null_resource" "function_binary" {
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}"
  }
}

// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "function_archive" {
  depends_on = [null_resource.function_binary]

  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

# IAM Role for Lambda
resource "aws_iam_role" "lambda_exec_role" {
  name = "url-shortener-lambda-exec"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

# Lambda Function
resource "aws_lambda_function" "url_shortener_lambda" {
  function_name    = "url-shortener-lambda"
  handler          = local.binary_name
  runtime          = "go1.x"
  source_code_hash = data.archive_file.function_archive.output_base64sha256
  role             = aws_iam_role.lambda_exec_role.arn
  filename         = local.archive_path

  environment {
    variables = {
      AWS_DYNAMO_DB_TABLE    = var.dynamodb_table_name
      AWS_DYNAMO_DB_ENDPOINT = var.dynamodb_endpoint
      AWS_REGION             = var.region
    }
  }
}

# IAM Policy to access DynamoDB
resource "aws_iam_role_policy" "lambda_policy" {
  name = "url-shortener-lambda-policy"
  role = aws_iam_role.lambda_exec_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem"
        ],
        Resource = aws_dynamodb_table.short_urls_table.arn
      },
      {
        Effect : "Allow",
        Action : [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource : "arn:aws:logs:*:*:*"
      }
    ]
  })
}

# DynamoDB Table
resource "aws_dynamodb_table" "short_urls_table" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "Hash"
    type = "S"
  }

  hash_key = "Hash"
}

# API Gateway REST API
resource "aws_api_gateway_rest_api" "api" {
  name = "url-shortener-api"
}

# Resource for /shorten
resource "aws_api_gateway_resource" "shorten" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "shorten"
}

# Resource for /{hash}
resource "aws_api_gateway_resource" "shortcode" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = "{hash}"
}

# POST /shorten method
resource "aws_api_gateway_method" "post_shorten" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.shorten.id
  http_method   = "POST"
  authorization = "NONE"
}

# GET /{hash} method
resource "aws_api_gateway_method" "get_redirect" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.shortcode.id
  http_method   = "GET"
  authorization = "NONE"
}

# Lambda Integration for POST /shorten
resource "aws_api_gateway_integration" "post_shorten_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.shorten.id
  http_method             = aws_api_gateway_method.post_shorten.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.url_shortener_lambda.invoke_arn
}

# Lambda Integration for GET /{hash}
resource "aws_api_gateway_integration" "get_redirect_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.shortcode.id
  http_method             = aws_api_gateway_method.get_redirect.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.url_shortener_lambda.invoke_arn
}

# Lambda permissions for API Gateway to invoke
resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.url_shortener_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api.execution_arn}/*/*"
}

# Deployment
resource "aws_api_gateway_deployment" "deployment" {
  depends_on = [
    aws_api_gateway_integration.post_shorten_integration,
    aws_api_gateway_integration.get_redirect_integration
  ]
  rest_api_id = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_stage" "dev" {
  stage_name    = "dev"
  rest_api_id   = aws_api_gateway_rest_api.api.id
  deployment_id = aws_api_gateway_deployment.deployment.id
  variables = {
    env = "development"
  }
}

# URL Shortener Serverless API (Go + AWS)

A lightweight, fully serverless URL shortener built in Go, powered by AWS Lambda, DynamoDB, and API Gateway.
Includes in-memory caching for the top 10 most accessed URLs to ensure ultra-fast redirects.
Infrastructure managed via Terraform for easy deployment.

## ✨ Features

- Shorten any URL with a simple POST request

- Fast redirects using DynamoDB and Lambda memory cache

- Tracks click counts per URL

- Serverless architecture (zero server maintenance)

- Infrastructure-as-Code with Terraform

## 🏗️ Architecture

- Go: Lambda function runtime

- AWS Lambda: compute logic

- DynamoDB: persistent URL storage and click tracking

- API Gateway: HTTP interface

- Terraform: deploys all AWS resources

## 🚀 Getting Started

Prerequisites

- Go installed (>= 1.19)

- Terraform installed (>= 1.0)

- AWS CLI configured (aws configure)

### 1. Clone the repository

```
git clone https://github.com/your-username/url-shortener-serverless.git
cd url-shortener-serverless
```

### 2. Build the Lambda function

```
cd lambda
GOOS=linux GOARCH=amd64 go build -o main main.go
zip main.zip main
cd ..
```

### 3. Deploy with Terraform

```
terraform init
terraform apply
```

Confirm with yes when prompted.

## 📬 API Usage

### 1. Shorten a URL

```
// Request

POST /shorten
Content-Type: application/json

{
  "url": "https://example.com"
}

// Response
{
  "shortCode": "abc123"
}
```

### 2. Redirect to the original URL

```
// Request

GET /{shortcode}

Will HTTP 302 Redirect to the original URL.
```

## 📦 Project Structure

```
url-shortener/
├── app/             # Go source code
│   ├── main.go
│   └── utils
│       └── utils.go # hash method to generate short url
├── main.tf          # Terraform AWS resources
├── variables.tf     # Terraform variables
└── README.md
```

## 🛡 License

This project is open-source and available under the [MIT License](LICENSE).

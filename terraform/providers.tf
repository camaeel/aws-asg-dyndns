terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    null = {
      source = "hashicorp/null"
      version = "~> 3.1"
    }
    http = {
      source = "hashicorp/http"
      version = "~> 2.1"
    }
  }
}

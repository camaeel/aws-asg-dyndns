variable "name" {
  type = string
  description = "AWS dynamic DNS lambda name"
}

variable "tags" {
  type = map(string)
  description = "List of tags to apply on all resources"
  default = {}
}

variable "log_retention" {
  type = number
  description = "Log retention for lambda function (in days)"
  default = 7
}

variable "logs_kms_key_id" {
  type = string
  description = "Log encryption KMS key"
  default = null
}

variable "sqs_kms_key_id" {
  type = string
  description = "Sqs queue encryption key. If left as null, SSE-SQS key is used"
  default = null
}

variable "dns_provider" {
  type = string
  description = "Name of dns provider. Possible values: cloudflare"
  default = "cloudflare"

  validation {   
    condition     = can(regex("^(cloudflare)$", var.dns_provider))    
    error_message = "Valid values are: cloudflare."  
  }
}

variable "zone_name" {
  type = string
  description = "Zone name"
  validation {   
    condition     = can(regex("^[a-z0-9-_.]+$", var.zone_name))    
    error_message = "Should be a valid dns name."  
  }
}

variable "lambda_timeout" {
  type = number
  default = 20
  description = "Timeout for lambda function execution"
}

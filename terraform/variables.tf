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

variable "dns_providers" {

  type = map(object({
    private_zone  = bool
    provider_name = string
    zone_name     = string
  }))

  description = "Map zone_name => provider"
  default = {}
}

variable "lambda_timeout" {
  type = number
  default = 10
  description = "Timeout for lambda function execution"
}

variable "github_lambda_version" {
  type = string
  default = "pre0.0.15"
  description = "Github lambda function version. If set to null will try to upload local zip file (from a directory one level above the module)"
}

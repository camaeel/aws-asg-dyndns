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

variable "tags" {
  type = map(string)
  description = "Tags"
}

variable "provider_name" {
  type = string
  description = "Provider name. Possible values: cloudflare, route53"
  
  validation {   
    condition     = can(regex("^(cloudflare|route53)$", var.provider_name))    
    error_message = "Valid values are: cloudflare route53."  
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

variable "lambda_role" {
  type = string
  description = "AWS role to grant access to SSM parameters"
}

variable "role_prefix_name" {
  type = string
  description = "Prefix role names"
}

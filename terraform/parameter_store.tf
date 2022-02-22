resource "aws_ssm_parameter" "cloudflare-token" {
  name        = "/dyn-dns/${var.zone_name}/cloudflare/token"
  description = "Cloudflare's token for AWS dynamic DNS"
  type        = "SecureString"
  value       = cloudflare_api_token.aws-dyn-dns[count.index].value
  tier        = "Standard"
  count = var.dns_provider == "cloudflare" ? 1 : 0

  tags = merge(
    var.tags,
    { Name = "/dyn-dns/${var.zone_name}/cloudflare/token" }
  ) 
}

resource "aws_ssm_parameter" "zone-provider" {
  name        = "/dyn-dns/${var.zone_name}/provider"
  description = "Cloudflare's token for AWS dynamic DNS"
  type        = "String"
  value       = var.dns_provider
  tier        = "Standard"

  tags = merge(
    var.tags,
    { Name = "/dyn-dns/${var.zone_name}/provider" }
  ) 
}

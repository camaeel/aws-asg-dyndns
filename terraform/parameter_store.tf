resource "aws_ssm_parameter" "cloudflare-token" {
  name        = "/dyn-dns/cloudflare/${var.zone_name}/token"
  description = "Cloudflare's token for AWS dynamic DNS"
  type        = "SecureString"
  value       = cloudflare_api_token.aws-dyn-dns[count.index].value
  tier        = "Standard"
  count = var.zone_name != "" ? 1 : 0

  tags = merge(
    var.tags,
    { Name = "/dyn-dns/cloudflare/${var.zone_name}/token" }
  ) 
}

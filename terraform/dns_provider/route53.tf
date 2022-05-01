data "aws_route53_zone" "zone" {
  name         = var.zone_name
  private_zone = var.private_zone

  count = var.provider_name == "route53" ? 1 : 0
}

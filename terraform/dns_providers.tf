module "dns_provider" {
  source = "./dns_provider"

  tags = var.tags

  provider_name = each.value.provider_name
  zone_name = each.value.zone_name
  lambda_role = aws_iam_role.lambda.name
  role_prefix_name = var.name
  private_zone = each.value.private_zone

  for_each = var.dns_providers
}

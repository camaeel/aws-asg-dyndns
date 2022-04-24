module "dns_provider" {
  source = "./dns_provider"

  tags = var.tags

  provider_name = each.value
  zone_name = each.key
  lambda_role = aws_iam_role.lambda.name
  role_prefix_name = var.name

  for_each = var.dns_providers
}

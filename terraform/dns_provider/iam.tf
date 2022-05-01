# SSM Parameter Store - Cloudflare
resource "aws_iam_role_policy_attachment" "lambda-ssm-parameter-store" {
  role       = var.lambda_role
  policy_arn = aws_iam_policy.lambda-ssm-parameter-store.arn
}

resource "aws_iam_policy" "lambda-ssm-parameter-store" {
  name  = "${var.role_prefix_name}-ssm-${var.provider_name}-${var.zone_name}"

  policy = data.aws_iam_policy_document.lambda-ssm-parameter-store.json

  tags = merge(
    var.tags,
    { Name = "${var.role_prefix_name}-ssm-${var.provider_name}-${var.zone_name}" }
  ) 
}

data "aws_iam_policy_document" "lambda-ssm-parameter-store" {
  statement {
    sid = "describeParams"
    effect = "Allow"
    actions = [
      "ssm:DescribeParameters",
    ]

    resources = [
      "*",
    ]
  }

  dynamic "statement" {
    for_each = aws_ssm_parameter.cloudflare-token
    content {
      sid = "accessSsmCfToken"
      effect = "Allow"
      actions = [
        "ssm:GetParameter",
      ]

      resources = [
        statement.value.arn,
      ]      
    }
  }

  dynamic "statement" {
    for_each = data.aws_route53_zone.zone
    content {
      sid = "accessSsmRoute53Access"
      effect = "Allow"
      actions = [
        "route53:ChangeResourceRecordSets",
      ]

      resources = [
        statement.value.arn,
      ]      
    }
  }


  statement {
    sid = "accessSsmZoneProvider"
    effect = "Allow"
    actions = [
      "ssm:GetParameter",
    ]

    resources = [
      aws_ssm_parameter.zone-provider.arn,
    ]
  }
}

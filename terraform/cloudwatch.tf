resource "aws_cloudwatch_log_group" "logs" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = var.log_retention
  kms_key_id = var.logs_kms_key_id

  tags = merge(
    var.tags,
    { Name = "${var.name}-role" }
  ) 
}

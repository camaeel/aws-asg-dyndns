resource "aws_lambda_function" "aws-asg-dyndns" {
  filename      = "${path.module}/../main.zip"
  function_name = var.name
  role          = aws_iam_role.aws-asg-dyndns.arn
  handler       = "main"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("${path.module}/../main.zip")

  runtime = "go1.x"

  # environment {
  # variables = {
  #    foo = "bar"
  #  }
  #}
  tags = merge(
    var.tags,
    { Name = var.name }
  ) 
  # TODO publish function
}

resource "aws_lambda_event_source_mapping" "aws-asg-dyndns-sqs" {
  event_source_arn = aws_sqs_queue.aws-asg-dyndns-events.arn
  function_name    = aws_lambda_function.aws-asg-dyndns.arn

  depends_on = [
    aws_iam_role_policy_attachment.aws-asg-dyndns-sqs
  ]
}

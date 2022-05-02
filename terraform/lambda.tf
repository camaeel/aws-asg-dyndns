resource "null_resource" "lambda_zip_download" {
  count = local.download_from_github ? 1 : 0
  triggers = {
    on_version_change = "${var.github_lambda_version}"
  }

  provisioner "local-exec" {
    command = "curl -o main.zip ${local.download_from_github_url}"
  }
}

data "http" "lambda_zip_sha256" {
  count = local.download_from_github ? 1 : 0
  url   = local.sha256_from_github_url
}


resource "aws_lambda_function" "lambda" {
  filename      = local.lambda_filename 
  function_name = var.name
  role          = aws_iam_role.lambda.arn
  handler       = "main"

  source_code_hash = local.download_from_github ? data.http.lambda_zip_sha256[0].body : filebase64sha256("${path.module}/../main.zip")

  runtime = "go1.x"

  timeout = var.lambda_timeout
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

  depends_on = [
    null_resource.lambda_zip_download
  ]
}

resource "aws_lambda_event_source_mapping" "lambda-sqs" {
  event_source_arn = aws_sqs_queue.events.arn
  function_name    = aws_lambda_function.lambda.arn

  depends_on = [
    aws_iam_role_policy_attachment.lambda-sqs
  ]
}

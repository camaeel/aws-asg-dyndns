locals {
  download_from_github = (var.github_lambda_version != "")
  lambda_filename = local.download_from_github ? "${path.module}/main.zip" : "${path.module}/../main.zip"
  download_from_github_url = "https://github.com/camaeel/aws-asg-dyndns/releases/download/${var.github_lambda_version}/aws-asg-dyndns-linux-amd64.zip"
  sha256_from_github_url = "https://github.com/camaeel/aws-asg-dyndns/releases/download/${var.github_lambda_version}/aws-asg-dyndns-linux-amd64.zip.sha256"
}

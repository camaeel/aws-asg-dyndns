# Lambda function role
resource "aws_iam_role" "lambda" {
  name = "${var.name}-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-role" }
  ) 
}

# Logs
resource "aws_iam_role_policy_attachment" "lambda-logs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda-logs.arn
}

resource "aws_iam_policy" "lambda-logs" {
  name = "${var.name}-logs"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "${aws_cloudwatch_log_group.logs.arn}:*"
        }
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-logs" }
  ) 
}

# SQS
resource "aws_iam_role_policy_attachment" "lambda-sqs" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda-sqs.arn
}

resource "aws_iam_policy" "lambda-sqs" {
  name = "${var.name}-sqs"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "sqs:ReceiveMessage",
            "sqs:DeleteMessage",
            "sqs:GetQueueAttributes"
          ],
          "Resource": "${aws_sqs_queue.events.arn}" 
        }
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-sqs" }
  ) 
}

#EC2 

resource "aws_iam_role_policy_attachment" "lambda-ec2-access" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda-ec2-access.arn
}

resource "aws_iam_policy" "lambda-ec2-access" {
  name = "${var.name}-ec2-access"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "ec2:DescribeNetworkInterfaces"
          ],
          "Resource": "*"         
        },
        {
          "Effect": "Allow",
          "Action": [
            "ec2:CreateTags"          
          ],
          "Resource": "*",
          "Condition": {
            "StringEquals": {
                "ec2:ResourceTag/aws-asg-dynds": "aws-asg-dynds"
            }
          }          
        },
        {
          "Effect": "Allow",
          "Action": [

            "ec2:DescribeTags"
          ],
          "Resource": "*"
        }       
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-ec2-read" }
  ) 
}

resource "aws_iam_role_policy_attachment" "lambda-asg-lifecycle" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda-asg-lifecycle.arn
}

resource "aws_iam_policy" "lambda-asg-lifecycle" {
  name = "${var.name}-asg-lifecycle"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "autoscaling:CompleteLifecycleAction"
          ],
          "Resource": "*",
          "Condition": {
              "StringEquals": {
                  "autoscaling:ResourceTag/aws-asg-dynds": "aws-asg-dynds"
              }
          }          
        }
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-asg-lifecycle" }
  ) 
}

# SSM Parameter Store
resource "aws_iam_role_policy_attachment" "lambda-ssm-paramter-store" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda-ssm-paramter-store.arn
  count      = var.zone_name != "" ? 1 : 0
}

resource "aws_iam_policy" "lambda-ssm-paramter-store" {
  name  = "${var.name}-ssm-paramter-store"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "ssm:DescribeParameters"
          ],
          "Resource": "*" 
        },
        {
          "Effect": "Allow",
          "Action": [
            "ssm:GetParameter"
          ],
          "Resource": "${aws_ssm_parameter.zone-provider.arn}" 
        }
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-ssm-paramter-store" }
  ) 
}


resource "aws_iam_role_policy_attachment" "lambda-ssm-paramter-store-cloudflare" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda-ssm-paramter-store-cloudflare[count.index].arn
  count      = var.zone_name != "" ? 1 : 0
}

resource "aws_iam_policy" "lambda-ssm-paramter-store-cloudflare" {
  name  = "${var.name}-ssm-paramter-store-cloudflare"
  count = var.dns_provider == "cloudflare" ? 1 : 0

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "ssm:GetParameter"
          ],
          "Resource": "${aws_ssm_parameter.cloudflare-token[count.index].arn}" 
        }
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-ssm-paramter-store-cloudflare" }
  ) 
}

######################## Role for writing the queue for asg lifecycyle hook
resource "aws_iam_role" "aws-asg-dyndns-sqs-writer" {
  name = "${var.name}-sqs-writer-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "autoscaling.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-sqs-writer-role" }
  ) 
}

resource "aws_iam_role_policy_attachment" "aws-asg-dyndns-sqs-writer" {
  role       = aws_iam_role.aws-asg-dyndns-sqs-writer.name
  policy_arn = aws_iam_policy.aws-asg-dyndns-sqs-writer.arn
}


resource "aws_iam_policy" "aws-asg-dyndns-sqs-writer" {
  name = "${var.name}-sqs-writer"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "sqs:SendMessage",
            "sqs:GetQueueUrl"
          ],
          "Resource": "${aws_sqs_queue.events.arn}" 
        }
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-sqs-writer" }
  ) 
}

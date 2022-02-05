resource "aws_iam_role" "aws-asg-dyndns" {
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
resource "aws_iam_role_policy_attachment" "aws-asg-dyndns-logs" {
  role       = aws_iam_role.aws-asg-dyndns.name
  policy_arn = aws_iam_policy.aws-asg-dyndns-logs.arn
}

resource "aws_iam_policy" "aws-asg-dyndns-logs" {
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
resource "aws_iam_role_policy_attachment" "aws-asg-dyndns-sqs" {
  role       = aws_iam_role.aws-asg-dyndns.name
  policy_arn = aws_iam_policy.aws-asg-dyndns-sqs.arn
}

resource "aws_iam_policy" "aws-asg-dyndns-sqs" {
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
          "Resource": "${aws_sqs_queue.aws-asg-dyndns-events.arn}" 
        }
    ]
}
EOF

  tags = merge(
    var.tags,
    { Name = "${var.name}-sqs" }
  ) 
}

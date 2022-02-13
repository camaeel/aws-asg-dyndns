output "sqs_arn" {
  value = aws_sqs_queue.aws-asg-dyndns-events.arn
}

output "sqs_writer_role_arn" {
  value = aws_iam_role.aws-asg-dyndns-sqs-writer.arn
}

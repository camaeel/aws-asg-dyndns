resource "aws_sqs_queue" "events" {
  name                      = "${var.name}-events"
  delay_seconds             = 0 # default
  max_message_size          = 2048
  message_retention_seconds = 300
  receive_wait_time_seconds = 20

  # Encryption
  sqs_managed_sse_enabled = (var.sqs_kms_key_id ==null ? true : false)
  kms_master_key_id       = var.sqs_kms_key_id

  tags = merge(
    var.tags,
    { Name = var.name }
  ) 
}

resource "aws_sqs_queue_policy" "events" {
  queue_url = aws_sqs_queue.events.id

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Id": "allow-send-for-asg-events",
    "Statement": [{
        "Effect": "Allow",
        "Principal": {
            "Service": ["autoscaling.amazonaws.com", "sqs.amazonaws.com"]
        },
        "Action": "sqs:SendMessage",
        "Resource": ["${aws_sqs_queue.events.arn}"]
    }]
}
POLICY
}
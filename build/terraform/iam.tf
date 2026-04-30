# GHIN-API IAM Role and Policy Attachments
resource "aws_iam_role" "role_ghin_api" {
  name = "ghin-api-lambda-runner-${terraform.workspace}"

  assume_role_policy = jsonencode({
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
  })

  tags = var.common_tags
}

resource "aws_iam_policy" "policy_ghin_api_runner" {
  name   = "ghin-policy-api-runner-${terraform.workspace}"
  description = "Allow access to services for APIs"
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
           "Effect" : "Allow",
           "Action" : ["secretsmanager:*"],
           "Resource" : "*"
      },
      {
           "Effect" : "Allow",
           "Action" : ["cloudwatch:*"],
           "Resource" : "*"
      },
      {
           "Effect" : "Allow",
           "Action" : [
             "dynamodb:GetItem",
             "dynamodb:PutItem"
           ],
           "Resource" : aws_dynamodb_table.table_ghin_api_cache.arn
      },
    ]
  })
  tags = var.common_tags
}

resource "aws_iam_role_policy_attachment" "ghin_api_lambda_ap_policy_attachment" {
  role = aws_iam_role.role_ghin_api.id
  policy_arn = aws_iam_policy.policy_ghin_api_runner.arn
}

# API-Gateway-Cloudwatch-Logging -- Allows API Gateway to write logs to CloudWatch
resource "aws_iam_role" "role_api_gateway_cloudwatch_logging" {
  name = "api-gateway-cloudwatch-logging-${terraform.workspace}"

  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "apigateway.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  })

  tags = var.common_tags
}
################################################################################
# API Definition

resource "aws_api_gateway_rest_api" "rest_api_ghin" {
  name = "ghin-api-${terraform.workspace}"
  description = "Rest API for accessing GHIN data"
  api_key_source = "HEADER"
  tags = var.common_tags
}


# Permission for API Gateway to invoke the GHIN lambda
resource "aws_lambda_permission" "lambda_permission_ghin" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_ghin_api[count.index].function_name
  principal     = "apigateway.amazonaws.com"
  count = terraform.workspace == "prod" ? 1 : 0
}

resource "aws_wafv2_web_acl" "WebService_APIWebACL" {
  depends_on = []
  name        = "ghin-api-web-acl-${terraform.workspace}"
  description = "WAF Web ACL for the GHIN API"
  scope       = "REGIONAL"
  
  default_action {
    allow {}
  }

  visibility_config {
    cloudwatch_metrics_enabled = false
    metric_name                = "ghin-api-web-acl"
    sampled_requests_enabled   = false
  }
}

resource "aws_wafv2_web_acl_association" "ghin_rest_api_web_acl_association" {
  resource_arn = aws_api_gateway_stage.stage-ghin-api.arn
  web_acl_arn  = aws_wafv2_web_acl.WebService_APIWebACL.arn
}


# Configure Stages, including WAF support
resource "aws_api_gateway_deployment" "WebService_API_Gateway_Deployment" {
  rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id

  triggers = {
      redeployment = sha1(jsonencode([
          file("apigateway_ghin.tf"),
      ]))
  }
  depends_on = [
    aws_api_gateway_integration.integration_courses_get
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_method_settings" "settings_all" {
  rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
  stage_name  = aws_api_gateway_stage.stage-ghin-api.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled = false
    # optional but commonly enabled alongside metrics:
    # logging_level    = "INFO"
    # data_trace_enabled = false
  }
}

resource "aws_api_gateway_stage" "stage-ghin-api" {
  deployment_id = aws_api_gateway_deployment.WebService_API_Gateway_Deployment.id
  rest_api_id   = aws_api_gateway_rest_api.rest_api_ghin.id
  stage_name    = lower(terraform.workspace)

  variables = {

  }
}

// API Keys
resource "aws_api_gateway_usage_plan" "api_key_usage_plan_ghin" {
  name        = "ghin-usage-plan-${terraform.workspace}"
  description = "usage plan for GHIN API"

  api_stages {
    api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    stage = aws_api_gateway_stage.stage-ghin-api.stage_name
  }

  throttle_settings {
    burst_limit = 5
    rate_limit  = 10
  }
}

resource "aws_api_gateway_api_key" "api_key_LeagueManager" {
  name = "ghin-api-key-leaguemanager-${terraform.workspace}"
}

resource "aws_api_gateway_usage_plan_key" "usage_plan_key_query" {
  key_id        = aws_api_gateway_api_key.api_key_LeagueManager.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.api_key_usage_plan_ghin.id
}

# resource "aws_api_gateway_account" "account_cloudwatch_logging" {
#   cloudwatch_role_arn = aws_iam_role.role_api_gateway_cloudwatch_logging.arn
# }

resource "aws_api_gateway_gateway_response" "response_401_unauthorized" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api_ghin.id
  status_code   = "401"
  response_type = "UNAUTHORIZED"

  response_templates = {
    "application/json" = "{\"errorMessage\":\"unauthorized to perform this ghin operation\",\"status\":\"error\"}"
  }

  response_parameters = {
    "gatewayresponse.header.Authorization" = "'Basic'"
  }
}

################################################################################
# Endpoints and Resources

# /ghin/
resource "aws_api_gateway_resource" "resource_crm_api_ghin" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    parent_id = aws_api_gateway_rest_api.rest_api_ghin.root_resource_id
    path_part = "ghin"
}

# /ghin/courses
resource "aws_api_gateway_resource" "resource_crm_api_ghin_courses" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    parent_id = aws_api_gateway_resource.resource_crm_api_ghin.id
    path_part = "courses"
}

resource "aws_api_gateway_method" "method_courses_get" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_courses.id
    http_method = "GET"
    authorization = "NONE"
    api_key_required = true
}

resource "aws_api_gateway_integration" "integration_courses_get" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_courses.id
    http_method = aws_api_gateway_method.method_courses_get.http_method
    # POST is needed for lambda integrations. See Lambda Proxy example https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.lambda_ghin_api[count.index].invoke_arn
    count = terraform.workspace == "prod" ? 1 : 0
}

module "api_gateway_enable_cors_ghin_api_courses" {
  source  = "squidfunk/api-gateway-enable-cors/aws"
  version = "0.3.3"
  api_id          = aws_api_gateway_rest_api.rest_api_ghin.id
  api_resource_id = aws_api_gateway_resource.resource_crm_api_ghin_courses.id
  allow_headers = ["Authorization", "Content-Type", "X-Amz-Date", "X-Amz-Security-Token", "X-Api-Key", "environment"]
}

# /ghin/courses/{id}
resource "aws_api_gateway_resource" "resource_crm_api_ghin_courses_id" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    parent_id = aws_api_gateway_resource.resource_crm_api_ghin_courses.id
    path_part = "{id}"
}

resource "aws_api_gateway_method" "method_courses_get_id" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_courses_id.id
    http_method = "GET"
    authorization = "NONE"
    api_key_required = true
}

resource "aws_api_gateway_integration" "integration_courses_get_id" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_courses_id.id
    http_method = aws_api_gateway_method.method_courses_get_id.http_method
    # POST is needed for lambda integrations. See Lambda Proxy example https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.lambda_ghin_api[count.index].invoke_arn
    count = terraform.workspace == "prod" ? 1 : 0
}

module "api_gateway_enable_cors_ghin_api_courses_id" {
  source  = "squidfunk/api-gateway-enable-cors/aws"
  version = "0.3.3"
  api_id          = aws_api_gateway_rest_api.rest_api_ghin.id
  api_resource_id = aws_api_gateway_resource.resource_crm_api_ghin_courses_id.id
  allow_headers = ["Authorization", "Content-Type", "X-Amz-Date", "X-Amz-Security-Token", "X-Api-Key", "environment"]
}

# /ghin/golfers/{id}
resource "aws_api_gateway_resource" "resource_crm_api_ghin_golfers" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    parent_id = aws_api_gateway_resource.resource_crm_api_ghin.id
    path_part = "golfers"
}

resource "aws_api_gateway_resource" "resource_crm_api_ghin_golfers_id" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    parent_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers.id
    path_part = "{id}"
}

resource "aws_api_gateway_method" "method_golfers_get_id" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id.id
    http_method = "GET"
    authorization = "NONE"
    api_key_required = true
}

resource "aws_api_gateway_integration" "integration_golfers_get_id" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id.id
    http_method = aws_api_gateway_method.method_golfers_get_id.http_method
    # POST is needed for lambda integrations. See Lambda Proxy example https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.lambda_ghin_api[count.index].invoke_arn
    count = terraform.workspace == "prod" ? 1 : 0
}

module "api_gateway_enable_cors_ghin_api_golfers_id" {
  source  = "squidfunk/api-gateway-enable-cors/aws"
  version = "0.3.3"
  api_id          = aws_api_gateway_rest_api.rest_api_ghin.id
  api_resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id.id
  allow_headers = ["Authorization", "Content-Type", "X-Amz-Date", "X-Amz-Security-Token", "X-Api-Key", "environment"]
}

# /ghin/golfers/{id}/handicaps
resource "aws_api_gateway_resource" "resource_crm_api_ghin_golfers_id_handicaps" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    parent_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id.id
    path_part = "handicaps"
}

resource "aws_api_gateway_method" "method_golfers_get_id_handicaps" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id_handicaps.id
    http_method = "GET"
    authorization = "NONE"
    api_key_required = true
}

resource "aws_api_gateway_integration" "integration_golfers_get_id_handicaps" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id_handicaps.id
    http_method = aws_api_gateway_method.method_golfers_get_id_handicaps.http_method
    # POST is needed for lambda integrations. See Lambda Proxy example https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.lambda_ghin_api[count.index].invoke_arn
    count = terraform.workspace == "prod" ? 1 : 0
}

module "api_gateway_enable_cors_ghin_api_golfers_id_handicaps" {
  source  = "squidfunk/api-gateway-enable-cors/aws"
  version = "0.3.3"
  api_id          = aws_api_gateway_rest_api.rest_api_ghin.id
  api_resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id_handicaps.id
  allow_headers = ["Authorization", "Content-Type", "X-Amz-Date", "X-Amz-Security-Token", "X-Api-Key", "environment"]
}

# /ghin/golfers/{id}/revisions
resource "aws_api_gateway_resource" "resource_crm_api_ghin_golfers_id_revisions" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    parent_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id.id
    path_part = "revisions"
}

resource "aws_api_gateway_method" "method_golfers_get_id_revisions" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id_revisions.id
    http_method = "GET"
    authorization = "NONE"
    api_key_required = true
}

resource "aws_api_gateway_integration" "integration_golfers_get_id_revisions" {
    rest_api_id = aws_api_gateway_rest_api.rest_api_ghin.id
    resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id_revisions.id
    http_method = aws_api_gateway_method.method_golfers_get_id_revisions.http_method
    # POST is needed for lambda integrations. See Lambda Proxy example https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_integration
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = aws_lambda_function.lambda_ghin_api[count.index].invoke_arn
    count = terraform.workspace == "prod" ? 1 : 0
}

module "api_gateway_enable_cors_ghin_api_golfers_id_revisions" {
  source  = "squidfunk/api-gateway-enable-cors/aws"
  version = "0.3.3"
  api_id          = aws_api_gateway_rest_api.rest_api_ghin.id
  api_resource_id = aws_api_gateway_resource.resource_crm_api_ghin_golfers_id_revisions.id
  allow_headers = ["Authorization", "Content-Type", "X-Amz-Date", "X-Amz-Security-Token", "X-Api-Key", "environment"]
}
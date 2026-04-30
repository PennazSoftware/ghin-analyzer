# GHIN-API
data "archive_file" "zip_ghin_api" {
   type        = "zip"
   source_dir  = "${var.code_location}/ghin-api/${terraform.workspace}/"
   output_path = "${var.code_location}/ghin-api/${terraform.workspace}/bootstrap.zip"
}

resource "aws_lambda_function" "lambda_ghin_api" {
   filename         = "${var.code_location}/ghin-api/${terraform.workspace}/bootstrap.zip"
   description      = "Rest API for GHIN"
   function_name    = "ghin-api-${terraform.workspace}"
   role             = aws_iam_role.role_ghin_api.arn
   source_code_hash = filebase64sha256("${var.code_location}/ghin-api/${terraform.workspace}/bootstrap.zip")
   handler          = "main"
   timeout          = "60"
   memory_size = "2048"
   runtime = "provided.al2023"
   architectures = ["arm64"]
   tags = var.common_tags
   count = terraform.workspace == "prod" ? 1 : 0

   environment {
     variables = {
         ENVIRONMENT = "${terraform.workspace}"
     }
   }
}

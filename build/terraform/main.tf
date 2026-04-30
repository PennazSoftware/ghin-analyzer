terraform {
 required_providers {
   aws = {
     source = "hashicorp/aws"
     version = ">= 5.42"
   }
 }

   backend "s3" {
    profile        = "pennaz"
    bucket         = "handicap-aws-terraform-state"
    key            = "terraform.tfstate"
    workspace_key_prefix = "handicap"
    region         = "us-west-2"
    dynamodb_table = "handicap-aws-terraform-lock"
  }
}
    
provider "aws" {
  region = "us-west-2"
  shared_credentials_files = ["$HOME/.aws/credentials"]
  profile = "pennaz"
}

variable "code_location" {
  type        = string
  description = "location of the zipped golang executables"
  default = "../bin"
}


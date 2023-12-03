provider "aws" {
  region = "us-east-1" # Change this to your desired AWS region
  default_tags {
    tags = {
      Management = "Terraform"
    }
  }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.12.0"
    }
  }
}

terraform {
  backend "s3" {
    # Replace this with your bucket name!
    bucket = "terraform-state-21151"
    key    = "global/s3/okta-event-hooks.tfstate"
    region = "us-east-1"

    # Replace this with your DynamoDB table name!
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}
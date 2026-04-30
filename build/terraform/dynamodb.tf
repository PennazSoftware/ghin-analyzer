################################################################################
# DynamoDB Tables
################################################################################

resource "aws_dynamodb_table" "table_handicap_golfers" {
    name = format("handicap-golfers-%s", lower(terraform.workspace))
    billing_mode = "PAY_PER_REQUEST"
    hash_key = "golferID"

    attribute {
        name = "golferID"
        type = "N"
    }

    tags = var.common_tags
}

resource "aws_dynamodb_table" "table_handicap_scores" {
    name           = format("handicap-scores-%s", lower(terraform.workspace))
    billing_mode   = "PAY_PER_REQUEST"
    hash_key       = "id"

    attribute {
        name = "id"
        type = "N"
    }

    attribute {
        name = "golfer_id"
        type = "N"
    }

    attribute {
        name = "played_at"
        type = "S"
    }

    global_secondary_index {
        name               = "golfer_id-index"
        key_schema           {
            attribute_name = "golfer_id"
            key_type       = "HASH"
        }
        projection_type    = "ALL"
    }

    global_secondary_index {
        name               = "played_at-index"
        key_schema           {
            attribute_name = "played_at"
            key_type       = "HASH"
        }
        projection_type    = "ALL"
    }

    tags = var.common_tags
}

resource "aws_dynamodb_table" "table_ghin_api_cache" {
    name         = format("ghin-api-cache-%s", lower(terraform.workspace))
    billing_mode = "PAY_PER_REQUEST"
    hash_key     = "cacheKey"

    attribute {
        name = "cacheKey"
        type = "S"
    }

    ttl {
        attribute_name = "ttl"
        enabled        = true
    }

    tags = var.common_tags
}

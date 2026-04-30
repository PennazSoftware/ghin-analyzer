# Copy the suspicion-scoring.html file to s3://prod.pennaz.com/suspicion_scoring.html
aws s3 cp suspicion-scoring.html s3://prod.pennaz.com/suspicion_scoring.html --profile pennaz --region us-west-2 --no-cli-pager

# Invalidate the CloudFront distribution for the apidocs path
aws cloudfront create-invalidation --distribution-id E1E98PZICJGP8H --paths "/suspicion_scoring.html" --profile pennaz --region us-west-2 --no-cli-pager
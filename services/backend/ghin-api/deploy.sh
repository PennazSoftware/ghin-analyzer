#!/bin/sh
# Variables
component=ghin-api
PROFILE=pennaz
REGION=us-west-2

# Verify Command Line Parameters
if [ -z "$1" ]
then
    echo You must specify an environment as the first parameter. Options are Dev or Prod
    exit 1
fi

envir=$1

if [ "$envir" = "Dev" ]
then
    envir="dev"
fi

if [ "$envir" = "dev" ]
then
    echo Deploying a development version of $component is not supported as it is production only.
    exit 1
fi

if [ "$envir" = "Prod" ]
then
    envir="prod"
fi

APP=$component-$envir
bash build.sh $1

# build
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap
zip deploy.zip ./bootstrap

echo Updating lambda: $APP

# Update the lambda

if [ "$(aws lambda get-function --profile $PROFILE --region $REGION --function-name=$APP)" ]; then
    aws lambda update-function-code --profile $PROFILE --no-cli-pager --region $REGION --function-name=$APP --zip-file=fileb://deploy.zip
else
    echo This Lambda must have been deployed via Terraform first
    echo Could not find Lambda: $APP in region: $REGION with profile: $PROFILE and environment: $envir
fi
rm deploy.zip
rm bootstrap

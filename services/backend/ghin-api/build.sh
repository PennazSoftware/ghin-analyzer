#!/bin/sh
component=ghin-api
echo Building $component...

# Verify Command Line Parameters
if [ -z "$1" ]
then
    echo You must specify an environment as the first parameter. Options are dev or prod
    exit 1
fi

envir=$1

if [ "$envir" = "Dev" ]
then
    envir="dev"
fi

if [ "$envir" = "dev" ]
then
    echo Building a development version of $component is not supported as it is production only.
    exit 1
fi

if [ "$envir" = "Prod" ]
then
    envir="prod"
fi

# Cleanup any previous zip file
rm ../../../build/bin/$component/$envir/bootstrap*

# Build Query Lamba into a zip file for deployment
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o ../../../build/bin/$component/$envir/bootstrap

rm lambda.zip
go build -o main
zip -X -r lambda.zip main
aws lambda update-function-code --function-name getCachedScore --zip-file fileb://lambda.zip --architectures "x86_64"

aws lambda update-function-configuration --function-name getCachedScore --handler main \
                                --timeout 300 --environment "Variables={GIT_PAT=$GIT_PAT}" --runtime go1.x
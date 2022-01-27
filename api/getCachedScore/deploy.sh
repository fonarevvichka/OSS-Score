rm lambda.zip
go build main.go
zip -X -r lambda.zip main
aws lambda update-function-code --function-name getCachedScore --zip-file fileb://lambda.zip
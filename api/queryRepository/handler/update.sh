# handler lambda
rm main lambda.zip
go build -o main
zip -X -r lambda.zip main mongo_cert.pem
aws lambda update-function-code --function-name queryScoreHandler --zip-file fileb://lambda.zip --architectures "x86_64" > /dev/null
sleep 3s
aws lambda update-function-configuration --function-name queryScoreHandler --handler main \
                                --timeout 300 --environment "Variables={GIT_PAT=$GIT_PAT_6, MONGO_URI=$MONGO_URI, SHELF_LIFE=$SHELF_LIFE, QUERY_QUEUE=$REPO_QUERY_QUEUE, DYNAMODB_TABLE=$DYNAMODB_TABLE}" \
                                --runtime go1.x > /dev/null
# handler lambda
rm main lambda.zip
go build -o main
cp ../../util/mongo_cert.pem ./mongo_cert.pem
zip -X -r lambda.zip main mongo_cert.pem
aws lambda update-function-code --function-name queryScoreHandler --zip-file fileb://lambda.zip --architectures "x86_64" > /dev/null
sleep 3s
aws lambda update-function-configuration --function-name queryScoreHandler --handler main \
                                --timeout 300 --environment "Variables={MONGO_URI=$MONGO_URI, SHELF_LIFE=$SHELF_LIFE, QUERY_QUEUE=$QUERY_QUEUE}" \
                                --runtime go1.x > /dev/null
rm mongo_cert.pem
rm main
rm lambda.zip

cp ../util/mongo_cert.pem ./mongo_cert.pem
cp -r ../util/scores ./scores

go build -o main
zip -X -r lambda.zip main mongo_cert.pem scores

aws lambda update-function-code --function-name getMetric --zip-file fileb://lambda.zip --architectures "x86_64" > /dev/null
sleep 3s
aws lambda update-function-configuration --function-name getMetric --handler main \
                                --timeout 300 --environment "Variables={MONGO_URI=$MONGO_URI, SHELF_LIFE=$SHELF_LIFE, DYNAMODB_TABLE=$DYNAMODB_TABLE}" \
                                --runtime go1.x > /dev/null

rm mongo_cert.pem
rm -r scores
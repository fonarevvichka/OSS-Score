rm main
rm lambda.zip
cp ../util/mongo_cert.pem ./mongo_cert.pem
go build -o main
zip -X -r lambda.zip main mongo_cert.pem
aws lambda update-function-code --function-name getCachedScore --zip-file fileb://lambda.zip --architectures "x86_64" > /dev/null

aws lambda update-function-configuration --function-name getCachedScore --handler main \
                                --timeout 300 --environment "Variables={MONGO_URI=$MONGO_URI, SHELF_LIFE=$SHELF_LIFE}" \
                                --runtime go1.x > /dev/null

rm mongo_cert.pem
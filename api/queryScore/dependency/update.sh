# processing lambda
rm main lambda.zip
go build -o main
cp ../../util/mongo_cert.pem ./mongo_cert.pem
cp -r ../../util/queries ./queries
cp -r ../../util/scores ./scores

zip -X -r lambda.zip main mongo_cert.pem queries scores
aws lambda update-function-code --function-name queryDependency --zip-file fileb://lambda.zip --architectures "x86_64" > /dev/null
sleep 3s
aws lambda update-function-configuration --function-name queryDependency --handler main \
                                --timeout 300 --environment "Variables={GIT_PAT_1=$GIT_PAT_1, GIT_PAT_2=$GIT_PAT_2, GIT_PAT_3=$GIT_PAT_3, GIT_PAT_4=$GIT_PAT_4, GIT_PAT_5=$GIT_PAT_5, MONGO_URI=$MONGO_URI, SHELF_LIFE=$SHELF_LIFE}" \
                                --runtime go1.x > /dev/null

rm mongo_cert.pem
rm -r scores
rm -r queries
.PHONY: build clean deploy

build:
	cd getScore && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/getScore getScore.go && cd ..
	cd getMetric && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/getMetric getMetric.go && cd ..
	cd queryScore/handler && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/queryScoreHandler queryScoreHandler.go && cd ../..
	cd queryScore/processing && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/queryScoreProcessing queryScoreProcessing.go && cd ../..
	cd queryScore/dependency && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/queryScoreDependency queryScoreDependency.go && cd ../..

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
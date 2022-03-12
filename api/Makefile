.PHONY: build clean deploy

build:
	cd getScore && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/getScore getScore.go && cd ..
	cd getMetric && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/getMetric getMetric.go && cd ..
	cd queryRepository/handler && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/queryRepoHandler queryRepoHandler.go && cd ../..
	cd queryRepository/processing && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/queryRepoProcessing queryRepoProcessing.go && cd ../..
	cd queryRepository/dependency && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../bin/queryDependency queryDependency.go && cd ../..

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

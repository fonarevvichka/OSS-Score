getScore:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getScore ./getScore/getScore.go
	
getMetric:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getMetric ./getMetric/getMetric.go

handler:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/queryRepoHandler queryRepository/handler/queryRepoHandler.go

processing:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/queryRepoProcessing queryRepository/processing/queryRepoProcessing.go

dependency:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/queryDependency queryRepository/dependency/queryDependency.go

build: getScore getMetric handler processing dependency

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

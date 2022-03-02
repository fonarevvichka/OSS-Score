.PHONY: build clean deploy

build:
	cd getScore && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/getScore getScore.go && cd ..

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
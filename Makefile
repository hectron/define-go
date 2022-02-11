.PHONY: clean build

build: clean
	go build -o bin/go-define define.go

clean:
	rm -rf ./bin

test:
	go test ./...

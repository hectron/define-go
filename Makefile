.PHONY: clean build

build: clean
	go build -o bin/go-define .

clean:
	rm -rf ./bin

test:
	go test ./...

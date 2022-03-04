.PHONY: clean build

build: clean update-remote
	go build -o bin/go-define -ldflags="-X 'main.Version=$$(git describe --tags --abbrev=0)'"

update-remote:
	git remote update

clean:
	rm -rf ./bin

test:
	go test ./...

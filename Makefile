
build:
	go build -a -tags netgo -ldflags '-w -extldflags "-static"'

echo:
	@echo "$(GOPATH)"

fmt:
	find . -path ./.go -prune -o -name "*.go" -exec gofmt -w {} \;
	find . -path ./.go -prune -o -name "*.i2pkeys" -exec rm {} \;

install:
	install -m755 i2pkeys /usr/local/bin/i2pkeys

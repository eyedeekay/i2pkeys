
build:
	go build -a -tags netgo -ldflags '-w -extldflags "-static"'

install:
	install -m755 i2pkeys /usr/local/bin/i2pkeys

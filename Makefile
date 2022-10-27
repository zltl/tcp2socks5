.PHONY: tcp2socks5 clean

VERSION="1.0.0"

tcp2socks5:
	go build -ldflags '-X main.VERSION=${VERSION}' ./cmd/tcp2socks5.go

clean:
	rm -rf tcp2socks5

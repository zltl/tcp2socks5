.PHONY: tcp2socks5 clean

tcp2socks5:
	go build ./cmd/tcp2socks5.go

clean:
	rm -rf tcp2socks5

# tcp2socks5

Tunnel TCP port to socks5 proxy.

```sh
make

./tcp2socks5 --local 0.0.0.0:4444 --socks5 127.0.0.1:1080 --target www.google.com:80

# curl curl -H "Host: www.google.com" localhost:4444
```

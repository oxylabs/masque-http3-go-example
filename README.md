# MASQUE HTTP/3 client example

This is an example in Go on how to proxy HTTP/3 requests through a MASQUE CONNECT-UDP capable proxy.


To build and run from source, have Go installed (https://go.dev/doc/install), then:

```shell
git clone https://github.com/oxylabs/masque-http3-go-example
cd masque-http3-go-example
go build -o masque-h3 .
./masque-h3 -u USERNAME -p PASSWORD [-t ipv4.oxylabs.io]
```

Usage (type `./masque-h3 --usage` after build):

```shell
  -insecure-origin
        skip TLS verification for the origin (default: false)
  -insecure-proxy
        skip TLS certificate verification for the proxy (default: false)
  -p string
        proxy password
  -proxy string
        proxy host (default "masque.oxylabs.io:50000")
  -t string
        target host (default "cloudflare-quic.com")
  -timeout duration
        timeout (default 30s)
  -tp int
        target port (default 443)
  -u string
        proxy user
```
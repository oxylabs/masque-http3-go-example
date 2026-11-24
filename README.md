# MASQUE HTTP/3 client example

<a href="https://oxylabs.io"><img src="https://img.shields.io/badge/🚀_Get_Started-FF6B35?style=for-the-badge" alt="Get started"></a>
<a href="https://developers.oxylabs.io"><img src="https://img.shields.io/badge/📚_Documentation-4A90E2?style=for-the-badge" alt="Documentation"></a>
<a href="https://dashboard.oxylabs.io"><img src="https://img.shields.io/badge/📊_Dashboard-2EA44F?style=for-the-badge" alt="Dashboard"></a>
<a href="https://discord.gg/Pds3gBmKMH"><img src="https://img.shields.io/badge/💬_Discord-5865F2?style=for-the-badge" alt="Discord"></a>
<a href="https://www.youtube.com/@oxylabs"><img src="https://img.shields.io/badge/▶️_YouTube-FF0000?style=for-the-badge" alt="YouTube"></a>

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

*Google is a restricted target to use with UDP connection.*

[Residential Proxies](https://github.com/oxylabs/residential-proxies)

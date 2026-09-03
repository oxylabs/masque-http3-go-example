package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/oxylabs/masque-http3-go-example/client"
)

var (
	defaultProxyHost  = "masque.oxylabs.io:50000" // MASQUE proxy (HTTP/3)
	defaultTargetHost = "cloudflare-quic.com"     // Destination target reached
	defaultTargetPort = 443
	defaultTimeout    = 30 * time.Second
)

func main() {
	proxyUser := flag.String("u", "", "proxy user")
	proxyPass := flag.String("p", "", "proxy password")
	proxyHost := flag.String("proxy", defaultProxyHost, "proxy host")
	targetHost := flag.String("t", defaultTargetHost, "target host")
	targetPort := flag.Int("tp", defaultTargetPort, "target port")
	timeout := flag.Duration("timeout", defaultTimeout, "timeout")
	insecureProxy := flag.Bool("insecure-proxy", false, "skip TLS certificate verification for the proxy (default: false)")
	insecureOrigin := flag.Bool("insecure-origin", false, "skip TLS verification for the origin (default: false)")

	flag.Parse()

	start := time.Now()
	resp, err := client.HTTP3OverMasque(*proxyUser, *proxyPass, *proxyHost, *targetHost, *targetPort, *insecureProxy, *insecureOrigin, *timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Request took: ", time.Since(start).String(), "\n")
	fmt.Println("Status from target: ", resp.Status, "\n")
	fmt.Println("Headers: ")
	for k, v := range resp.Header {
		fmt.Println(k, ": ", v[0])
	}
	fmt.Println()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 200))
	fmt.Println("Body (first 200 chars): \n", string(body))
}

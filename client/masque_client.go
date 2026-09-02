package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/quicvarint"
)

func HTTP3OverMasque(proxyUsername, proxyPassword, proxyHost, targetHost string, targetPort int, insecureProxy, insecureOrigin bool, timeout time.Duration) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		time.Sleep(timeout)
		cancel()
	}()

	proxyTLS := &tls.Config{
		NextProtos:         []string{http3.NextProtoH3},
		InsecureSkipVerify: insecureProxy,
	}
	proxyQuicConf := &quic.Config{
		EnableDatagrams:         true,
		DisablePathMTUDiscovery: true,
		InitialPacketSize:       1392, // Must be about 40 bytes greater than inner QUIC
	}

	proxyConn, err := quic.DialAddr(ctx, proxyHost, proxyTLS, proxyQuicConf)
	if err != nil {
		return nil, fmt.Errorf("dial proxy quic: %w", err)
	}

	h3ProxyConn := (&http3.Transport{EnableDatagrams: true}).NewClientConn(proxyConn)

	reqURL := &url.URL{
		Scheme: "https",
		Host:   proxyHost,
		Path:   fmt.Sprintf("/.well-known/masque/udp/%s/%d/", targetHost, targetPort),
	}
	req := (&http.Request{
		Method: http.MethodConnect,
		URL:    reqURL,
		Host:   proxyHost,
		Header: make(http.Header),
	}).WithContext(ctx)
	if proxyUsername != "" && proxyPassword != "" {
		req.Header.Set("Proxy-Authorization", basicAuth(proxyUsername, proxyPassword))
	}
	req.Proto = "connect-udp" // extended CONNECT :protocol

	rstr, err := h3ProxyConn.OpenRequestStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("open request stream: %w", err)
	}
	if err := rstr.SendRequestHeader(req); err != nil {
		return nil, fmt.Errorf("send CONNECT-UDP header: %w", err)
	}

	resp, err := rstr.ReadResponse()
	if err != nil {
		return nil, fmt.Errorf("read CONNECT-UDP response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("proxy rejected CONNECT-UDP: %s", resp.Status)
	}
	log.Printf("MASQUE tunnel established (status %d)", resp.StatusCode)

	pc := &masqueUDPConn{
		stream: *rstr,
		target: &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: targetPort},
	}

	innerTLS := &tls.Config{
		ServerName:         targetHost,
		NextProtos:         []string{http3.NextProtoH3},
		InsecureSkipVerify: insecureOrigin,
	}
	innerTr := &http3.Transport{
		EnableDatagrams: false,
		Dial: func(ctx context.Context, _ string, tlsConf *tls.Config, qConf *quic.Config) (*quic.Conn, error) {
			qConf.InitialPacketSize = 1352 // This must be lower than the outer/proxy QUIC, 1352 ir pretty much the
			// maximum with the outer at 1392 as there needs to be room for QUIC overhead
			qConf.DisablePathMTUDiscovery = true // PMTUD must be disabled for the inner QUIC
			// as it's not aware of the outer/proxy QUIC and will mismeasure the MTU
			return quic.Dial(ctx, pc, pc.target, tlsConf, qConf)
		},
		TLSClientConfig: innerTLS,
	}
	go func() {
		time.Sleep(timeout)
		innerTr.Close()
	}()

	client := &http.Client{Transport: innerTr}

	targetURL := fmt.Sprintf("https://%s:%d/", targetHost, targetPort)
	greq, _ := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	greq.Header.Set("User-Agent", "masque-h3-example/1.0")

	gresp, err := client.Do(greq)
	if err != nil {
		return gresp, fmt.Errorf("target request: %w", err)
	}

	return gresp, nil
}

type masqueUDPConn struct {
	stream http3.RequestStream
	target net.Addr
}

func (c *masqueUDPConn) ReadFrom(p []byte) (int, net.Addr, error) {
	for {
		dgram, err := c.stream.ReceiveDatagram(context.Background())
		if err != nil {
			return 0, nil, err
		}
		r := bytes.NewReader(dgram)
		ctxID, err := quicvarint.Read(r)
		if err != nil {
			continue
		}
		if ctxID != 0 {
			// Non-zero context IDs are for extensions we don't handle here.
			continue
		}
		n, _ := r.Read(p)
		return n, c.target, nil
	}
}

func (c *masqueUDPConn) WriteTo(p []byte, _ net.Addr) (int, error) {
	buf := make([]byte, 0, len(p)+8)
	buf = quicvarint.Append(buf, 0) // context ID 0 => raw UDP payload
	buf = append(buf, p...)
	if err := c.stream.SendDatagram(buf); err != nil {
		return 0, err
	}
	return len(p), nil
}

func basicAuth(user, pass string) string {
	creds := user + ":" + pass
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}

func (c *masqueUDPConn) Close() error                       { c.stream.Close(); return nil }
func (c *masqueUDPConn) LocalAddr() net.Addr                { return &net.UDPAddr{} }
func (c *masqueUDPConn) SetDeadline(t time.Time) error      { return nil }
func (c *masqueUDPConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *masqueUDPConn) SetWriteDeadline(t time.Time) error { return nil }

package sankaku

import (
	"encoding/base64"
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/proxy"
)

var (
	client = &fasthttp.Client{
		Name:                          "Sankaku",
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}
)

func UseProxy(addr string) error {
	URL, err := url.Parse(addr)
	if err != nil {
		return err
	}

	switch URL.Scheme {
	case "http":
		client.Dial = dialerHTTPProxy(URL)
	case "socks5", "socks5h":
		client.Dial = dialerSOCKS5Proxy(URL)
	default:
		return errors.New("only http and socks5 proxy are supported")
	}

	return nil
}

func dialerHTTPProxy(p *url.URL) fasthttp.DialFunc {
	authHeader := "\r\n"
	if p.User != nil {
		auth := base64.StdEncoding.EncodeToString([]byte(p.User.String()))
		authHeader = "Proxy-Authorization: Basic " + auth + "\r\n\r\n"
	}
	host := p.Host
	return func(addr string) (net.Conn, error) {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			return nil, err
		}
		_, err = conn.Write([]byte("CONNECT " + addr + " HTTP/1.1\r\nHost: " + addr + "\r\n" + authHeader))
		if err != nil {
			conn.Close()
			return nil, err
		}
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if !strings.HasPrefix(string(buf[:n]), "HTTP/1.1 200 Connection established") {
			conn.Close()
			return nil, errors.New("proxy error: \n" + string(buf[:n]))
		}
		return conn, nil
	}
}

func dialerSOCKS5Proxy(p *url.URL) fasthttp.DialFunc {
	var auth *proxy.Auth
	if p.User != nil {
		auth = &proxy.Auth{User: p.User.Username()}
		if pass, ok := p.User.Password(); ok {
			auth.Password = pass
		}
	}
	dialer, err := proxy.SOCKS5("tcp", p.Host, auth, proxy.Direct)
	if err != nil {
		return nil
	}
	return func(addr string) (net.Conn, error) {
		return dialer.Dial("tcp", addr)
	}
}

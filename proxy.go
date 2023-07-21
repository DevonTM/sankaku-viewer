package sankaku

import (
	"errors"
	"net/http"
	"net/url"
)

var (
	transport = &http.Transport{}
	client    = &http.Client{
		Transport: transport,
	}
)

func UseProxy(addr string) error {
	URL, err := url.Parse(addr)
	if err != nil {
		return err
	}
	if URL.Scheme != "http" && URL.Scheme != "https" && URL.Scheme != "socks5" {
		return errors.New("only http, https and socks5 proxy are supported")
	}
	transport.Proxy = http.ProxyURL(URL)
	return nil
}

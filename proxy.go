package sankaku

import (
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
		return url.InvalidHostError(URL.Scheme)
	}
	transport.Proxy = http.ProxyURL(URL)
	return nil
}

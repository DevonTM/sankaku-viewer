package sankaku

import (
	"errors"
	"net/url"

	"github.com/valyala/fasthttp"
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
	if URL.Scheme != "http" && URL.Scheme != "https" && URL.Scheme != "socks5" {
		return errors.New("only http, https and socks5 proxy are supported")
	}
	return nil
}

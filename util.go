package sankaku

import (
	"errors"
	"net/url"
	"path"

	"github.com/patrickmn/go-cache"
	"github.com/valyala/fasthttp"
)

func getBaseURL(ctx *fasthttp.RequestCtx) string {
	scheme := string(ctx.URI().Scheme())
	host := string(ctx.URI().Host())
	if ctx.Request.Header.Peek("X-Forwarded-Proto") != nil {
		scheme = string(ctx.Request.Header.Peek("X-Forwarded-Proto"))
	}
	if ctx.Request.Header.Peek("X-Forwarded-Host") != nil {
		host = string(ctx.Request.Header.Peek("X-Forwarded-Host"))
	}
	return scheme + "://" + host + "/"
}

func getData(id string) (data *PostData, err error) {
	d, ok := c.Get(id)
	if ok {
		data, ok = d.(*PostData)
		if !ok {
			err = errors.New("unable to get data from cache")
		}
	} else {
		data, err = GetPost(id)
		if err == nil {
			if data.URL != "" {
				c.Set(id, data, cache.DefaultExpiration)
			} else {
				err = errors.New("login required")
			}
		}
	}
	return
}

func getName(tags []Tag) string {
	var names [2]string
	for _, tag := range tags {
		if tag.Type == 3 {
			names[1] = tag.Name
		} else if tag.Type == 4 {
			names[0] = tag.Name
		}
		if names[0] != "" && names[1] != "" {
			break
		}
	}
	switch {
	case names[0] != "" && names[1] != "":
		return names[0] + " - " + names[1]
	case names[0] == "":
		return names[1]
	case names[1] == "":
		return names[0]
	}
	return "Sankaku Content"
}

func getID(rawURL string) (string, error) {
	URL, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("cannot parse URL")
	}
	if path.Dir(URL.Path) != "/post/show" {
		return "", errors.New("invalid URL")
	}
	return path.Base(URL.Path), nil
}

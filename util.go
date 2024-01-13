package sankaku

import (
	"errors"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/valyala/fasthttp"
)

var (
	cacheTTL time.Duration = 10 * time.Minute
	c                      = cache.New(cache.DefaultExpiration, 1*time.Hour)
)

func SetCacheDuration(seconds int) {
	if seconds > 3600 {
		seconds = 3600
	}
	cacheTTL = time.Duration(seconds) * time.Second
}

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
			err = errors.New("Unable to get data from cache")
		}
	} else {
		data, err = GetPost(id)
		if err == nil {
			if data.URL != "" {
				c.Set(id, data, cacheTTL)
			} else {
				err = errors.New("Login required")
			}
		}
	}
	return
}

func getName(tags []Tag) string {
	var series, names string
	var characters []string
	for _, tag := range tags {
		if tag.Type == 3 {
			series = tag.Name
		} else if tag.Type == 4 {
			characters = append(characters, tag.Name)
		}
	}
	if len(characters) > 0 {
		names = strings.Join(characters, " - ")
	}
	switch {
	case series != "" && names != "":
		return names + " - " + series
	case series != "":
		return series
	case names != "":
		return names
	}
	return "Sankaku Content"
}

func getID(rawURL string) (string, error) {
	URL, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("Cannot parse URL: " + err.Error())
	}
	paths := path.Clean(URL.Path)
	_, paths, ok := strings.Cut(paths, "/posts/")
	if !ok {
		return "", errors.New("Invalid URL")
	}
	paths = strings.Split(paths, "/")[0]
	return paths, nil
}

func isValidExt(fileName string) bool {
	extensions := []string{".css", ".ico", ".js", ".png"}
	for _, ext := range extensions {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}

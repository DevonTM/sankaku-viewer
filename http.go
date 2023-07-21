package sankaku

import (
	"errors"
	"html/template"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/valyala/fasthttp"
)

type PageData struct {
	Title  string
	Loc    string
	URL    string
	Poster string
	Type   string
	Format string
	Ori    string
	Error  string
	Width  int
	Height int
	Size   int
	ID     int
}

var c = cache.New(10*time.Minute, 1*time.Hour)

func ListenAndServe(addr string) error {
	ln, err := listen(addr)
	if err != nil {
		return err
	}
	server := &fasthttp.Server{
		Name:            "Sankaku",
		Handler:         fasthttp.CompressHandler(handler),
		GetOnly:         true,
		CloseOnShutdown: true,
	}
	err = server.Serve(ln)
	return err
}

func listen(addr string) (net.Listener, error) {
	network := "tcp"
	if strings.HasPrefix(addr, "unix:") {
		addr = strings.TrimPrefix(addr, "unix:")
		network = "unix"
		if err := os.RemoveAll(addr); err != nil {
			return nil, err
		}
	}
	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	if network == "unix" {
		if err := os.Chmod(addr, 0o666); err != nil {
			return nil, err
		}
	}
	return ln, nil
}

func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/favicon.ico":
		ctx.SendFile("./static/favicon.ico")
		return
	case "/logo.png":
		ctx.SendFile("./static/logo.png")
		return
	}

	if ctx.QueryArgs().Has("id") && ctx.QueryArgs().Has("type") {
		id := string(ctx.QueryArgs().Peek("id"))
		data, err := getData(id)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			_, _ = ctx.WriteString(err.Error())
			return
		}
		typ := string(ctx.QueryArgs().Peek("type"))
		switch typ {
		case "url":
			ctx.Redirect(data.URL, fasthttp.StatusFound)
		case "surl":
			ctx.Redirect(data.Sample, fasthttp.StatusFound)
		case "purl":
			ctx.Redirect(data.Preview, fasthttp.StatusFound)
		}
		return
	}

	scheme := string(ctx.URI().Scheme())
	if ctx.Request.Header.Peek("X-Forwarded-Proto") != nil {
		scheme = string(ctx.Request.Header.Peek("X-Forwarded-Proto"))
	}
	loc := scheme + "://" + string(ctx.URI().Host()) + "/"

	if ctx.QueryArgs().Has("url") {
		URL := string(ctx.QueryArgs().Peek("url"))
		id, err := getID(URL)
		if err != nil {
			render(ctx, PageData{
				Loc:   loc,
				Error: err.Error(),
			})
			return
		}
		data, err := getData(id)
		if err != nil {
			render(ctx, PageData{
				Loc:   loc,
				Error: err.Error(),
			})
			return
		}
		var typ string
		if strings.Contains(data.Content, "/") {
			typ = strings.Split(data.Content, "/")[0]
		} else {
			typ = data.Content
		}
		ori := strings.Split(URL, "?")[0]
		render(ctx, PageData{
			Loc:    loc,
			Type:   typ,
			Ori:    ori,
			Title:  data.Name,
			URL:    data.URL,
			Poster: data.Preview,
			Format: data.Content,
			Width:  data.Width,
			Height: data.Height,
			Size:   data.Size,
			ID:     data.ID,
		})
	} else {
		render(ctx, PageData{
			Loc: loc,
		})
	}
}

func render(ctx *fasthttp.RequestCtx, data interface{}) {
	page := "./static/index.html"
	t, err := template.ParseFiles(page)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = ctx.WriteString("Failed to parse template")
		return
	}
	ctx.SetContentType("text/html; charset=utf-8")
	if err = t.Execute(ctx, data); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = ctx.WriteString("Failed to execute template")
		return
	}
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

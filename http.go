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
	SURL   string
	PURL   string
	Type   string
	Format string
	Ori    string
	Error  string
	Width  int
	Height int
	Size   int
}

var c = cache.New(10*time.Minute, 1*time.Hour)

func ListenAndServe(addr string) error {
	ln, err := Listen(addr)
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

func Listen(addr string) (net.Listener, error) {
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
	if string(ctx.Path()) == "/favicon.ico" {
		ctx.SendFile("./static/favicon.ico")
		return
	}
	if string(ctx.Path()) == "/logo.png" {
		ctx.SendFile("./static/logo.png")
		return
	}

	loc := string(ctx.URI().Scheme()) + "://" + string(ctx.URI().Host()) + "/"
	typ := "website"
	if ctx.QueryArgs().Has("url") {
		URL := string(ctx.QueryArgs().Peek("url"))
		id, err := getID(URL)
		if err != nil {
			render(ctx, PageData{
				Loc:   loc,
				Type:  typ,
				Error: err.Error(),
			})
			return
		}
		var data *PostData
		d, ok := c.Get(id)
		if !ok {
			data, err = getData(id)
			if err != nil {
				render(ctx, PageData{
					Loc:   loc,
					Type:  typ,
					Error: err.Error(),
				})
				return
			}
			c.Set(id, data, cache.DefaultExpiration)
		} else {
			data = d.(*PostData)
		}
		if data.URL == "" {
			render(ctx, PageData{
				Loc:   loc,
				Type:  typ,
				Error: "login required",
			})
			return
		}
		if strings.Contains(data.Content, "/") {
			typ = strings.Split(data.Content, "/")[0]
		} else {
			typ = data.Content
		}
		ori := strings.Split(URL, "?")[0]
		loc += "?url=" + url.QueryEscape(ori)
		title := getName(data)
		render(ctx, PageData{
			Title:  title,
			Loc:    loc,
			Type:   typ,
			Ori:    ori,
			URL:    data.URL,
			SURL:   data.Sample,
			PURL:   data.Preview,
			Format: data.Content,
			Width:  data.Width,
			Height: data.Height,
			Size:   data.Size,
		})
	} else {
		render(ctx, PageData{
			Loc:  loc,
			Type: typ,
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

func getName(data *PostData) string {
	var names [2]string
	for _, tag := range data.Tags {
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

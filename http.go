package sankaku

import (
	"html/template"
	"net"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)

type PageData struct {
	Title  string
	Loc    string
	URL    string
	Poster string
	Type   string
	Format string
	Ext    string
	Ori    string
	Error  string
	Width  int
	Height int
	Size   int
	ID     int
}

var serveFile = fileHandler("static")

func ListenAndServe(addr string) error {
	ln, err := listen(addr)
	if err != nil {
		return err
	}
	server := &fasthttp.Server{
		Name:            "Sankaku",
		Handler:         fasthttp.CompressHandlerBrotliLevel(requestHandler, 3, 9),
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

func requestHandler(ctx *fasthttp.RequestCtx) {
	switch path := string(ctx.Path()); {
	case path == "/get":
		handleGet(ctx)
	case strings.HasPrefix(path, "/redir"):
		handleRedir(ctx)
	case isValidExt(path):
		serveFile(ctx)
	default:
		loc := getBaseURL(ctx)
		render(ctx, PageData{Loc: loc})
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

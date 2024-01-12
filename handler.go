package sankaku

import (
	"strings"

	"github.com/valyala/fasthttp"
)

var CacheCompressed bool

func fileHandler(root string) fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		Root:            root,
		CompressRoot:    root + "/.cache",
		Compress:        CacheCompressed,
		CompressBrotli:  true,
		AcceptByteRange: true,
	}
	return fs.NewRequestHandler()
}

func handleRedir(ctx *fasthttp.RequestCtx) {
	id := string(ctx.QueryArgs().Peek("id"))
	if id == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = ctx.WriteString("Missing Post ID")
		return
	}
	typ := string(ctx.QueryArgs().Peek("type"))
	if typ == "" {
		ctx.Redirect(APIPosts+"?tags=id:"+id, fasthttp.StatusPermanentRedirect)
		return
	}
	data, err := getData(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	switch typ {
	case "url":
		ctx.Redirect(data.URL, fasthttp.StatusFound)
	case "surl":
		ctx.Redirect(data.Sample, fasthttp.StatusFound)
	case "purl":
		ctx.Redirect(data.Preview, fasthttp.StatusFound)
	default:
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = ctx.WriteString("Invalid type")
	}
}

func handleGet(ctx *fasthttp.RequestCtx) {
	loc := getBaseURL(ctx)
	var id, ori string
	if ctx.QueryArgs().Has("id") {
		id = string(ctx.QueryArgs().Peek("id"))
		ori = "https://sankaku.app/posts/" + id
	} else if ctx.QueryArgs().Has("url") {
		URL := string(ctx.QueryArgs().Peek("url"))
		var err error
		id, err = getID(URL)
		if err != nil {
			render(ctx, PageData{
				Loc:   loc,
				Error: err.Error(),
			})
			return
		}
		ori = strings.Split(URL, "?")[0]
	} else {
		render(ctx, PageData{
			Loc:   loc,
			Error: "Missing Post ID or URL",
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
	var typ, ext string
	if strings.Contains(data.Content, "/") {
		typ = strings.Split(data.Content, "/")[0]
		ext = "." + strings.Split(data.Content, "/")[1]
	} else {
		typ = data.Content
	}
	render(ctx, PageData{
		Loc:    loc,
		Type:   typ,
		Ext:    ext,
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
}

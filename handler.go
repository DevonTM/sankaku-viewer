package sankaku

import (
	"strings"

	"github.com/valyala/fasthttp"
)

func handleRedir(ctx *fasthttp.RequestCtx) {
	id := string(ctx.QueryArgs().Peek("id"))
	if id == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		_, _ = ctx.WriteString("Missing id")
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
}

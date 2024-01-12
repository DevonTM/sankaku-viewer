package sankaku

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/valyala/fasthttp"
)

const APIPosts = "https://capi-v2.sankakucomplex.com/posts"

type Tag struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type PostData struct {
	Name    string `json:"-"`
	Sample  string `json:"sample_url"`
	Preview string `json:"preview_url"`
	URL     string `json:"file_url"`
	Content string `json:"file_type"`
	Tags    []Tag  `json:"tags"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Size    int    `json:"file_size"`
	ID      int    `json:"id"`
}

func GetPost(id string) (*PostData, error) {
	payload := url.Values{}
	payload.Set("tags", "id:"+id)
	URL := APIPosts + "?" + payload.Encode()

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(URL)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("SCChannelApp/4.0")
	req.Header.Set("Accept", "application/vnd.sankaku.api+json;v=2")
	req.Header.Set("Origin", "https://sankaku.app")
	req.Header.SetReferer("https://sankaku.app/posts/" + id)
	if Token != "" {
		req.Header.Set("Authorization", Token)
	}

	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, errors.New("Failed to send request: " + err.Error())
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New("Failed to get data, HTTP status: " + string(resp.Header.StatusMessage()))
	}

	var data []PostData
	if err = json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, errors.New("Failed to unmarshal response body: " + err.Error())
	}

	data[0].Name = getName(data[0].Tags)

	return &data[0], nil
}

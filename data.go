package sankaku

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const APIPosts = "https://capi-v2.sankakucomplex.com/posts"

type PostData struct {
	Sample  string `json:"sample_url"`
	Preview string `json:"preview_url"`
	URL     string `json:"file_url"`
	Content string `json:"file_type"`
	Tags    []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"tags"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Size   int `json:"file_size"`
}

func getData(id string) (*PostData, error) {
	payload := url.Values{}
	payload.Set("lang", "en")
	payload.Set("page", "1")
	payload.Set("limit", "1")
	if len(id) == 32 {
		payload.Set("tags", "md5:"+id)
	} else {
		payload.Set("tags", "id_range:"+id)
	}

	query := payload.Encode()
	URL := fmt.Sprintf("%s?%s", APIPosts, query)

	req, err := http.NewRequest(http.MethodGet, URL, http.NoBody)
	if err != nil {
		return nil, errors.New("failed to create request")
	}

	req.Header.Set("Accept", "application/vnd.sankaku.api+json;v=2")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://sankaku.app")
	req.Header.Set("Referer", "https://sankaku.app/")
	req.Header.Set("User-Agent", "SCChannelApp/4.0")
	if Token != "" {
		req.Header.Set("Authorization", Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to get data, HTTP status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body")
	}

	var data []PostData
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, errors.New("Failed to parse response body: " + err.Error())
	}

	return &data[0], nil
}

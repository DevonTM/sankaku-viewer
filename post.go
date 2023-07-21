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

type Tag struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	ID   int    `json:"id"`
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

	data[0].Name = getName(data[0].Tags)

	return &data[0], nil
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
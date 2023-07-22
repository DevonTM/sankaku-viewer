package sankaku

import (
	"encoding/json"
	"errors"

	"github.com/valyala/fasthttp"
)

const APIAuth = "https://capi-v2.sankakucomplex.com/auth/token"

var Token string

func Login(username, password string) error {
	payload := map[string]interface{}{
		"login":    username,
		"password": password,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.New("Failed to marshal login payload: " + err.Error())
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(APIAuth)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetUserAgent("SCChannelApp/4.0")
	req.Header.SetContentType("application/json")
	req.Header.Set("Accept", "application/vnd.sankaku.api+json;v=2")
	req.SetBodyRaw(jsonPayload)

	resp := fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return errors.New("Failed to send login request: " + err.Error())
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		if resp.StatusCode() == fasthttp.StatusUnauthorized || resp.StatusCode() == fasthttp.StatusForbidden {
			return errors.New("Invalid username or password")
		}
		return errors.New("Failed to login, HTTP status: " + string(resp.Header.StatusMessage()))
	}

	var response map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return errors.New("Failed to unmarshal login response: " + err.Error())
	}

	Token = "Bearer " + response["access_token"].(string)

	return nil
}

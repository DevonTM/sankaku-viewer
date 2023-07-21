package sankaku

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
		return errors.New("failed to marshal login payload")
	}

	req, err := http.NewRequest(http.MethodPost, APIAuth, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errors.New("failed to create login request")
	}

	req.Header.Set("Accept", "application/vnd.sankaku.api+json;v=2")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SCChannelApp/4.0")

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("failed to send login request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return errors.New("invalid username or password")
		}
		return fmt.Errorf("login request failed with status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return errors.New("failed to decode login token")
	}

	Token = fmt.Sprintf("Bearer %v", response["access_token"])
	return nil
}

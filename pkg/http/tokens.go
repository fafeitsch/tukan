package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api"
	"net/http"
	"time"
)

type PhoneClient interface {
	FetchToken()
}

func fetchToken(ip string, login string, password string) (*string, error) {
	url := fmt.Sprintf("http://{ip}/Login")
	credentials := api.Credentials{
		Login:    login,
		Password: password,
	}
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	payload, _ := json.Marshal(credentials)
	reader := bytes.NewBuffer(payload)
	resp, err := client.Post(url, "application/json", reader)
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	tokenResp := api.TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal received token: %v", err)
	}
	return &tokenResp.Token, nil
}

func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("authentication error, status code: %d, message: \"%s\"", resp.StatusCode, resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d with message \"%s\"", resp.StatusCode, resp.Status)
	}
	return nil
}

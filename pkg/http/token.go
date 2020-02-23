package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api"
	"net/http"
)

type tokener interface {
	fetchToken(ip string) (*string, error)
	logout(ip string, token string) error
}

type tokenerImpl struct {
	login    string
	password string
	port     int
	client   *http.Client
}

func (t *tokenerImpl) fetchToken(ip string) (*string, error) {
	url := fmt.Sprintf("http://%s:%d/Login", ip, t.port)
	credentials := api.Credentials{
		Login:    t.login,
		Password: t.password,
	}
	payload, _ := json.Marshal(credentials)
	reader := bytes.NewBuffer(payload)
	resp, err := t.client.Post(url, "application/json", reader)
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tokenResp := api.TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal token from %s: %v", ip, err)
	}
	return &tokenResp.Token, nil
}

func (t *tokenerImpl) logout(ip string, token string) error {
	url := fmt.Sprintf("http://%s:%d/Logout", ip, t.port)
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+token)
	resp, err := t.client.Do(request)
	err = checkResponse(resp, err)
	return err
}

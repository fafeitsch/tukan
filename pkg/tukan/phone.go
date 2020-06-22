package tukan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// A phone represents a http client that talks to exactly on
// physical telephone. A phone needs to be created with the method following method:
//   Connect(client,address,username,password)
// It is strongly recommended to defer calling the method Phone#Logout() because
// most IP620/630 only allow one active token at a time.
type Phone struct {
	client  *http.Client
	token   string
	address string
}

// Tries to connect to the phone at the given address and obtain an bearer token
// with the given credentials. If this is successful, a Phone is returned, otherwise
// this method returns an error.
func Connect(client *http.Client, address, username, password string) (*Phone, error) {
	url := fmt.Sprintf("%s/Login", address)
	credentials := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    username,
		Password: password,
	}
	payload, _ := json.Marshal(credentials)
	reader := bytes.NewBuffer(payload)
	resp, err := client.Post(url, "application/json", reader)
	err = checkResponse(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tokenResp := struct {
		Token string `json:"token"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal token from %s: %v", address, err)
	}
	return &Phone{
		client:  client,
		token:   tokenResp.Token,
		address: address,
	}, nil
}

func (p *Phone) Host() string {
	return p.address
}

func (p *Phone) Token() string {
	return p.token
}

// Sends a logout request to the phone. If the request passes without error
// then the token of the phone is reset. Further usage of the phone struct
// will most likely not work. If and error is returned, then the token stored
// in this telephone may or may not be used again, depending on the error.
func (p *Phone) Logout() error {
	url := fmt.Sprintf("%s/Logout", p.address)
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+p.token)
	resp, err := p.client.Do(request)
	err = checkResponse(resp, err)
	if err == nil {
		p.token = ""
	}
	return err
}

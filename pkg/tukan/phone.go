package tukan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/up"
	"net/http"
)

type Connector struct {
	Client   *http.Client
	UserName string
	Password string
}

// Tries to log in to a specific telephone identified by its address.
// On success, returns a phone Client, otherwise, an error is returned.
func (c *Connector) SingleConnect(address string) (*Phone, error) {
	url := fmt.Sprintf("%s/Login", address)
	credentials := up.Credentials{
		Login:    c.UserName,
		Password: c.Password,
	}
	payload, _ := json.Marshal(credentials)
	reader := bytes.NewBuffer(payload)
	resp, err := c.Client.Post(url, "application/json", reader)
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
		client:  c.Client,
		token:   tokenResp.Token,
		address: address,
	}, nil
}

// Result type for a multiple connect request. It carries either a phone or an error.
// The intention behind this struct is to simulate a channel having two types.
type ConnectResult struct {
	Phone *Phone
	Error error
}

// Connects to all phones given by the addresses parameter parallely and
// returns the results in the channel. Depending on how fast the real telephones answer
// to the login, the order in which the telephones are put in the channel differs
// from the order defined in the addresses parameter.
//
// This method returns immediately; the results channel is closed once all telephones
// have been contacted.
func (c *Connector) MultipleConnect(results chan<- ConnectResult, addresses ...string) {
	done := make(chan bool)
	for index, address := range addresses {
		go func(index int, address string) {
			phone, err := c.SingleConnect(address)
			if err != nil {
				err = fmt.Errorf("could not connect to address %d \"%s\": %v", index, address, err)
			}
			results <- ConnectResult{
				Phone: phone,
				Error: err,
			}
			done <- true
		}(index, address)
	}
	go func() {
		finished := 0
		for finished < len(addresses) {
			<-done
			finished = finished + 1
		}
		close(results)
	}()
}

// A phone represents a http Client that talks to exactly on
// physical telephone. A phone needs to be created with a Connector (see example).
// It is strongly recommended to defer calling the method Phone#Logout() because
// most IP620/630 only allow one active token at a time.
type Phone struct {
	client  *http.Client
	token   string
	address string
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

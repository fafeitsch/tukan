package tukan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/up"
	"net"
	"net/http"
	"sync"
)

// A connector is used to obtain a login token from a telephone and
// perform REST actions on the telephone. The connector can be used to either
// connect to only one telephone or to a bunch of telephones at the same time.
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

// Creates number addresses to connect to phones. This first address is given by startIp, which
// is always inclusive. If the startIp is not a valid IP, then an empty slice is returned.
// The resulting slice can be used in the MultipleConnect function of the Connector.
//
// Attention: This method does not contain any IP sub net logic, it just increments the ip addresses
// "stupidly": IP addresses like 10.20.30.255 or 10.20.255.0 may occur, and depending on the subnet
// these are valid host IP addresses or not. The method never returns syntactically wrong IP addresses.
func CreateAddresses(protocol, startIp string, port, number int) []string {
	incrementIp := func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}
	currentIp := net.ParseIP(startIp)
	if currentIp == nil {
		return []string{}
	}
	result := make([]string, 0, number)
	for i := 0; i < number; i++ {
		address := fmt.Sprintf("%s://%s:%d", protocol, currentIp.String(), port)
		result = append(result, address)
		incrementIp(currentIp)
	}
	return result
}

// Result type for a multiple connect request. It carries either a phone or an error.
// The intention behind this struct is to simulate a channel having two types.
type ConnectResult struct {
	Address string
	Phone   *Phone
	Error   error
}

type ConnectResults chan ConnectResult

// Connects to all phones given by the addresses parameter parallely and
// returns the results in the channel. Depending on how fast the real telephones answer
// to the login, the order in which the telephones are put in the channel differs
// from the order defined in the addresses parameter.
//
// This method returns immediately; the results channel is closed once all telephones
// have been contacted.
func (c *Connector) MultipleConnect(addresses ...string) ConnectResults {
	var wg sync.WaitGroup
	results := make(chan ConnectResult)
	for index, address := range addresses {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()
			phone, err := c.SingleConnect(address)
			if err != nil {
				err = fmt.Errorf("could not connect to address %d \"%s\": %v", index, address, err)
			}
			results <- ConnectResult{
				Address: address,
				Phone:   phone,
				Error:   err,
			}
		}(index, address)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
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

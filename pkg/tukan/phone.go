package tukan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/up"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// A connector is used to obtain a login token from a telephone and
// perform REST actions on the telephone. The connector can be used to either
// connect to only one telephone or to a bunch of telephones at the same time.
type Connector struct {
	Client    *http.Client
	UserName  string
	Password  string
	Addresses []string
}

// Tries to log in to a specific telephone identified by its Address.
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
		Address: address,
	}, nil
}

// Expands IP Addresses. If an passed Address cannot be parsed, then it is returned as is.
func ExpandAddresses(protocol string, addresses ...string) []string {
	result := make([]string, 0, len(addresses))
	for _, originalAddress := range addresses {
		address, port, number := splitAddress(originalAddress)
		created := CreateAddresses(protocol, address, port, number+1)
		for _, expanded := range created {
			result = append(result, expanded)
		}
		if len(created) == 0 {
			result = append(result, originalAddress)
		}
	}
	return result
}

func splitAddress(address string) (string, int, int) {
	modifier := strings.Index(address, "+")
	number := 0
	if modifier != -1 {
		number, _ = strconv.Atoi(address[modifier+1:])
		number = int(math.Max(float64(number), 0))
	} else {
		modifier = len(address)
	}
	colon := strings.Index(address, ":")
	port := 80
	if colon != -1 {
		port, _ = strconv.Atoi(address[colon+1 : modifier])
	} else {
		colon = modifier
	}
	return address[0:colon], port, number
}

// Creates number addresses to connect to phones. This first Address is given by startIp, which
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

type PhoneResult struct {
	Address string
	Error   error
}

type ResultCallback func(p *PhoneResult)

type PhoneAction func(p *Phone)

func (c *Connector) Run(loginCallback ResultCallback, operation PhoneAction, logoutCallback ResultCallback) {
	var wg sync.WaitGroup
	for index, address := range c.Addresses {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()
			phone, err := c.SingleConnect(address)
			loginCallback(&PhoneResult{Address: address, Error: err})
			if err != nil || phone == nil {
				return
			}
			defer func() {
				err = phone.Logout()
				logoutCallback(&PhoneResult{Address: address, Error: err})
			}()
			operation(phone)
		}(index, address)
	}
	wg.Wait()
}

// A phone represents a http Client that talks to exactly on
// physical telephone. A phone needs to be created with a Connector (see example).
// It is strongly recommended to defer calling the method Phone#Logout() because
// most IP620/630 only allow one active token at a time.
type Phone struct {
	client  *http.Client
	token   string
	Address string
	invalid bool
}

func (p *Phone) Host() string {
	return p.Address
}

func (p *Phone) Token() string {
	return p.token
}

// Sends a logout request to the phone. If the request passes without error
// then the token of the phone is reset. Further usage of the phone struct
// will most likely not work. If and error is returned, then the token stored
// in this telephone may or may not be used again, depending on the error.
func (p *Phone) Logout() error {
	url := fmt.Sprintf("%s/Logout", p.Address)
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+p.token)
	resp, err := p.client.Do(request)
	err = checkResponse(resp, err)
	if err == nil {
		p.token = ""
	}
	return err
}

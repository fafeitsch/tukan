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

// Expands IP Addresses. If an passed address cannot be parsed, then it is returned as is.
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
// The intention behind this struct is to simulate a channel having several types.
type Connection struct {
	Address string
	Phone   *Phone
	Error   error
}

type Connections chan *Phone

type PhoneResult struct {
	Comment string
	Address string
	Error   error
}

func (p *PhoneResult) String() string {
	return fmt.Sprintf("%s: %t (%s)", p.Address, p.Error == nil, p.Comment)
}

type ResultCallback func(p *PhoneResult)

// Connects to all phones given by the addresses parameter parallely and
// returns the results in the channel. Depending on how fast the real telephones answer
// to the login, the order in which the telephones are put in the channel differs
// from the order defined in the addresses parameter.
//
// This method returns immediately; the results channel is closed once all telephones
// have been contacted. Erroneous connections are reported asynchronously via the onError callback.
func (c *Connector) MultipleConnect(onError ResultCallback, addresses ...string) Connections {
	var wg sync.WaitGroup
	results := make(chan *Phone)
	for index, address := range addresses {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()
			phone, err := c.SingleConnect(address)
			if err != nil {
				onError(&PhoneResult{Comment: err.Error(), Error: err, Address: address})
				return
			}
			results <- phone
		}(index, address)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func (c *Connector) Run(loginCallback ResultCallback, logoutCallback ResultCallback, operations ...func(p *Phone)) {
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
			for _, operation := range operations {
				operation(phone)
			}
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
	address string
	invalid bool
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

// Logouts from all phones in the channel and blocks until the channel is closed.
// Results (successful and erroneous) are reported asynchronously via the callback.
func (p Connections) Logout(onFinish ResultCallback) {
	var wg sync.WaitGroup
	for phone := range p {
		wg.Add(1)
		go func(phone *Phone) {
			defer wg.Done()
			err := phone.Logout()
			if err != nil {
				onFinish(&PhoneResult{Comment: err.Error(), Address: phone.address, Error: err})
			} else {
				onFinish(&PhoneResult{Comment: "logout successful", Address: phone.address})
			}
		}(phone)
	}
	wg.Wait()
}

func (p Connections) loop(singleAction func(phone *Phone), end func()) {
	var wg sync.WaitGroup
	for connection := range p {
		wg.Add(1)
		go func(phone *Phone) {
			defer wg.Done()
			singleAction(phone)
		}(connection)
	}
	go func() {
		wg.Wait()
		end()
	}()
}

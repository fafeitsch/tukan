package tukan

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func ExamplePhone() {
	connector := Connector{Client: http.DefaultClient, UserName: username, Password: password}
	phone, _ := connector.SingleConnect("http://example.com:8080")
	book, _ := phone.DownloadPhoneBook()
	fmt.Print(book)
}

var password = "ken sent me"
var username = "larry"

func TestConnector_SingleConnect(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	server := httptest.NewServer(handler)
	defer server.Close()
	connector := Connector{
		Client:   http.DefaultClient,
		UserName: username,
		Password: password,
	}
	t.Run("success", func(t *testing.T) {
		phone, err := connector.SingleConnect(server.URL)
		assert.NoError(t, err, "no error expected")
		assert.Equal(t, *telephone.Token, phone.token, "token not equal")
	})
	t.Run("invalid logins", func(t *testing.T) {
		connector.Password = ""
		phone, err := connector.SingleConnect(server.URL)
		assert.EqualError(t, err, "authentication error, status code: 403 with message \"403 Forbidden\"", "error message not as expected")
		assert.Nil(t, phone, "phone should be nil in case of error")
	})
}

func ExampleCreateAddresses() {
	addresses := CreateAddresses("http", "10.1.254.254", 8081, 3)
	fmt.Printf("Length of addresses is %d.\n", len(addresses))
	fmt.Printf("%s\n", addresses[0])
	fmt.Printf("%s\n", addresses[1])
	fmt.Printf("%s\n", addresses[2])
	// Output: Length of addresses is 3.
	// http://10.1.254.254:8081
	// http://10.1.254.255:8081
	// http://10.1.255.0:8081
}

func TestCreateAddresses(t *testing.T) {
	t.Run("invalid address", func(t *testing.T) {
		addresses := CreateAddresses("http", "not an ip", 8080, 4)
		assert.Empty(t, addresses, "addresses should be empty if the ip is not valid")
	})
	t.Run("success", func(t *testing.T) {
		addresses := CreateAddresses("http", "10.1.254.253", 8081, 5)
		want := []string{"http://10.1.254.253:8081", "http://10.1.254.254:8081", "http://10.1.254.255:8081", "http://10.1.255.0:8081", "http://10.1.255.1:8081"}
		assert.Equal(t, want, addresses, "generated addresses not correct")
	})
}

func TestConnector_MultipleConnect(t *testing.T) {
	handler1, telephone1 := mock.CreatePhone(username, password)
	handler2, telephone2 := mock.CreatePhone(username, password)
	server1 := httptest.NewServer(handler1)
	defer server1.Close()
	server2 := httptest.NewServer(handler2)
	defer server2.Close()
	connector := Connector{Client: http.DefaultClient, UserName: username, Password: password}

	var firstFound, secondFound, thirdFound bool
	onError := func(result *PhoneResult) {
		if result.Address == "htp://invalid_url" {
			assert.EqualError(t, result.Error, "Post \"htp://invalid_url/Login\": unsupported protocol scheme \"htp\"", "error message of third telephone is wrong")
			thirdFound = true
		}
	}

	results := connector.MultipleConnect(onError, server2.URL, server1.URL, "htp://invalid_url")
	for result := range results {
		if result.address == server1.URL {
			assert.Equal(t, *telephone1.Token, result.token, "token of first telephone wrong")
			firstFound = true
		} else if result.address == server2.URL {
			assert.Equal(t, *telephone2.Token, result.token, "token of second telephone wrong")
			secondFound = true
		} else {
			assert.Fail(t, "unkown address %s", result.address)
		}
	}
	assert.True(t, firstFound, "first telephone must be handled")
	assert.True(t, secondFound, "second telephone must be handled")
	assert.True(t, thirdFound, "third (erroneous) telephone must be handled")
}

func TestPhone_Logout(t *testing.T) {
	handler, _ := mock.CreatePhone(username, password)
	server := httptest.NewServer(handler)
	defer server.Close()
	connector := Connector{
		Client:   http.DefaultClient,
		UserName: username,
		Password: password,
	}
	t.Run("success", func(t *testing.T) {
		phone, err := connector.SingleConnect(server.URL)
		assert.NoError(t, err, "no error expected")
		err = phone.Logout()
		assert.NoError(t, err, "no error expected")
	})
	t.Run("invalid token", func(t *testing.T) {
		phone, err := connector.SingleConnect(server.URL)
		assert.NoError(t, err, "no error expected")
		phone.token = "invalid"
		err = phone.Logout()
		assert.EqualError(t, err, "authentication error, status code: 401 with message \"401 Unauthorized\"", "no error expected")
	})
}

func ExampleExpandAddresses() {
	addresses := ExpandAddresses("http", "127.0.0.1", "not an ip", "10.20.30.40+2", "20.20.20.20:8080+1", "30.30.30.30:1234")
	for _, address := range addresses {
		fmt.Printf("%s\n", address)
	}
	// Output: http://127.0.0.1:80
	// not an ip
	// http://10.20.30.40:80
	// http://10.20.30.41:80
	// http://10.20.30.42:80
	// http://20.20.20.20:8080
	// http://20.20.20.21:8080
	// http://30.30.30.30:1234
}

func phonesTestSuite(t *testing.T, successMessage string, underTest func(Connections, ResultCallback) Connections, phoneSetup func(telephone *mock.Telephone)) *mock.Telephone {
	handler1, telephone1 := mock.CreatePhone(username, password)
	phoneSetup(telephone1)
	handler2, _ := mock.CreatePhone(username, password)
	server1 := httptest.NewServer(handler1)
	defer server1.Close()
	server2 := httptest.NewServer(handler2)
	defer server2.Close()

	fail := func(result *PhoneResult) {
		assert.Fail(t, "connecting should not fail: %v", result)
	}
	connector := Connector{UserName: username, Password: password, Client: http.DefaultClient}
	channel := connector.MultipleConnect(fail, server1.URL, server2.URL)
	transformed := make(Connections)
	go func() {
		for phone := range channel {
			if phone.address == server2.URL {
				phone.token = "faked"
			}
			transformed <- phone
		}
		close(transformed)
	}()
	counter := int32(0)
	onProcess := func(result *PhoneResult) {
		atomic.AddInt32(&counter, 1)
		if result.Address == server2.URL {
			assert.Equal(t, "authentication error, status code: 401 with message \"401 Unauthorized\"", result.Comment)
		} else if result.Address == server1.URL {
			assert.Equal(t, successMessage, result.Comment)
		} else {
			assert.Fail(t, "unexpected result server URL: %v", result)
		}
	}
	underTest(transformed, onProcess).Logout(func(p *PhoneResult) {})
	assert.Equal(t, 2, int(counter), "two results should be reported to onProcess")
	return telephone1
}

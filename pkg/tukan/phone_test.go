package tukan

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	t.Run("invalid Address", func(t *testing.T) {
		addresses := CreateAddresses("http", "not an ip", 8080, 4)
		assert.Empty(t, addresses, "addresses should be empty if the ip is not valid")
	})
	t.Run("success", func(t *testing.T) {
		addresses := CreateAddresses("http", "10.1.254.253", 8081, 5)
		want := []string{"http://10.1.254.253:8081", "http://10.1.254.254:8081", "http://10.1.254.255:8081", "http://10.1.255.0:8081", "http://10.1.255.1:8081"}
		assert.Equal(t, want, addresses, "generated addresses not correct")
	})
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

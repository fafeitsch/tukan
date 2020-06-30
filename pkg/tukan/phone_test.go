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

	results := connector.MultipleConnect(server2.URL, server1.URL, "htp://invalid_url")
	var firstFound, secondFound, thirdFound bool
	for result := range results {
		if result.Phone != nil && result.Phone.address == server1.URL {
			assert.Equal(t, *telephone1.Token, result.Phone.token, "token of first telephone wrong")
			firstFound = true
		}
		if result.Phone != nil && result.Phone.address == server2.URL {
			assert.Equal(t, *telephone2.Token, result.Phone.token, "token of second telephone wrong")
			secondFound = true
		}
		if result.Phone == nil {
			assert.EqualError(t, result.Error, "could not connect to address 2 \"htp://invalid_url\": Post \"htp://invalid_url/Login\": unsupported protocol scheme \"htp\"", "error message of third telephone is wrong")
			thirdFound = true
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

func TestConnectResults_Scan(t *testing.T) {
	handler1, _ := mock.CreatePhone(username, password)
	handler2, _ := mock.CreatePhone(username, password)
	handler3, _ := mock.CreatePhone("abc", "123")
	server1 := httptest.NewServer(handler1)
	defer server1.Close()
	server2 := httptest.NewServer(handler2)
	defer server2.Close()
	server3 := httptest.NewServer(handler3)
	defer server3.Close()

	connector := Connector{UserName: username, Password: password, Client: http.DefaultClient}
	channel := connector.MultipleConnect(server1.URL, server2.URL, server3.URL)
	transformed := make(Connections)
	got := transformed.Scan()
	for connectionResult := range channel {
		if connectionResult.Address == server2.URL {
			connectionResult.Phone.token = "faked"
		}
		transformed <- connectionResult
	}
	close(transformed)
	var phone1, phone2, phone3 bool
	counter := 0
	for result := range got {
		switch result.Address {
		case server1.URL:
			assert.True(t, result.Success, "the first telephone should be scanned successfully")
			assert.Equal(t, "connection established and login successful", result.Comment, "comment of first telephone is wrong")
			phone1 = true
			break
		case server2.URL:
			assert.False(t, result.Success, "the second telephone should not be scanned successfully")
			assert.Equal(t, "connection established and login successful, but logout not: authentication error, status code: 401 with message \"401 Unauthorized\"", result.Comment, "comment of the second telephone is wrong")
			phone2 = true
			break
		case server3.URL:
			assert.False(t, result.Success, "the third telephone should not be scanned successfully")
			expected := fmt.Sprintf("could not connect: could not connect to address 2 \"%s\": authentication error, status code: 403 with message \"403 Forbidden\"", server3.URL)
			assert.Equal(t, expected, result.Comment, "comment of the third telephone is wrong")
			phone3 = true
			break
		}
		counter = counter + 1
	}
	assert.True(t, phone1, "first phone should be detected")
	assert.True(t, phone2, "second phone should be detected")
	assert.True(t, phone3, "third phone should be detected")
	assert.Equal(t, 3, counter, "there should be exactly three scan results")
}

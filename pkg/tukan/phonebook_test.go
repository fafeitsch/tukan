package tukan

import (
	"github.com/fafeitsch/Tukan/pkg/tukan/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestPhone_UploadPhoneBook(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	server := httptest.NewServer(handler)
	defer server.Close()
	payloadContent := "<phonebook>some dummy entries</phonebook>"
	connector := Connector{Client: http.DefaultClient, UserName: username, Password: password}
	phone, err := connector.SingleConnect(server.URL)
	defer func() { _ = phone.Logout() }()
	require.NoError(t, err, "no error expected")
	t.Run("success", func(t *testing.T) {
		assert.NotEqual(t, payloadContent+"\n", telephone.Phonebook, "before the upload, the phone book must not be equal to the expected one")
		err = phone.UploadPhoneBook(payloadContent)
		require.NoError(t, err, "no error expected")
		assert.Equal(t, payloadContent+"\n", telephone.Phonebook, "uploaded phone book not correct")
	})
}

func TestPhone_DownloadPhoneBook(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	server := httptest.NewServer(handler)
	defer server.Close()
	connector := Connector{Client: http.DefaultClient, UserName: username, Password: password}
	phone, err := connector.SingleConnect(server.URL)
	require.NoError(t, err, "no error expected")
	telephone.Phonebook = "this is a telephone book"
	t.Run("success", func(t *testing.T) {
		book, err := phone.DownloadPhoneBook()
		require.NoError(t, err, "no error expected")
		assert.Equal(t, telephone.Phonebook, *book, "downloaded phone book is wrong")
	})
	t.Run("error", func(t *testing.T) {
		phone.token = "abc"
		book, err := phone.DownloadPhoneBook()
		require.EqualError(t, err, "authentication error, status code: 401 with message \"401 Unauthorized\"", "error expected")
		assert.Nil(t, book, "result should be nil in case of an error")
	})
}

func TestConnections_UploadPhoneBook(t *testing.T) {
	handler1, telephone1 := mock.CreatePhone(username, password)
	telephone1.Phonebook = ""
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
			assert.Equal(t, "Upload successful", result.Comment)
		} else {
			assert.Fail(t, "unexpected result server URL: %v", result)
		}
	}
	transformed.UploadPhoneBook(onProcess, "this is my phone book").
		Logout(func(p *PhoneResult) {})
	assert.Equal(t, 2, int(counter), "two results should be reported to onProcess")
	assert.Equal(t, "this is my phone book\n", telephone1.Phonebook, "phone book of first telephone should be changed")
}

func TestConnections_DownloadPhoneBook(t *testing.T) {
	handler1, telephone1 := mock.CreatePhone(username, password)
	telephone1.Phonebook = "book of telephone 1"
	handler2, _ := mock.CreatePhone(username, password)
	server1 := httptest.NewServer(handler1)
	defer server1.Close()
	server2 := httptest.NewServer(handler2)
	defer server2.Close()

	connector := Connector{UserName: username, Password: password, Client: http.DefaultClient}
	fail := func(result *PhoneResult) {
		assert.Fail(t, "connecting should not fail: %v", result)
	}
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
			assert.Equal(t, "Download successful", result.Comment)
		} else {
			assert.Fail(t, "unexpected result server URL: %v", result)
		}
	}
	successCounter := int32(0)
	onSuccess := func(result *PhoneBookResult) {
		atomic.AddInt32(&successCounter, 1)
		assert.Equal(t, telephone1.Phonebook, result.PhoneBook, "expected phone book wrong")
	}
	transformed.DownloadPhoneBook(onProcess, onSuccess).
		Logout(func(p *PhoneResult) {})
	assert.Equal(t, 2, int(counter), "expected two results to be reported")
	assert.Equal(t, 1, int(successCounter), "there should be one success")
}

package tukan

import (
	"github.com/fafeitsch/Tukan/pkg/tukan/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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

func TestPrepareUploadPhoneBook(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	telephone.Phonebook = ""
	server := httptest.NewServer(handler)
	defer server.Close()
	var result *PhoneResult
	operation := PrepareUploadPhoneBook("this is my phone book", func(phoneResult *PhoneResult) { result = phoneResult })
	connector := &Connector{Client: http.DefaultClient, UserName: username, Password: password}
	phone, err := connector.SingleConnect(server.URL)
	require.NoError(t, err, "no error expected")
	t.Run("success", func(t *testing.T) {
		operation(phone)
		assert.Equal(t, server.URL, result.Address, "result callback not called properly")
		assert.NoError(t, result.Error, "no error expected")
		assert.Equal(t, "this is my phone book\n", telephone.Phonebook)
	})
	t.Run("failure", func(t *testing.T) {
		phone.token = "invalid"
		operation(phone)
		assert.EqualError(t, result.Error, "authentication error, status code: 401 with message \"401 Unauthorized\"", "error not set correctly in callback result")
	})
}

func TestPrepareDownloadPhoneBook(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	telephone.Phonebook = "This is a cool book."
	server := httptest.NewServer(handler)
	defer server.Close()
	var result *PhoneBookResult
	operation := PrepareDownloadPhoneBook(func(phoneResult *PhoneBookResult) { result = phoneResult })
	connector := &Connector{Client: http.DefaultClient, UserName: username, Password: password}
	phone, err := connector.SingleConnect(server.URL)
	require.NoError(t, err, "no error expected")
	t.Run("success", func(t *testing.T) {
		operation(phone)
		assert.Equal(t, server.URL, result.PhoneResult.Address, "result callback not called properly")
		assert.NoError(t, result.Error, "no error expected")
		assert.Equal(t, "This is a cool book.", *result.PhoneBook)
	})
	t.Run("failure", func(t *testing.T) {
		phone.token = "invalid"
		operation(phone)
		assert.EqualError(t, result.Error, "authentication error, status code: 401 with message \"401 Unauthorized\"", "error not set correctly in callback result")
	})
}

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
	phoneSetup := func(phone *mock.Telephone) {
		phone.Phonebook = ""
	}
	underTest := func(connections Connections, callback ResultCallback) Connections {
		return connections.UploadPhoneBook(callback, "this is my phone book")
	}
	phone := phonesTestSuite(t, "Upload successful", underTest, phoneSetup)
	assert.Equal(t, "this is my phone book\n", phone.Phonebook, "phone book of first telephone should be changed")
}

func TestConnections_DownloadPhoneBook(t *testing.T) {
	phoneSetup := func(phone *mock.Telephone) {
		phone.Phonebook = "book of telephone 1"
	}
	successCounter := int32(0)
	onSuccess := func(result *PhoneBookResult) {
		atomic.AddInt32(&successCounter, 1)
		assert.Equal(t, "book of telephone 1", result.PhoneBook, "expected phone book wrong")
	}
	underTest := func(connections Connections, callback ResultCallback) Connections {
		return connections.DownloadPhoneBook(callback, onSuccess)
	}
	_ = phonesTestSuite(t, "Download successful", underTest, phoneSetup)
	assert.Equal(t, 1, int(successCounter), "there should be one success")
}

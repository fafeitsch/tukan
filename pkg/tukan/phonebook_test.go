package tukan

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/tukan/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
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
	got := transformed.UploadPhoneBook("this is my phone book")
	for connectionResult := range channel {
		if connectionResult.Address == server2.URL {
			connectionResult.Phone.token = "faked"
		}
		transformed <- connectionResult
	}
	close(transformed)
	want1 := SimpleResult{Address: server1.URL, Success: true, Comment: "Upload successful"}
	want2 := SimpleResult{Address: server2.URL, Success: false, Comment: "authentication error, status code: 401 with message \"401 Unauthorized\""}
	want3 := SimpleResult{Address: server3.URL, Success: false, Comment: fmt.Sprintf("could not connect to address 2 \"%s\": authentication error, status code: 403 with message \"403 Forbidden\"", server3.URL)}
	checkSimpleResults(t, got, want1, want2, want3)
	assert.Equal(t, "this is my phone book\n", telephone1.Phonebook, "phone book of first telephone should be changed")
}

func TestConnections_DownloadPhoneBook(t *testing.T) {
	handler1, telephone1 := mock.CreatePhone(username, password)
	telephone1.Phonebook = "book of telephone 1"
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
	got := transformed.DownloadPhoneBook()
	for connectionResult := range channel {
		if connectionResult.Address == server2.URL {
			connectionResult.Phone.token = "faked"
		}
		transformed <- connectionResult
	}
	close(transformed)
	simpleResults := make(chan SimpleResult)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for bookResult := range got {
			if bookResult.Address == server1.URL {
				assert.Equal(t, "book of telephone 1", *bookResult.PhoneBook, "book of telephone 1 should be downloaded correctly")
			}
			simpleResults <- bookResult.SimpleResult
		}
		close(simpleResults)
	}()
	want1 := SimpleResult{Address: server1.URL, Success: true, Comment: "Download successful"}
	want2 := SimpleResult{Address: server2.URL, Success: false, Comment: "authentication error, status code: 401 with message \"401 Unauthorized\""}
	want3 := SimpleResult{Address: server3.URL, Success: false, Comment: fmt.Sprintf("could not connect to address 2 \"%s\": authentication error, status code: 403 with message \"403 Forbidden\"", server3.URL)}
	checkSimpleResults(t, simpleResults, want1, want2, want3)
	wg.Wait()
}

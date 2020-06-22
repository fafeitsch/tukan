package tukan

import (
	"github.com/fafeitsch/Tukan/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPhone_UploadPhoneBook(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	srv := httptest.NewServer(handler)
	defer srv.Close()
	payloadContent := "<phonebook>some dummy entries</phonebook>"
	phone, err := Connect(http.DefaultClient, srv.URL, username, password)
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
	srv := httptest.NewServer(handler)
	defer srv.Close()
	phone, err := Connect(http.DefaultClient, srv.URL, username, password)
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
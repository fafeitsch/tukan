package tukan

import (
	"github.com/fafeitsch/Tukan/pkg/tukan/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var password = "ken sent me"
var username = "larry"
var token = "abc1234"

func TestConnect(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	server := httptest.NewServer(handler)
	defer server.Close()
	t.Run("success", func(t *testing.T) {
		phone, err := Connect(http.DefaultClient, server.URL, username, password)
		assert.NoError(t, err, "no error expected")
		assert.Equal(t, *telephone.Token, phone.token, "token not equal")
	})
	t.Run("invalid logins", func(t *testing.T) {
		phone, err := Connect(http.DefaultClient, server.URL, "", "")
		assert.EqualError(t, err, "authentication error, status code: 403 with message \"403 Forbidden\"", "error message not as expected")
		assert.Nil(t, phone, "phone should be nil in case of error")
	})
}

func TestPhone_Logout(t *testing.T) {
	handler, _ := mock.CreatePhone(username, password)
	server := httptest.NewServer(handler)
	defer server.Close()
	t.Run("success", func(t *testing.T) {
		phone, err := Connect(http.DefaultClient, server.URL, username, password)
		assert.NoError(t, err, "no error expected")
		err = phone.Logout()
		assert.NoError(t, err, "no error expected")
	})
	t.Run("invalid token", func(t *testing.T) {
		phone, err := Connect(http.DefaultClient, server.URL, username, password)
		assert.NoError(t, err, "no error expected")
		phone.token = "invalid"
		err = phone.Logout()
		assert.EqualError(t, err, "authentication error, status code: 401 with message \"401 Unauthorized\"", "no error expected")
	})
}

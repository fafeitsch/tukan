package http

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogout(t *testing.T) {
	token := "token for testing"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer "+token {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	tokener := tokenerImpl{
		client: srv.Client(),
		port:   port,
	}
	t.Run("successful", func(t *testing.T) {
		err := tokener.logout(ip, token)
		require.NoError(t, err, "no error expected")
	})
	t.Run("not successful", func(t *testing.T) {
		err := tokener.logout(ip, "wrong token")
		assert.EqualError(t, err, "authentication error, status code: 401 with message \"401 Unauthorized\"", "error message not as expected")
	})
}

func TestGetTokenSuccess(t *testing.T) {
	telephone := mock.Telephone{Login: "a_user", Password: "a_password"}
	srv := httptest.NewServer(http.HandlerFunc(telephone.AttemptLogin))
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	tokener := tokenerImpl{
		client:   srv.Client(),
		port:     port,
		login:    "a_user",
		password: "a_password",
	}
	token, err := tokener.fetchToken(ip)
	require.Nil(t, err, "no error expected")
	assert.NotNil(t, token, "token should not be nil")
	defer func() { _ = tokener.logout(ip, *token) }()
	assert.Equal(t, telephone.Token, token, "tokens should be equal")
}

func TestGetTokenUnparsableAnswer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "this is not json")
	}))
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	tokener := tokenerImpl{
		client:   srv.Client(),
		port:     port,
		login:    "a_user",
		password: "a_password",
	}
	token, err := tokener.fetchToken(ip)
	assert.Nil(t, token, "token should be nil")
	assert.EqualError(t, err, "could not unmarshal token from 127.0.0.1: invalid character 'h' in literal true (expecting 'r')", "error message not correct")
}

func TestGetTokenWrongPassword(t *testing.T) {
	telephone := mock.Telephone{Login: "a_user", Password: "a_password"}
	srv := httptest.NewServer(http.HandlerFunc(telephone.AttemptLogin))
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	tokener := tokenerImpl{client: srv.Client(), port: port, login: "a_user", password: "wrong_password"}
	token, err := tokener.fetchToken(ip)
	assert.Nil(t, token, "token should be nil")
	assert.EqualError(t, err, "authentication error, status code: 403 with message \"403 Forbidden\"", "error message is wrong")
}

package http

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/mock"
	"github.com/stretchr/testify/assert"
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
	pc := PhoneClient{
		Client: srv.Client(),
		Port:   port,
	}
	t.Run("successful", func(t *testing.T) {
		logger := prepareLogger()
		defer resetLogger()
		pc.logout(ip, token)
		got := logger.String()
		want := fmt.Sprintf("logging out of %s (http://%s:%d/Logout) …\nlogout of %s successful\n", ip, ip, port, ip)
		assert.Equal(t, want, got, "logging result wrong")
	})
	t.Run("not successful", func(t *testing.T) {
		logger := prepareLogger()
		defer resetLogger()
		pc.logout(ip, "wrong token")
		got := logger.String()
		want := fmt.Sprintf("logging out of %s (http://%s:%d/Logout) …\ncould not logout from %s: authentication error, status code: 401 with message \"401 Unauthorized\"\n", ip, ip, port, ip)
		assert.Equal(t, want, got, "logging result wrong")
	})
}

func TestGetTokenSuccess(t *testing.T) {
	telephone := mock.Telephone{Login: "a_user", Password: "a_password"}
	srv := httptest.NewServer(http.HandlerFunc(telephone.AttemptLogin))
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	pc := PhoneClient{
		Client:   srv.Client(),
		Port:     port,
		Login:    "a_user",
		Password: "a_password",
	}
	logger := prepareLogger()
	defer resetLogger()
	token := pc.fetchToken(ip)
	assert.NotNil(t, token, "token should not be nil")
	defer pc.logout(ip, *token)
	assert.Equal(t, telephone.Token, token, "tokens should be equal")
	logged := logger.String()
	want := fmt.Sprintf("fetching token for %s (http://%s:%d/Login) …\nFetching token for %s successful\n", ip, ip, port, ip)
	assert.Equal(t, want, logged, "logging was wrong")
}

func TestGetTokenUnparsableAnswer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "this is not json")
	}))
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	pc := PhoneClient{
		Client:   srv.Client(),
		Port:     port,
		Login:    "a_user",
		Password: "a_password",
	}
	logger := prepareLogger()
	defer resetLogger()
	token := pc.fetchToken(ip)
	assert.Nil(t, token, "token should not nil")
	logged := logger.String()
	want := fmt.Sprintf("fetching token for %s (http://%s:%d/Login) …\ncould not unmarshal token from %s: invalid character 'h' in literal true (expecting 'r')\n", ip, ip, port, ip)
	assert.Equal(t, want, logged, "logging was wrong")
}

func TestGetToken(t *testing.T) {
	telephone := mock.Telephone{Login: "a_user", Password: "a_password"}
	srv := httptest.NewServer(http.HandlerFunc(telephone.AttemptLogin))
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	tokener := tokenerImpl{port: port, login: "a_user", password: "wrong_password"}
	logger := prepareLogger()
	defer resetLogger()
	token, err := tokener.fetchToken(ip)
	assert.Nil(t, token, "token should be nil")
	assert.EqualError(t, "", err, "error message is wrong")
	logged := logger.String()
	want := fmt.Sprintf("fetching token for %s (http://%s:%d/Login) …\nauthentication error, status code: 403 with message \"403 Forbidden\"\n", ip, ip, port)
	assert.Equal(t, want, logged, "logging was wrong")
}

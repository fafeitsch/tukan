package http

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"testing"
)

var pattern = regexp.MustCompile("http://\\[?([^:]+)]?:(\\d*)")

func parseTestServerURL(url string) (string, int) {
	match := pattern.FindStringSubmatch(url)
	if len(match) != 3 {
		log.Fatalf("url \"%s\" has not the required format", url)
	}
	port, _ := strconv.Atoi(match[2])
	return match[1], port
}

func resetLogger() {
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stdout)
}

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
		logger := bytes.Buffer{}
		log.SetOutput(&logger)
		log.SetFlags(0)
		defer resetLogger()
		pc.logout(ip, token)
		got := logger.String()
		want := fmt.Sprintf("logging out of %s (http://%s:%d/Logout) …\nlogout of %s successful\n", ip, ip, port, ip)
		assert.Equal(t, want, got, "logging result wrong")
	})
	t.Run("not successful", func(t *testing.T) {
		logger := bytes.Buffer{}
		log.SetOutput(&logger)
		log.SetFlags(0)
		defer resetLogger()
		pc.logout(ip, "wrong token")
		got := logger.String()
		want := fmt.Sprintf("logging out of %s (http://%s:%d/Logout) …\ncould not logout from %s: unexpected status code: 401 with message \"401 Unauthorized\"\n", ip, ip, port, ip)
		assert.Equal(t, want, got, "logging result wrong")
	})
}

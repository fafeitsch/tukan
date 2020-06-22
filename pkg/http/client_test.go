package http

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"os"
	"testing"
)

func prepareLogger() *bytes.Buffer {
	logger := &bytes.Buffer{}
	log.SetOutput(logger)
	log.SetFlags(0)
	return logger
}

func resetLogger() {
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stdout)
}

func TestPhoneClient_Scan(t *testing.T) {
	t.SkipNow()
	pc := BuildPhoneClient(8080, "username", "password", 5)
	logWriter := &bytes.Buffer{}
	logger := log.New(logWriter, "LOGGING ", 0)
	pc.Logger = logger
	result := pc.Scan("10.10.40.1", 4)
	require.Equal(t, 4, len(result), "number of results not correct")
	require.Equal(t, "phone is reachable, login worked", result["10.10.40.1"], "result of 10.10.40.1 incorrect")
	require.Equal(t, "login failed", result["10.10.40.2"], "result of 10.10.40.2 incorrect")
	require.Equal(t, "phone is reachable, login worked", result["10.10.40.3"], "result of 10.10.40.3 incorrect")
	require.Equal(t, "logout failed", result["10.10.40.4"], "result of 10.10.40.4 incorrect")
	want := `LOGGING fetching token for 10.10.40.1…
LOGGING 10.10.40.1 is reachable and login is possible
LOGGING logging out of 10.10.40.1…
LOGGING fetching token for 10.10.40.2…
LOGGING fetching token for 10.10.40.2 failed: authentication failed
LOGGING fetching token for 10.10.40.3…
LOGGING 10.10.40.3 is reachable and login is possible
LOGGING logging out of 10.10.40.3…
LOGGING fetching token for 10.10.40.4…
LOGGING 10.10.40.4 is reachable and login is possible
LOGGING logging out of 10.10.40.4…
LOGGING could not logout from 10.10.40.4
`
	require.Equal(t, want, logWriter.String(), "logger output wrong")
}

func TestCheckResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		err        error
		wantError  error
	}{
		{name: "all ok", statusCode: http.StatusOK, message: "OK", err: nil, wantError: nil},
		{name: "no connect", statusCode: http.StatusOK, message: "OK", err: fmt.Errorf("no route to host"), wantError: fmt.Errorf("no route to host")},
		{name: "unauthorized", statusCode: http.StatusUnauthorized, message: "Not authorized", err: nil, wantError: fmt.Errorf("authentication error, status code: %d with message \"Not authorized\"", http.StatusUnauthorized)},
		{name: "server error", statusCode: http.StatusInternalServerError, message: "Not available", err: nil, wantError: fmt.Errorf("unexpected status code: %d with message \"Not available\"", http.StatusInternalServerError)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := http.Response{StatusCode: tt.statusCode, Status: tt.message}
			got := checkResponse(&resp, tt.err)
			assert.Equal(t, tt.wantError, got)
		})
	}
}

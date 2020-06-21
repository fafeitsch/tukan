package http

import (
	"bytes"
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/domain"
	"github.com/fafeitsch/Tukan/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestPhoneClient_UploadPhoneBook(t *testing.T) {
	handler, telephone := mock.CreatePhone("username", "password")
	srv := httptest.NewServer(handler)
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	delimiter := "DELIMITER"
	payloadContent := "<phonebook>some dummy entries</phonebook>"
	payload := domain.InsertIntoTemplate(payloadContent, delimiter)
	pc := BuildPhoneClient(port, "username", "password", 5)
	pc.client = srv.Client()
	t.Run("success", func(t *testing.T) {
		logWriter := &bytes.Buffer{}
		logger := log.New(logWriter, "TESTING ", 0)
		pc.Logger = logger
		result := pc.UploadPhoneBook(ip, 1, payload, delimiter)
		assert.Equal(t, payloadContent+"\n", telephone.Phonebook, "uploaded phone book wrong")
		assert.Equal(t, 1, len(result), "number of results wrong")
		assert.Equal(t, "uploading phone book successful", result[ip], "result string wrong")
		wantTemplate := "TESTING fetching token for %s…\n" +
			"TESTING starting up of phone book to %s\n" +
			"TESTING uploaded phone book successfully to %s\n" +
			"TESTING logging out of %s…\n"
		want := fmt.Sprintf(wantTemplate, ip, ip, ip, ip)
		assert.Equal(t, want, logWriter.String(), "logging output is wrong")
	})
	t.Run("no success", func(t *testing.T) {
		logWriter := &bytes.Buffer{}
		logger := log.New(logWriter, "TESTING ", 0)
		pc.Logger = logger
		result := pc.UploadPhoneBook(ip, 1, "wrong payload", delimiter)
		assert.Equal(t, 1, len(result), "number of results wrong")
		assert.Equal(t, "uploading phone book failed", result[ip], "result string wrong")
		wantTemplate := "TESTING fetching token for %s…\n" +
			"TESTING starting up of phone book to %s\n" +
			"TESTING could not up phone book to %s: unexpected status code: 400 with message \"400 Bad Request\"\n" +
			"TESTING logging out of %s…\n"
		want := fmt.Sprintf(wantTemplate, ip, ip, ip, ip)
		assert.Equal(t, want, logWriter.String(), "logging output is wrong")
	})
}

func TestPhoneClient_DownloadPhoneBook(t *testing.T) {
	handler, telephone := mock.CreatePhone("username", "password")
	srv := httptest.NewServer(handler)
	defer srv.Close()
	ip, port := parseTestServerURL(srv.URL)
	telephone.Phonebook = "this is a telephone book"
	pc := BuildPhoneClient(port, "username", "password", 5)
	pc.client = srv.Client()
	t.Run("success", func(t *testing.T) {
		logWriter := &bytes.Buffer{}
		logger := log.New(logWriter, "TESTING ", 0)
		pc.Logger = logger
		resultMap, result := pc.DownloadPhoneBook(ip)
		assert.Equal(t, telephone.Phonebook, result, "uploaded phone book wrong")
		assert.Equal(t, 1, len(resultMap), "number of results wrong")
		assert.Equal(t, "downloading phone book successful", resultMap[ip], "result string wrong")
		wantTemplate := "TESTING fetching token for %s…\n" +
			"TESTING start phone book download from %s…\n" +
			"TESTING phone book download from %s successful\n" +
			"TESTING logging out of %s…\n"
		want := fmt.Sprintf(wantTemplate, ip, ip, ip, ip)
		assert.Equal(t, want, logWriter.String(), "logging output is wrong")
	})
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

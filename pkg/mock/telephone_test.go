package mock

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTelephone_AttemptLogin(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		header     string
		user       string
		pw         string
		wantStatus int
		wantMsg    string
	}{
		{name: "success", method: "POST", header: "application/json", user: "username", pw: "pass", wantStatus: http.StatusOK, wantMsg: ""},
		{name: "wrong method", method: "GET", header: "text/plain", user: "username", pw: "pass", wantStatus: http.StatusMethodNotAllowed, wantMsg: "Unsupported method \"GET\", want method \"POST\""},
		{name: "wrong header", method: "POST", header: "text/plain", user: "username", pw: "pass", wantStatus: http.StatusUnsupportedMediaType, wantMsg: "Header \"Content-Type\" must begin with value \"application/json\", but was \"text/plain\""},
		{name: "wrong password", method: "POST", header: "application/json", user: "username", pw: "wrong", wantStatus: http.StatusForbidden, wantMsg: "provided credentials not valid"},
		{name: "wrong username", method: "POST", header: "application/json", user: "wrong", pw: "pass", wantStatus: http.StatusForbidden, wantMsg: "provided credentials not valid"},
		{name: "too much elements", method: "POST", header: "application/json", user: "wrong\"}", pw: "pass", wantStatus: http.StatusUnprocessableEntity, wantMsg: "request body contained more than one json object, which is not allowed"},
		{name: "unprocessable entity", method: "POST", header: "application/json", user: "wrong\"", pw: "pass", wantStatus: http.StatusUnprocessableEntity, wantMsg: "invalid character '\"' after object key:value pair"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := fmt.Sprintf("{\"login\":\"%s\", \"password\":\"%s\"}", tt.user, tt.pw)
			request := httptest.NewRequest(tt.method, "/Login", strings.NewReader(payload))
			request.Header.Add("Content-Type", tt.header)
			telephone := Telephone{Login: "username", Password: "pass"}
			recorder := httptest.NewRecorder()
			assert.Nil(t, telephone.Token, "token should be nil before anything happens on it")
			telephone.AttemptLogin(recorder, request)
			response := recorder.Result()
			data, _ := ioutil.ReadAll(response.Body)
			assert.Equal(t, tt.wantStatus, response.StatusCode, "the status code is wrong")
			if tt.wantStatus == http.StatusOK {
				require.NotNil(t, telephone.Token, "token should be set now")
				want := fmt.Sprintf("{\"token\":\"%s\"}", *telephone.Token)
				assert.Equal(t, want, string(data), "expected token result wrong")
			} else {
				assert.Nil(t, telephone.Token, "token should be nil if login was not successful")
				assert.Equal(t, tt.wantMsg, string(data), "response payload is wrong")
			}
		})
	}
}

func TestTelephone_PostPhoneBook(t *testing.T) {
	payload := domain.InsertIntoTemplate("hooray, a phonebook", "BOUNDARY-42")
	tests := []struct {
		name       string
		method     string
		header     string
		payload    string
		wantStatus int
		wantMsg    string
	}{
		{name: "success", method: "POST", header: "multipart/form-data; boundary=BOUNDARY-42", payload: payload, wantStatus: http.StatusNoContent, wantMsg: ""},
		{name: "wrong method", method: "GET", header: "multipart/form-data; boundary=BOUNDARY-42", payload: payload, wantStatus: http.StatusMethodNotAllowed, wantMsg: "Unsupported method \"GET\", want method \"POST\""},
		{name: "incomplete header", method: "POST", header: "multipart/form-data", payload: payload, wantStatus: http.StatusUnsupportedMediaType, wantMsg: "Header \"Content-Type\" must begin with value \"multipart/form-data; boundary=\", but was \"multipart/form-data\""},
		{name: "unparsable", method: "POST", header: "multipart/form-data; boundary=a", payload: payload, wantStatus: http.StatusBadRequest, wantMsg: "could not parse multipart-form"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/LocalPhoneBook", strings.NewReader(payload))
			request.Header.Add("Content-Type", tt.header)
			telephone := Telephone{}
			assert.Empty(t, telephone.Phonebook, "phonebook should be empty before anything happens")
			recorder := httptest.NewRecorder()
			telephone.PostPhoneBook(recorder, request)
			response := recorder.Result()
			assert.Equal(t, tt.wantStatus, response.StatusCode, "status code is wrong")
			data, _ := ioutil.ReadAll(response.Body)
			assert.Equal(t, tt.wantMsg, string(data), "received payload is wrong")
		})
	}
}

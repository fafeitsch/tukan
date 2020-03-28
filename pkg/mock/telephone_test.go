package mock

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api"
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
			status, data := getStatusAndData(recorder)
			assert.Equal(t, tt.wantStatus, status, "the status code is wrong")
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

func getStatusAndData(recorder *httptest.ResponseRecorder) (int, string) {
	response := recorder.Result()
	data, _ := ioutil.ReadAll(response.Body)
	return response.StatusCode, string(data)
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
			status, data := getStatusAndData(recorder)
			assert.Equal(t, tt.wantStatus, status, "status code is wrong")
			assert.Equal(t, tt.wantMsg, string(data), "received payload is wrong")
		})
	}
}

func TestTelephone_SaveLocalPhoneBook(t *testing.T) {
	phonebook := "this is the phonebook"
	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantMsg    string
	}{
		{name: "success", method: "GET", wantStatus: http.StatusOK, wantMsg: phonebook},
		{name: "wrong method", method: "POST", wantStatus: http.StatusMethodNotAllowed, wantMsg: "Unsupported method \"POST\", want method \"GET\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			telephone := Telephone{Phonebook: phonebook}
			request := httptest.NewRequest(tt.method, "/SaveLocalPhoneBook", strings.NewReader(""))
			recorder := httptest.NewRecorder()
			telephone.SaveLocalPhoneBook(recorder, request)
			status, data := getStatusAndData(recorder)
			assert.Equal(t, tt.wantStatus, status, "status code is wrong")
			assert.Equal(t, tt.wantMsg, string(data), "received phonebook wrong")
		})
	}
}

func TestTelephone_HandleParameters_GET(t *testing.T) {
	keys := []api.FunctionKey{
		{DisplayName: "Shep Alves", PhoneNumber: "854", CallPickupCode: "***"},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
		{DisplayName: "Koren Wolledge", PhoneNumber: "294", CallPickupCode: "##"},
		{DisplayName: "Ossi Lisimore", PhoneNumber: "929", CallPickupCode: "##"},
		{DisplayName: "Jordana Jeromson", PhoneNumber: "245", CallPickupCode: "##"},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
	}
	wantBytes, _ := ioutil.ReadFile("../mockdata/parameters.json")
	telephone := Telephone{Keys: api.FunctionKeys{FunctionKeys: keys}}
	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantMsg    string
	}{
		{name: "get parameters successfully", method: "GET", wantStatus: http.StatusOK, wantMsg: string(wantBytes)},
		{name: "wrong method", method: "PUT", wantStatus: http.StatusMethodNotAllowed, wantMsg: "The method \"PUT\" is not allowed. Only \"GET\" and \"POST\" are supported"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/Parameters", strings.NewReader(""))
			recorder := httptest.NewRecorder()
			telephone.HandleParameters(recorder, request)
			status, data := getStatusAndData(recorder)
			assert.Equal(t, tt.wantStatus, status, "status code is wrong")
			assert.Equal(t, tt.wantMsg, string(data), "received phonebook wrong")
		})
	}
}
package mock

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/tukan/params"
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
			telephone.attemptLogin(recorder, request)
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

const payloadTemplate = `--%s
Content-Disposition: form-data; name="file"; filename="LocalPhonebook.xml"
Content-Type: text/xml

%s

--%s--`

func TestTelephone_PostPhoneBook(t *testing.T) {
	payload := fmt.Sprintf(payloadTemplate, "BOUNDARY-42", "hooray, a phonebook", "BOUNDARY-42")
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
			telephone.postPhoneBook(recorder, request)
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
			telephone.saveLocalPhoneBook(recorder, request)
			status, data := getStatusAndData(recorder)
			assert.Equal(t, tt.wantStatus, status, "status code is wrong")
			assert.Equal(t, tt.wantMsg, string(data), "received phonebook wrong")
		})
	}
}

func TestTelephone_HandleParameters_GET(t *testing.T) {
	keys := []params.FunctionKey{
		{DisplayName: "Shep Alves", PhoneNumber: "854", CallPickupCode: "***"},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
		{DisplayName: "Koren Wolledge", PhoneNumber: "294", CallPickupCode: "##"},
		{DisplayName: "Ossi Lisimore", PhoneNumber: "929", CallPickupCode: "##"},
		{DisplayName: "Jordana Jeromson", PhoneNumber: "245", CallPickupCode: "##"},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
		{DisplayName: "", PhoneNumber: "", CallPickupCode: ""},
	}
	telephone := Telephone{Parameters: params.Parameters{FunctionKeys: keys}}
	t.Run("success", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/Parameters", strings.NewReader(""))
		recorder := httptest.NewRecorder()
		telephone.handleParameters(recorder, request)
		status, data := getStatusAndData(recorder)
		assert.Equal(t, http.StatusOK, status, "status code is wrong")
		got := params.Parameters{}
		err := json.Unmarshal([]byte(data), &got)
		require.NoError(t, err, "no error expected")
		assert.Equal(t, "Shep Alves", got.FunctionKeys[0].DisplayName)
		assert.Equal(t, "929", got.FunctionKeys[3].PhoneNumber)
		assert.Equal(t, "##", got.FunctionKeys[4].CallPickupCode)
	})
	t.Run("failure", func(t *testing.T) {
		request := httptest.NewRequest("PUT", "/Parameters", strings.NewReader(""))
		recorder := httptest.NewRecorder()
		telephone.handleParameters(recorder, request)
		status, data := getStatusAndData(recorder)
		assert.Equal(t, http.StatusMethodNotAllowed, status, "status code is wrong")
		assert.Equal(t, "The method \"PUT\" is not allowed. Only \"GET\" and \"POST\" are supported", string(data), "received phonebook wrong")

	})
}

func TestTelephone_HandleParameters_POST(t *testing.T) {
	payload := params.Parameters{FunctionKeys: []params.FunctionKey{{}, {DisplayName: "Ossi Lisimore"}}}
	marsh, err := json.Marshal(&payload)
	require.NoError(t, err, "no error expected")
	require.NoError(t, err, "no error expected")
	tests := []struct {
		name        string
		method      string
		contentType string
		payload     string
		wantStatus  int
		wantMsg     string
	}{
		{name: "success", method: "POST", contentType: "application/json", payload: string(marsh), wantStatus: http.StatusNoContent},
		{name: "wrong media type", method: "POST", contentType: "application/xml", payload: string(marsh), wantStatus: http.StatusUnsupportedMediaType, wantMsg: "contentType \"application/xml\" not supported"},
		{name: "invalid json", method: "POST", contentType: "application/json", payload: "{", wantStatus: http.StatusBadRequest, wantMsg: "could not deserialize json: unexpected EOF\n"},
		{name: "double json", method: "POST", contentType: "application/json", payload: "{}{}", wantStatus: http.StatusBadRequest, wantMsg: "request body contained more than one json object, which is not allowed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := []params.FunctionKey{
				{DisplayName: "Shep Alves", PhoneNumber: "854", CallPickupCode: "***"},
				{},
			}
			telephone := Telephone{Parameters: params.Parameters{FunctionKeys: keys}}
			request := httptest.NewRequest(tt.method, "/Parameters", strings.NewReader(tt.payload))
			request.Header.Add("Content-Type", tt.contentType)
			recorder := httptest.NewRecorder()
			telephone.handleParameters(recorder, request)
			status, data := getStatusAndData(recorder)
			assert.Equal(t, tt.wantStatus, status, "status code is wrong")
			assert.Equal(t, tt.wantMsg, data, "data is wrong")
			if status == http.StatusNoContent {
				assert.Equal(t, params.FunctionKey{}, telephone.Parameters.FunctionKeys[0])
				assert.Equal(t, "Ossi Lisimore", telephone.Parameters.FunctionKeys[1].DisplayName, "display name should be changed")
			} else {
				assert.True(t, telephone.Parameters.FunctionKeys[1].IsEmpty(), "display name should not be changed in case of an error")
			}
		})
	}
}

func TestTelephone_backup(t *testing.T) {
	telephone := Telephone{Backup: []byte("this is my telephone backup")}
	request := httptest.NewRequest("GET", "/SaveAllSettings", strings.NewReader(""))
	recorder := httptest.NewRecorder()
	telephone.backup(recorder, request)
	status, data := getStatusAndData(recorder)
	assert.Equal(t, http.StatusOK, status, "status is wrong")
	assert.Equal(t, "this is my telephone backup", data, "data is wrong")
}

func TestTelephone_restore(t *testing.T) {
	telephone := Telephone{}
	require.Empty(t, telephone.Backup, "the telephone backup must be empty before the test")
	t.Run("wrong content-type", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/RestoreSettings", strings.NewReader("restored backup"))
		recorder := httptest.NewRecorder()
		telephone.restore(recorder, request)
		status, data := getStatusAndData(recorder)
		assert.Equal(t, http.StatusUnsupportedMediaType, status, "status code is wrong")
		assert.Equal(t, "Header \"Content-Type\" must begin with value \"multipart/form-data; boundary=\", but was \"\"", data, "message is wrong")
	})
	require.Empty(t, telephone.Backup, "the telephone backup must be empty before the test")
	payload := fmt.Sprintf(payloadTemplate, "BOUNDARY-42", "hooray, a backup", "BOUNDARY-42")
	t.Run("incompatible boundary", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/RestoreSettings", strings.NewReader(payload))
		request.Header.Set("Content-Type", "multipart/form-data; boundary=a")
		recorder := httptest.NewRecorder()
		telephone.restore(recorder, request)
		status, data := getStatusAndData(recorder)
		assert.Equal(t, http.StatusBadRequest, status, "status code is wrong")
		assert.Equal(t, "could not parse multipart-form", data, "message is wrong")
	})
	require.Empty(t, telephone.Backup, "the telephone backup must be empty before the test")
	t.Run("incompatible boundary", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/RestoreSettings", strings.NewReader(payload))
		request.Header.Set("Content-Type", "multipart/form-data; boundary=BOUNDARY-42")
		recorder := httptest.NewRecorder()
		telephone.restore(recorder, request)
		status, data := getStatusAndData(recorder)
		assert.Equal(t, http.StatusOK, status, "status code is wrong")
		assert.Empty(t, data, "message is wrong")
		assert.Equal(t, "hooray, a backup\n", string(telephone.Backup), "phone must have received the backup")
	})
}

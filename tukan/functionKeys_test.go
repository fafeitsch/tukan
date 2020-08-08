package tukan

import (
	"github.com/fafeitsch/Tukan/tukan/mock"
	"github.com/fafeitsch/Tukan/tukan/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPhone_DownloadParameters(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	keys := []params.FunctionKey{
		{},
		{DisplayName: "Ellen", PhoneNumber: "42"},
		{DisplayName: "Alan", PhoneNumber: "50"},
		{},
		{},
		{},
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	connector := Connector{Client: http.DefaultClient, UserName: username, Password: password}
	phone, err := connector.SingleConnect(server.URL)
	require.NoError(t, err, "no error expected")
	telephone.Parameters = params.Parameters{FunctionKeys: keys}
	t.Run("success", func(t *testing.T) {
		got, err := phone.DownloadParameters()
		require.NoError(t, err, "error should be nil")
		require.NotNil(t, got, "response object should not be nil")
		assert.Equal(t, 3, len(got.FunctionKeys), "there should be three function keys after pruning")
		assert.Equal(t, "", got.FunctionKeys[0].DisplayName, "display name of first entry should be empty")
		assert.Equal(t, "Ellen", got.FunctionKeys[1].DisplayName, "display name of second function key is wrong")
		assert.Equal(t, "Alan", got.FunctionKeys[2].DisplayName, "display name of third function key is wrong")
		assert.Equal(t, "42", got.FunctionKeys[1].PhoneNumber, "phone number of second function key is wrong")
		assert.Equal(t, "50", got.FunctionKeys[2].PhoneNumber, "phone number of third function key is wrong")
	})
	t.Run("error", func(t *testing.T) {
		phone.token = "wrong"
		got, err := phone.DownloadParameters()
		require.EqualError(t, err, "authentication error, status code: 401 with message \"401 Unauthorized\" and content \"Token not valid.\"", "error message wrong")
		require.Nil(t, got, "result should be nil in case of an error")
	})
}

func TestPhone_UploadParameters(t *testing.T) {
	handler, telephone := mock.CreatePhone(username, password)
	telephone.Parameters.FunctionKeys[0] = params.FunctionKey{DisplayName: "Joe"}
	telephone.Parameters.FunctionKeys[1] = params.FunctionKey{CallPickupCode: "***"}
	keys := []params.FunctionKey{
		{},
		{DisplayName: "Ellen", PhoneNumber: "42"},
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	connector := Connector{Client: http.DefaultClient, UserName: username, Password: password}
	phone, err := connector.SingleConnect(server.URL)
	require.NoError(t, err, "no error expected")
	t.Run("success", func(t *testing.T) {
		assert.Equal(t, "", telephone.Parameters.FunctionKeys[1].DisplayName, "Display name of first function key should be empty before")
		err = phone.UploadParameters(params.Parameters{FunctionKeys: keys})
		require.NoError(t, err, "no error is expected")
		assert.Equal(t, "Joe", telephone.Parameters.FunctionKeys[0].DisplayName, "Display name of first function key should not have been changed")
		assert.Equal(t, "Ellen", telephone.Parameters.FunctionKeys[1].DisplayName, "Display name of first function key is wrong")
		assert.Equal(t, "42", telephone.Parameters.FunctionKeys[1].PhoneNumber, "Phone number of first function key is wrong")
		assert.Equal(t, "***", telephone.Parameters.FunctionKeys[1].CallPickupCode, "CallPickupCode should not have been changed")
	})
}

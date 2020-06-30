package http

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

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

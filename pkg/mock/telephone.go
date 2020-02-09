package mock

import (
	"encoding/json"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/fafeitsch/Tukan/pkg/api"
	"net/http"
	"strings"
)

type Telephone struct {
	login    string
	password string
	token    *string
}

func (t *Telephone) attemptLogin(w http.ResponseWriter, r *http.Request) {
	if t.token != nil {
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprintf(w, "Login is already consumed, log out first")
		return
	}
	if fail, status, msg := t.preconditionsFail(r, "application/json", "POST"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
	}
	creds := api.Credentials{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}
	if decoder.More() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "request body contained more than one json object, which is not allowed")
		return
	}
	if creds.Password != t.password || creds.Login != t.login {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "provided credentials not valid")
		return
	}
	token := uniuri.NewLen(64)
	t.token = &token
	payload, _ := json.Marshal(api.TokenResponse{Token: token})
	_, _ = w.Write(payload)
}

func (t *Telephone) preconditionsFail(r *http.Request, contentType string, method string) (bool, int, string) {
	if r.Method != method {
		return true, http.StatusBadRequest, fmt.Sprintf("Unsupported method \"%s\", want method \"%s", r.Method, method)
	}
	if r.Header.Get("Content-Type") != contentType {
		return true, http.StatusBadRequest, fmt.Sprintf("Header \"Content-Type\" must have value \"%s\", but was \"%s\"", contentType, r.Header.Get("Content-Type"))
	}
	return false, http.StatusOK, ""
}

func (t *Telephone) preconditionsFailWithAuth(r *http.Request, contentType string, method string) (bool, int, string) {
	if fail, status, msg := t.preconditionsFail(r, contentType, method); fail {
		return fail, status, msg
	}
	auth := r.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	if t.token == nil || token != *t.token {
		return true, http.StatusUnauthorized, fmt.Sprintf("Token not valid")
	}
	return false, http.StatusOK, ""
}

func (t *Telephone) logout(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFailWithAuth(r, "", "POST"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	t.token = nil
	w.WriteHeader(http.StatusNoContent)
}

package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/fafeitsch/Tukan/pkg/api"
	"io"
	"log"
	"net/http"
	"strings"
)

type Telephone struct {
	login     string
	password  string
	token     *string
	phonebook string
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
		return
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
	token := uniuri.NewLen(1)
	t.token = &token
	payload, _ := json.Marshal(api.TokenResponse{Token: token})
	_, _ = w.Write(payload)
}

func (t *Telephone) preconditionsFail(r *http.Request, contentType string, method string) (bool, int, string) {
	if r.Method != method {
		return true, http.StatusBadRequest, fmt.Sprintf("Unsupported method \"%s\", want method \"%s", r.Method, method)
	}
	declaredType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(declaredType, contentType) {
		return true, http.StatusBadRequest, fmt.Sprintf("Header \"Content-Type\" must begin with value \"%s\", but was \"%s\"", contentType, r.Header.Get("Content-Type"))
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

func (t *Telephone) postPhoneBook(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFailWithAuth(r, "multipart/form-data; boundary=", "POST"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	err := r.ParseMultipartForm(2 << 20) // 20 MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "could not parse multipart-form")
		log.Printf("Error while parsing multipart form from %s: %v", r.RemoteAddr, err)
		return
	}
	file, _, _ := r.FormFile("file")
	defer func() { _ = file.Close() }()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, file)
	t.phonebook = buf.String()
	w.WriteHeader(http.StatusNoContent)
}

func (t *Telephone) saveLocalPhoneBook(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFailWithAuth(r, "", "GET"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	r.Header.Set("Content-Type", "application/xml")
	_, _ = fmt.Fprintf(w, t.phonebook)
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

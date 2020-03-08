package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/fafeitsch/Tukan/pkg/api"
	"github.com/fafeitsch/Tukan/pkg/domain"
	"io"
	"log"
	"net/http"
	"strings"
)

type Telephone struct {
	Login      string
	Password   string
	Token      *string
	Phonebook  string
	Keys       api.FunctionKeys
	Parameters domain.Parameters
}

func (t *Telephone) AttemptLogin(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFail(r, "application/json", "POST"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	creds := api.Credentials{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}
	if decoder.More() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = fmt.Fprintf(w, "request body contained more than one json object, which is not allowed")
		return
	}
	if creds.Password != t.Password || creds.Login != t.Login {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "provided credentials not valid")
		return
	}
	token := uniuri.NewLen(32)
	t.Token = &token
	payload, _ := json.Marshal(api.TokenResponse{Token: token})
	_, _ = w.Write(payload)
}

func (t *Telephone) preconditionsFail(r *http.Request, contentType string, method string) (bool, int, string) {
	if r.Method != method {
		return true, http.StatusMethodNotAllowed, fmt.Sprintf("Unsupported method \"%s\", want method \"%s\"", r.Method, method)
	}
	declaredType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(declaredType, contentType) {
		return true, http.StatusUnsupportedMediaType, fmt.Sprintf("Header \"Content-Type\" must begin with value \"%s\", but was \"%s\"", contentType, r.Header.Get("Content-Type"))
	}
	return false, http.StatusOK, ""
}

func (t *Telephone) PostPhoneBook(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFail(r, "multipart/form-data; boundary=", "POST"); fail {
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
	t.Phonebook = buf.String()
	log.Printf("Saved phone book from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusNoContent)
}

func (t *Telephone) SaveLocalPhoneBook(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFail(r, "", "GET"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	r.Header.Set("Content-Type", "application/xml")
	_, _ = fmt.Fprintf(w, t.Phonebook)
}

func (t *Telephone) handleParameters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.getParameters(w, r)
		break
	case http.MethodPost:
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			_, _ = fmt.Fprintf(w, "contentType %s not supported", r.Header.Get("contentType"))
		}
		t.changeFunctionKeys(w, r.Body)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = fmt.Fprintf(w, "the method %s is not allowed. Only GET and POST are supported.", r.Method)
		return
	}
}

func (t *Telephone) changeFunctionKeys(w http.ResponseWriter, body io.ReadCloser) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&(t.Keys))
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
	log.Printf("Received function keys from: %v", t.Keys)
	w.WriteHeader(http.StatusNoContent)
}

func (t *Telephone) getParameters(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFail(r, "application/json", "GET"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	parameters := domain.Parameters{}
	keys := make([]domain.FunctionKey, 0, len(t.Keys.FunctionKeys))
	for _, key := range t.Keys.FunctionKeys {
		number := domain.Setting{Value: key.PhoneNumber}
		display := domain.Setting{Value: key.DisplayName}
		callpickup := domain.Setting{Value: key.CallPickupCode}
		domKey := domain.FunctionKey{PhoneNumber: number, DisplayName: display, CallPickupCode: callpickup}
		keys = append(keys, domKey)
	}
	parameters.FunctionKeys = keys
	payload, _ := json.Marshal(parameters)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func (t *Telephone) logout(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFail(r, "", "POST"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	t.Token = nil
	w.WriteHeader(http.StatusNoContent)
}

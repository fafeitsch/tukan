package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/fafeitsch/Tukan/pkg/api/down"
	"github.com/fafeitsch/Tukan/pkg/api/up"
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
	Parameters up.Parameters
}

func (t *Telephone) AttemptLogin(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFail(r, "application/json", "POST"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	creds := up.Credentials{}
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
	payload, _ := json.Marshal(down.TokenResponse{Token: token})
	w.Header().Add("Content-Type", "application/json")
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

func (t *Telephone) HandleParameters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.getParameters(w)
		break
	case http.MethodPost:
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			_, _ = fmt.Fprintf(w, "contentType \"%s\" not supported", r.Header.Get("Content-Type"))
			return
		}
		t.changeFunctionKeys(w, r.Body)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = fmt.Fprintf(w, "The method \"%s\" is not allowed. Only \"GET\" and \"POST\" are supported", r.Method)
		return
	}
}

func (t *Telephone) changeFunctionKeys(w http.ResponseWriter, body io.ReadCloser) {
	decoder := json.NewDecoder(body)
	keys := up.Parameters{}
	err := decoder.Decode(&(keys))
	if err != nil {
		http.Error(w, fmt.Sprintf("could not deserialize json: %v", err), http.StatusBadRequest)
		return
	}
	if decoder.More() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "request body contained more than one json object, which is not allowed")
		return
	}
	if len(t.Parameters.FunctionKeys) < len(keys.FunctionKeys) {
		msg := fmt.Sprintf("the request contained %d function keys, but the phone only has %d", len(keys.FunctionKeys), len(t.Parameters.FunctionKeys))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	for index, key := range keys.FunctionKeys {
		if t.Parameters.FunctionKeys[index] == nil {
			t.Parameters.FunctionKeys[index] = make(map[string]string)
		}
		for propertyName, value := range key {
			t.Parameters.FunctionKeys[index][propertyName] = value
		}
	}
	log.Printf("Received function keys")
	w.WriteHeader(http.StatusNoContent)
}

func (t *Telephone) getParameters(w http.ResponseWriter) {
	parameters := down.Parameters{}
	keys := make([]down.FunctionKey, 0, len(t.Parameters.FunctionKeys))
	for _, key := range t.Parameters.FunctionKeys {
		number := down.Setting{Value: key["PhoneNumber"], Options: []interface{}{}}
		display := down.Setting{Value: key["DisplayName"], Options: []interface{}{}}
		callpickup := down.Setting{Value: key["CallPickupCode"], Options: []interface{}{}}
		keyType := down.Setting{Value: down.KeyTypeBLF, Options: []interface{}{0, 1, 2, 3, down.KeyTypeBLF, 5, 6, 7}}
		if number.Value == "" && display.Value == "" {
			keyType.Value = down.KeyTypeNone
		}
		domKey := down.FunctionKey{DisplayName: display, PhoneNumber: number, CallPickupCode: callpickup, Type: keyType}
		keys = append(keys, domKey)
	}
	parameters.FunctionKeys = keys
	payload, _ := json.MarshalIndent(parameters, "", "  ")
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

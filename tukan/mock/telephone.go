package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/Tukan/tukan/params"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const KeyTypeBLF = "4"
const KeyTypeNone = "-1"

// A mock telephone for using in test environments. It has similar properties
// to those of the IP620/630 telephones and can be manipulated directly.
type Telephone struct {
	Login      string
	Password   string
	Token      *string
	Phonebook  string
	Parameters params.Parameters
}

func (t *Telephone) attemptLogin(w http.ResponseWriter, r *http.Request) {
	if fail, status, msg := t.preconditionsFail(r, "application/json", "POST"); fail {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, msg)
		return
	}
	creds := params.Credentials{}
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
	token := generateToken()
	t.Token = &token
	tokenObject := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	payload, _ := json.Marshal(tokenObject)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(payload)
}

func generateToken() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzåäö")
	length := 32
	var b strings.Builder
	for i := 0; i < length; i++ {
		// don't need crypt.rand here because it's a mock server for testing purposes
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
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

func (t *Telephone) postPhoneBook(w http.ResponseWriter, r *http.Request) {
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

func (t *Telephone) saveLocalPhoneBook(w http.ResponseWriter, r *http.Request) {
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
	keys := params.Parameters{}
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
	// The actual phones rather do a merge here (only overwrite values that exist in the request body),
	// but we want to keep it simple here.
	t.Parameters = keys
	log.Printf("Received function keys")
	w.WriteHeader(http.StatusNoContent)
}

func (t *Telephone) backup(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = fmt.Fprintf(w, "This is the backup")
	w.WriteHeader(http.StatusOK)
}

func (t *Telephone) restore(w http.ResponseWriter, req *http.Request) {
	content, _ := ioutil.ReadAll(req.Body)
	log.Printf("Content-Type: %s\n%s", req.Header.Get("Content-Type"), string(content))
	w.WriteHeader(http.StatusOK)
}

func (t *Telephone) getParameters(w http.ResponseWriter) {
	// Actually, this is not what the real telephones do!
	// Instead for every string setting, the telephones send a whole JSON object
	// {"value": "8080", "flags": 8, "validator": … , …}
	// However, this mock phone just sends the string value (same format as the POST method expects)
	// This is due to technical reasons because I don't want to copy the whole Parameters struct just
	// to define another Marshall method for the settings.
	payload, _ := json.MarshalIndent(t.Parameters, "", "  ")
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

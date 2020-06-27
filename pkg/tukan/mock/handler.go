package mock

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// Creates a new mocking phone as well as a http handler for the mock phone.
// The username and password parameters determine the credentials for the mock phone.
// Both entities can be used to set up a test environment for client functions.
// The mock telephone can be modified directly; and the handler can be given,
// for example, to an httptest server, which interacts with the telephone.
func CreatePhone(login string, password string) (http.Handler, *Telephone) {
	router := mux.NewRouter()
	tele := Telephone{
		Login:      login,
		Password:   password,
		Parameters: RawParameters{FunctionKeys: make([]map[string]string, 8)},
	}
	router.HandleFunc("/Login", tele.attemptLogin)
	router.Handle("/Logout", enforceTokenHandler(&tele, tele.logout))
	router.Handle("/LocalPhonebook", enforceTokenHandler(&tele, tele.postPhoneBook))
	router.Handle("/SaveLocalPhonebook", enforceTokenHandler(&tele, tele.saveLocalPhoneBook))
	router.Handle("/Parameters", enforceTokenHandler(&tele, tele.handleParameters))
	return router, &tele
}

func enforceTokenHandler(telephone *Telephone, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		if telephone.Token == nil || token != *telephone.Token {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprintf(w, "Token not valid.")
			return
		}
		next.ServeHTTP(w, r)
	})
}
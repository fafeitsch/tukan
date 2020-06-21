package mock

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api/up"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func CreatePhone(login string, password string) (http.Handler, *Telephone) {
	router := mux.NewRouter()
	tele := Telephone{
		Login:      login,
		Password:   password,
		Parameters: up.Parameters{FunctionKeys: make([]map[string]string, 8)},
	}
	router.HandleFunc("/Login", tele.AttemptLogin)
	router.Handle("/Logout", enforceTokenHandler(&tele, tele.logout))
	router.Handle("/LocalPhonebook", enforceTokenHandler(&tele, tele.PostPhoneBook))
	router.Handle("/SaveLocalPhonebook", enforceTokenHandler(&tele, tele.SaveLocalPhoneBook))
	router.Handle("/Parameters", enforceTokenHandler(&tele, tele.HandleParameters))
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

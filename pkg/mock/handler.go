package mock

import (
	"fmt"
	"github.com/fafeitsch/Tukan/pkg/api/up"
	"log"
	"net/http"
	"strings"
)

func StartHandler(port int, login string, password string) {
	tele := Telephone{
		Login:      login,
		Password:   password,
		Parameters: up.Parameters{FunctionKeys: make([]map[string]string, 8)},
	}
	http.HandleFunc("/Login", tele.AttemptLogin)
	http.Handle("/Logout", enforceTokenHandler(&tele, tele.logout))
	http.Handle("/LocalPhonebook", enforceTokenHandler(&tele, tele.PostPhoneBook))
	http.Handle("/SaveLocalPhonebook", enforceTokenHandler(&tele, tele.SaveLocalPhoneBook))
	http.Handle("/Parameters", enforceTokenHandler(&tele, tele.HandleParameters))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
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

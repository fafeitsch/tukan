package mock

import (
	"fmt"
	"log"
	"net/http"
)

func StartHandler(port int, login string, password string) {
	tele := Telephone{
		login:    login,
		password: password,
	}
	http.HandleFunc("/Login", tele.attemptLogin)
	http.HandleFunc("/Logout", tele.logout)
	http.HandleFunc("/LocalPhonebook", tele.postPhoneBook)
	http.HandleFunc("/SaveLocalPhonebook", tele.saveLocalPhoneBook)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

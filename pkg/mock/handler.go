package mock

import (
	"fmt"
	"log"
	"net/http"
)

func StartHandler(port int, login string, password string) {
	tele := Telephone{
		Login:    login,
		Password: password,
	}
	http.HandleFunc("/Login", tele.AttemptLogin)
	http.HandleFunc("/Logout", tele.logout)
	http.HandleFunc("/LocalPhonebook", tele.postPhoneBook)
	http.HandleFunc("/SaveLocalPhonebook", tele.saveLocalPhoneBook)
	http.HandleFunc("/Parameters", tele.changeFunctionKeys)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

package mock

import (
	"encoding/csv"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
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

func ParseFunctionKeysCsv(filename string) ([]map[string]string, error) {
	reader, err := os.Open(filename)
	defer func() { _ = reader.Close() }()
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(reader)
	all, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not parse csv content: %v", err)
	}
	result := make([]map[string]string, 0, len(all))
	for index, row := range all[1:] {
		if len(row) < 3 {
			return nil, fmt.Errorf("row %d has %d values, but exactly 3 are required", index, len(row))
		}
		result = append(result, map[string]string{"DisplayName": row[0], "PhoneNumber": row[1], "CallPickupCode": row[2]})
	}
	return result, nil
}

package up

import "fmt"

type Credentials struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type FunctionKey struct {
	DisplayName    string `json:"DisplayName,omitEmpty"`
	CallPickupCode string `json:"CallPickupCode,omitEmpty"`
	PhoneNumber    string `json:"PhoneNumber,omitEmpty"`
}

func (f *FunctionKey) String() string {
	return fmt.Sprintf("DisplayName: %s, Pickup: %s, Number: %s", f.DisplayName, f.CallPickupCode, f.PhoneNumber)
}

type FunctionKeys struct {
	FunctionKeys []map[string]string `json:"FunctionKeys"`
}

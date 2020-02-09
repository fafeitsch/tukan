package api

import "fmt"

type Credentials struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type FunctionKey struct {
	DisplayName    string `json:"DisplayName"`
	CallPickupCode string `json:"CallPickupCode"`
	PhoneNumber    string `json:"PhoneNumber"`
}

func (f *FunctionKey) String() string {
	return fmt.Sprintf("DisplayName: %s, Pickup: %s, Number: %s", f.DisplayName, f.CallPickupCode, f.PhoneNumber)
}

type FunctionKeys struct {
	FunctionKeys []FunctionKey `json:"FunctionKeys"`
}

package up

import "fmt"

type Credentials struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type FunctionKey struct {
	DisplayName    string `json:"DisplayName,omitempty"`
	CallPickupCode string `json:"CallPickupCode,omitempty"`
	PhoneNumber    string `json:"PhoneNumber,omitempty"`
}

func (f *FunctionKey) String() string {
	return fmt.Sprintf("DisplayName: %s, Pickup: %s, Number: %s", f.DisplayName, f.CallPickupCode, f.PhoneNumber)
}

type Parameters struct {
	FunctionKeys []FunctionKey `json:"FunctionKeys"`
}

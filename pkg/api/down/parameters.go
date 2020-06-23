package down

import (
	"fmt"
	"strings"
)

type Parameters struct {
	FunctionKeys FunctionKeys `json:"FunctionKeys"`
}

type FunctionKeys []FunctionKey

func (f FunctionKeys) String() string {
	keys := make([]string, 0, len(f))
	for _, key := range f {
		str := key.String()
		keys = append(keys, str)
	}
	return strings.Join(keys, ", ")
}

type FunctionKey struct {
	DisplayName    Setting `json:"DisplayName"`
	PhoneNumber    Setting `json:"PhoneNumber"`
	CallPickupCode Setting `json:"CallPickupCode"`
	Type           Setting `json:"Type"`
}

func (f *FunctionKey) IsEmpty() bool {
	return f.Type.Value == "-1"
}

func (f *FunctionKey) String() string {
	var typ string
	switch f.Type.Value {
	case KeyTypeBLF:
		typ = "BLF"
		break
	default:
		typ = "unknown"
	}
	return fmt.Sprintf("[\"%s\": %s (%s) (%s)]", f.DisplayName.Value, f.PhoneNumber.Value, f.CallPickupCode.Value, typ)
}

type Setting struct {
	Value     string      `json:"value"`
	Flags     int         `json:"flags"`
	Validator Validator   `json:"validator"`
	Options   interface{} `json:"options"`
}

type Validator struct {
	Regexp string `json:"regexp"`
}

const KeyTypeBLF = "4"
const KeyTypeNone = "-1"

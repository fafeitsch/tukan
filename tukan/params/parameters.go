package params

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
)

type Credentials struct {
	Login    string
	Password string
}

// Parameters describes all the settings of the VoIP phone. The phones have
// slightly different format for the parameter upload/download. While unmarshalling
// downloaded parameters, both formats can be unmarshalled.
// For marshalling, the upload format is used.
type Parameters struct {
	FunctionKeys FunctionKeys `json:"FunctionKeys"`
}

func (p *Parameters) TransformFunctionKeyNames(original, replace string) (Parameters, []int) {
	keys := make([]FunctionKey, 0, len(p.FunctionKeys))
	changed := make([]int, 0, 0)
	for index, fnKey := range p.FunctionKeys {
		var key = FunctionKey{}
		if fnKey.DisplayName.String() == original {
			key = FunctionKey{DisplayName: Setting(replace)}
			changed = append(changed, index)
		}
		keys = append(keys, key)
	}
	return Parameters{FunctionKeys: keys}, changed
}

type FunctionKeys []FunctionKey

func (f FunctionKeys) WriteCsvWithHeader(writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	_ = csvWriter.Write([]string{"DisplayName", "PhoneNumber", "CallPickupCode", "Type"})
	for _, fnKey := range f {
		_ = csvWriter.Write([]string{fnKey.DisplayName.String(), fnKey.PhoneNumber.String(), fnKey.CallPickupCode.String(), fnKey.Type.String()})
	}
	csvWriter.Flush()
	if csvWriter.Error() != nil {
		return fmt.Errorf("could not write: %v", csvWriter.Error())
	}
	return nil
}

type FunctionKey struct {
	DisplayName    Setting `json:"DisplayName"`
	PhoneNumber    Setting `json:"PhoneNumber"`
	CallPickupCode Setting `json:"CallPickupCode"`
	Type           Setting `json:"Type"`
}

func (f *FunctionKey) IsEmpty() bool {
	return (f.Type == "" && f.PhoneNumber == "" && f.DisplayName == "" && f.CallPickupCode == "") || f.Type == "-1"
}

func (f *FunctionKey) Merge(other FunctionKey) {
	if other.PhoneNumber != "" {
		f.PhoneNumber = other.PhoneNumber
	}
	if other.CallPickupCode != "" {
		f.CallPickupCode = other.PhoneNumber
	}
	if other.Type != "" {
		f.Type = other.Type
	}
	if other.DisplayName != "" {
		f.DisplayName = other.DisplayName
	}
}

type Setting string

func (s *Setting) String() string {
	return string(*s)
}

func (s *Setting) UnmarshalJSON(data []byte) error {
	setting := struct {
		Value string `json:"value"`
	}{}
	err := json.Unmarshal(data, &setting)
	// data is not download-format, try the upload format:
	if err != nil {
		str := ""
		err = json.Unmarshal(data, &str)
		if err != nil {
			return err
		}
		*s = Setting(str)
		return nil
	}
	got := Setting(setting.Value)
	*s = got
	return nil
}

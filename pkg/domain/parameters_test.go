package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParameters_PurgeTrailingFunctionKeys(t *testing.T) {
	fk1 := FunctionKey{Type: Setting{Value: "4"}}
	fk2 := FunctionKey{Type: Setting{Value: "4"}}
	fk3 := FunctionKey{Type: Setting{Value: "-1"}}
	fk4 := FunctionKey{Type: Setting{Value: "4"}}
	fk5 := FunctionKey{Type: Setting{Value: "-1"}}
	fk6 := FunctionKey{Type: Setting{Value: "-1"}}
	keys := []FunctionKey{fk1, fk2, fk3, fk4, fk5, fk6}
	cases := []struct {
		name       string
		input      []FunctionKey
		wantOutput []FunctionKey
	}{
		{name: "delete nothing", input: keys[0:4], wantOutput: keys[0:4]},
		{name: "delete last two", input: keys, wantOutput: keys[0:4]},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			params := Parameters{FunctionKeys: tt.input}
			params.PurgeTrailingFunctionKeys()
			assert.Equal(t, tt.wantOutput, params.FunctionKeys, "function keys after removal of empty trailing keys not correct.")
		})
	}
}

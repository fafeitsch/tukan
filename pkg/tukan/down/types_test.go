package down

import (
	"github.com/fafeitsch/Tukan/pkg/tukan/up"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParameters_TransformFunctionKeyNames(t *testing.T) {
	parameters := Parameters{FunctionKeys: FunctionKeys{
		FunctionKey{PhoneNumber: Setting{Value: "20"}, DisplayName: Setting{Value: "Ron Doe"}},
		FunctionKey{PhoneNumber: Setting{Value: "30"}, DisplayName: Setting{Value: "Ronald Gene"}},
		FunctionKey{PhoneNumber: Setting{Value: "40"}, DisplayName: Setting{Value: "Mary Hope"}},
		FunctionKey{PhoneNumber: Setting{Value: "42"}, DisplayName: Setting{Value: "Ronald Gene"}},
	}}
	got, ints := parameters.TransformFunctionKeyNames("Ronald Gene", "Belinda Fox")
	assert.Equal(t, []int{1, 3}, ints, "changed indices are wrong")
	assert.Equal(t, up.FunctionKey{}, got.FunctionKeys[0], "first function key wrong")
	assert.Equal(t, up.FunctionKey{DisplayName: "Belinda Fox"}, got.FunctionKeys[1], "second function key wrong")
	assert.Equal(t, up.FunctionKey{}, got.FunctionKeys[2], "forth function key wrong")
	assert.Equal(t, up.FunctionKey{DisplayName: "Belinda Fox"}, got.FunctionKeys[3], "fifth function key wrong")
}

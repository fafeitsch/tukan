package params

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParameters_TransformFunctionKeyNames(t *testing.T) {
	parameters := Parameters{FunctionKeys: FunctionKeys{
		FunctionKey{PhoneNumber: "20", DisplayName: "Ron Doe"},
		FunctionKey{PhoneNumber: "30", DisplayName: "Ronald Gene"},
		FunctionKey{PhoneNumber: "40", DisplayName: "Mary Hope"},
		FunctionKey{PhoneNumber: "42", DisplayName: "Ronald Gene"},
	}}
	got, ints := parameters.TransformFunctionKeyNames("Ronald Gene", "Belinda Fox")
	assert.Equal(t, "Ronald Gene", parameters.FunctionKeys[1].DisplayName, "original struct must not be changed")
	assert.Equal(t, "Ronald Gene", parameters.FunctionKeys[3].DisplayName, "original struct must not be changed")
	assert.Equal(t, []int{1, 3}, ints, "changed indices are wrong")
	assert.Equal(t, FunctionKey{}, got.FunctionKeys[0], "first function key wrong")
	assert.Equal(t, FunctionKey{DisplayName: "Belinda Fox"}, got.FunctionKeys[1], "second function key wrong")
	assert.Equal(t, FunctionKey{}, got.FunctionKeys[2], "forth function key wrong")
	assert.Equal(t, FunctionKey{DisplayName: "Belinda Fox"}, got.FunctionKeys[3], "fifth function key wrong")
}

package params

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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

func TestFunctionKey_IsEmpty(t *testing.T) {
	t.Run("type is -1", func(t *testing.T) {
		key := &FunctionKey{Type: "-1", DisplayName: "John"}
		assert.True(t, key.IsEmpty(), "with type = -1 it should always be empty")
	})
	t.Run("most properties are empty", func(t *testing.T) {
		key := &FunctionKey{}
		assert.True(t, key.IsEmpty(), "if all important fields are empty, the function key should be considered empty")
	})
	t.Run("not empty", func(t *testing.T) {
		key := &FunctionKey{Type: "4", DisplayName: "John"}
		assert.False(t, key.IsEmpty(), "key should not be empty")
	})
}

func TestFunctionKey_Merge(t *testing.T) {
	t.Run("merge phone number (empty)", func(t *testing.T) {
		key := &FunctionKey{}
		key2 := &FunctionKey{PhoneNumber: "555-Nose"}
		key.Merge(*key2)
		assert.Equal(t, key2.PhoneNumber, key.PhoneNumber, "if field is empty, merge")
	})
	t.Run("merge phone number", func(t *testing.T) {
		key := &FunctionKey{PhoneNumber: "123-Nothing"}
		key2 := &FunctionKey{PhoneNumber: ""}
		key.Merge(*key2)
		assert.Equal(t, "123-Nothing", key.PhoneNumber, "if field is empty, merge")
	})
}

func TestParameters_UnmarshalJSON(t *testing.T) {
	file, err := ioutil.ReadFile("../mock/mockdata/parameters.json")
	require.NoError(t, err, "no error expected")
	parameters := Parameters{}
	err = json.Unmarshal(file, &parameters)
	require.NoError(t, err, "no error expected")
	assert.Equal(t, "Donald HO", parameters.FunctionKeys[4].DisplayName)
	assert.Equal(t, "example.com/pbx", parameters.Sip[1].Domain)
	assert.Equal(t, "24h", parameters.TimeFormat)
}

package params

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestFunctionKeys_Transform(t *testing.T) {
	keys := FunctionKeys{
		FunctionKey{PhoneNumber: "20", DisplayName: "Ron Doe"},
		FunctionKey{PhoneNumber: "30", DisplayName: "Ronald Gene"},
		FunctionKey{PhoneNumber: "40", DisplayName: "Mary Hope"},
		FunctionKey{PhoneNumber: "42", DisplayName: "Ronald Gene"},
	}
	got, ints := keys.Transform(ReplaceDisplayName("Ronald Gene", "Belinda Fox"))
	assert.Equal(t, []int{1, 3}, ints, "changed indices are wrong")
	assert.Equal(t, keys[0], got[0], "first key must not be changed")
	assert.Equal(t, "30", got[1].PhoneNumber, "phone number of second function key wrong")
	assert.Equal(t, "Belinda Fox", got[1].DisplayName, "display name of second function key wrong")
	assert.Equal(t, "Ronald Gene", keys[1].DisplayName, "original display name must not be changed")
	assert.Equal(t, keys[2], got[2], "second key must not be changed")
	assert.Equal(t, "42", got[3].PhoneNumber, "phone number of forth function key wrong")
	assert.Equal(t, "Belinda Fox", got[3].DisplayName, "display name of forth function key wrong")
	assert.Equal(t, "Ronald Gene", keys[3].DisplayName, "original display name must not be changed")
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

func TestSips_Transform(t *testing.T) {
	sips := Sips{
		{DisplayName: "John"},
		{DisplayName: ""},
		{DisplayName: "Freyan"},
	}
	got, ints := sips.Transform(SipOverrideDisplayName("222 John"))
	assert.Equal(t, []int{0, 2}, ints, "changed indices are not correct")
	assert.Equal(t, "222 John", got[0].DisplayName)
	assert.Equal(t, "", got[1].DisplayName)
	assert.Equal(t, "222 John", got[2].DisplayName)
}

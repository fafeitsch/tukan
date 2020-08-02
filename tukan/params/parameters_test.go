package params

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, "Ronald Gene", parameters.FunctionKeys[1].DisplayName.String(), "original struct must not be changed")
	assert.Equal(t, "Ronald Gene", parameters.FunctionKeys[3].DisplayName.String(), "original struct must not be changed")
	assert.Equal(t, []int{1, 3}, ints, "changed indices are wrong")
	assert.Equal(t, FunctionKey{}, got.FunctionKeys[0], "first function key wrong")
	assert.Equal(t, FunctionKey{DisplayName: "Belinda Fox"}, got.FunctionKeys[1], "second function key wrong")
	assert.Equal(t, FunctionKey{}, got.FunctionKeys[2], "forth function key wrong")
	assert.Equal(t, FunctionKey{DisplayName: "Belinda Fox"}, got.FunctionKeys[3], "fifth function key wrong")
}

type mockWriter struct {
	do func(p []byte) (int, error)
}

func (m *mockWriter) Write(p []byte) (int, error) {
	return m.do(p)
}

func TestFunctionKeys_WriteCsvWithHeader(t *testing.T) {
	keys := FunctionKeys{
		{DisplayName: "Clarice Lecter", PhoneNumber: "44", CallPickupCode: "*0", Type: "BLF"},
		{DisplayName: "Hannibal Starling", PhoneNumber: "12", CallPickupCode: "#0", Type: "QuickDial"},
	}
	t.Run("success", func(t *testing.T) {
		buf := &bytes.Buffer{}
		err := keys.WriteCsvWithHeader(buf)
		require.NoError(t, err, "no error expected")
		got := buf.String()
		assert.Equal(t, "DisplayName,PhoneNumber,CallPickupCode,Type\nClarice Lecter,44,*0,BLF\nHannibal Starling,12,#0,QuickDial\n", got, "generated csv is wrong")
	})
	t.Run("failure 1", func(t *testing.T) {
		writer := mockWriter{do: func(p []byte) (int, error) {
			return 0, fmt.Errorf("forbidden to write the header")
		}}
		err := keys.WriteCsvWithHeader(&writer)
		assert.EqualError(t, err, "could not write: forbidden to write the header", "error message wrong")
	})
}

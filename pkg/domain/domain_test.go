package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTukanResult_String(t *testing.T) {
	tukan := map[string]string{
		"10.32.12.4":    "success",
		"110.24.25.2":   "error - could not login",
		"telephone.org": "404, not found",
	}
	got := TukanResult(tukan)
	want := `
==========
 Results
==========
10.32.12.4: success
110.24.25.2: error - could not login
telephone.org: 404, not found
`
	assert.Equal(t, want, got.String(), "string conversion not correct")
}

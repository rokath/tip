package tiptable

import (
	"testing"

	"github.com/tj/assert"
)

func Test_spaces(t *testing.T) {
	tt := []struct {
		l int    // given count
		s string // expected string
	}{
		{-2, ""},
		{-1, ""},
		{0, ""},
		{1, " "},
		{2, "  "},
	}
	for _, x := range tt {
		act := spaces(x.l)
		assert.Equal(t, x.s, act)
	}
}

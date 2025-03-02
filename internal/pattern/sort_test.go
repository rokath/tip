package pattern

import (
	"sync"
	"testing"

	"github.com/tj/assert"
)

func TestHistogram_SortKeysByIncrSize(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		exp    fields
	}{
		// test cases:
		{
			"", // name
			fields{map[string]Pat{}, &mu, []string{"bb11", "112233", "aa22"}},
			fields{map[string]Pat{}, &mu, []string{"aa22", "bb11", "112233"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.SortKeysByIncrSize()
			for i := range p.Key {
				assert.Equal(t, tt.exp.Keys[i], tt.fields.Keys[i])
			}
		})
	}
}

func TestSortByDescCountDescLength(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	pat := []Patt{
		{100, []byte{1, 2, 3, 1, 2, 3, 4}, "01020301020304"},
		{100, []byte{1, 2, 3, 4}, "01020304"},
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
		{900, []byte{1, 2}, "0102"},
		{100, []byte{8, 2, 3, 1, 2, 3}, "080203010203"},
		{300, []byte{1, 2, 3}, "010203"},
	}
	exp := []Patt{
		{900, []byte{1, 2}, "0102"},
		{300, []byte{1, 2, 3}, "010203"},
		{100, []byte{1, 2, 3, 1, 2, 3, 4}, "01020301020304"},
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
		{100, []byte{8, 2, 3, 1, 2, 3}, "080203010203"},
		{100, []byte{1, 2, 3, 4}, "01020304"},
	}
	act := SortByDescCountDescLength(pat)
	assert.Equal(t, exp, act)
}

// generated: ////////////////////////////////

package pattern

import (
	"sync"
	"testing"

	"github.com/tj/assert"
)

func TestHistogram_SortKeysByIncrSize(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]Pattern
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
			fields{map[string]Pattern{}, &mu, []string{"bb11", "112233", "aa22"}},
			fields{map[string]Pattern{}, &mu, []string{"aa22", "bb11", "112233"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Keys: tt.fields.Keys,
			}
			p.SortKeysByIncrSize()
			for i := range p.Keys {
				assert.Equal(t, tt.exp.Keys[i], tt.fields.Keys[i])
			}
		})
	}
}

func TestSortByDescCountDescLength(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	pat := []Pattern{
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10}},
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10, 20}},
	}
	exp := []Pattern{
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10, 20}},
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10}},
	}
	act := SortByDescCount(pat)
	assert.Equal(t, exp, act)
}

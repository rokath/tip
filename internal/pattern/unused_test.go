package pattern

import (
	"sync"
	"testing"

	"github.com/tj/assert"
)

func TestHistogram_SortKeysByDescSize(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]int
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		exp    fields
	}{
		// TODO: Add test cases.
		{
			"", // name
			fields{map[string]int{}, &mu, []string{"bb11", "112233", "aa22"}},
			fields{map[string]int{}, &mu, []string{"112233", "aa22", "bb11"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.SortKeysByDescSize()
			for i := range p.Key {
				assert.Equal(t, tt.exp.Keys[i], tt.fields.Keys[i])
			}
		})
	}
}

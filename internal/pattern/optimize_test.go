package pattern

import (
	"sync"
	"testing"

	"github.com/tj/assert"
)

func Test_countOverlapping2(t *testing.T) {
	type args struct {
		s   string
		sub string
	}
	tests := []struct {
		name string
		args args
		want int
	}{ // test cases:
		{"", args{"11111111", "1111111111"}, 0},
		{"", args{"11111111", "11111111"}, 1},
		{"", args{"11111111", "111111"}, 2},
		{"", args{"11111111", "1111"}, 3},
		{"", args{"11111111", "11"}, 4},
		{"", args{"1111", "111a"}, 0},
		{"", args{"1111", "1111"}, 1},
		{"", args{"111111", "1111"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countOverlapping2(tt.args.s, tt.args.sub); got != tt.want {
				t.Errorf("countOverlapping2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHistogram_ReduceOverlappingKeys(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Keys []string
	}
	type args struct {
		equalSize1stKey []string
		equalSize2ndKey []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    map[string]Pat
	}{ // test cases:
		{
			"",
			fields{map[string]Pat{"1111": {20, nil}, "11111111": {2, nil}}, &mu, nil},
			args{[]string{"11111111"}, []string{"1111"}},
			map[string]Pat{"1111": {14, nil}, "11111111": {2, nil}},
		},
		{
			"",
			fields{map[string]Pat{"1122": {10, nil}, "112233": {1, nil}}, &mu, nil},
			args{[]string{"112233"}, []string{"1122"}},
			map[string]Pat{"1122": {9, nil}, "112233": {1, nil}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.ReduceOverlappingKeys(tt.args.equalSize1stKey, tt.args.equalSize2ndKey)
			assert.Equal(t, tt.exp, p.Hist)
		})
	}
}

func TestHistogram_BalanceByteUsage(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Key  []string
	}
	type args struct {
		maxPatternSize int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    map[string]Pat
	}{ // test cases:
		{"", fields{map[string]Pat{"1122": {1, []int{0}}}, &m, nil}, args{8}, map[string]Pat{"1122": {500, []int{0}}}},
		{"", fields{map[string]Pat{"1122": {10, []int{0}}}, &m, nil}, args{8}, map[string]Pat{"1122": {5000, []int{0}}}},
		{"", fields{map[string]Pat{"112233": {10, []int{0}}}, &m, nil}, args{8}, map[string]Pat{"112233": {3333, []int{0}}}},
		{"", fields{map[string]Pat{"1111": {3, []int{0}}}, &m, nil}, args{8}, map[string]Pat{"1111": {1500, []int{0}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Key,
			}
			p.BalanceByteUsage(tt.args.maxPatternSize)
		})
	}
}

func TestHistogram_Reduce(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		exp    map[string]Pat
	}{ // test cases:
		{
			"",
			fields{map[string]Pat{"1122": {10, nil}, "112233": {1, nil}}, &mu, nil},
			map[string]Pat{"1122": {9, nil}, "112233": {1, nil}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.GetKeys()
			p.Reduce()
			assert.Equal(t, tt.exp, tt.fields.Hist)
		})
	}
}

// generated: ////////////////////////////////

package pattern

import (
	"fmt"
	"reflect"
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

func TestHistogram_DeletePosition(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Key  []string
	}
	type args struct {
		key      string
		position int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    fields
	}{
		// test cases:
		{"",
			fields{map[string]Pat{"1122": {3, []int{4, 5, 7}}, "112233": {1, nil}}, &mu, nil},
			args{"1122", 5},
			fields{map[string]Pat{"1122": {2, []int{4, 7}}, "112233": {1, nil}}, &mu, nil},
		},
		{"",
			fields{map[string]Pat{"1122": {1, []int{4}}, "112233": {1, nil}}, &mu, nil},
			args{"1122", 5},
			fields{map[string]Pat{"1122": {1, []int{4}}, "112233": {1, nil}}, &mu, nil},
		},
		{"",
			fields{map[string]Pat{"1122": {1, []int{4}}, "112233": {1, nil}}, &mu, nil},
			args{"1122", 4},
			fields{map[string]Pat{"1122": {0, []int{}}, "112233": {1, nil}}, &mu, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Key,
			}
			p.DeletePosition(tt.args.key, tt.args.position)
			e := tt.exp.Hist
			a := p.Hist
			fmt.Println("exp:", reflect.ValueOf(e).Type(), e)
			fmt.Println("act:", reflect.ValueOf(a).Type(), a)
			result := reflect.DeepEqual(e, a)
			assert.True(t, result)
		})
	}
}

func Test_positionIndexMatch(t *testing.T) {
	type args struct {
		a   []int
		idx int
		b   []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{ // test cases:
		{"", args{[]int{0, 30, 40}, 0, []int{22, 23, 0, 25}}, 0},
		{"", args{[]int{20, 30, 40}, 3, []int{22, 23, 24, 25}}, 23},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := positionIndexMatch(tt.args.a, tt.args.idx, tt.args.b); got != tt.want {
				t.Errorf("positionIndexMatch() = %v, want %v", got, tt.want)
			}
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
			fields{map[string]Pat{"1122": {2, []int{4, 9}}, "112233": {1, []int{4}}}, &mu, nil},
			map[string]Pat{"1122": {1, []int{9}}, "112233": {1, []int{4}}},
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
			fields{map[string]Pat{"1a1a": {3, []int{0, 1, 5}}, "1a1a1a": {1, []int{0}}}, &mu, nil}, // the histogram in p
			args{[]string{"1a1a1a"}, []string{"1a1a"}},                                          // the function arguments
			map[string]Pat{"1a1a": {1, []int{5}}, "1a1a1a": {1, []int{0}}},                      // the expected result in p !!!!!!!!!!!!!!!!!!
		},
		{
			"", // case: |xx1a1a1a1axx...|
			fields{map[string]Pat{"1a1a": {9, []int{0, 2, 3, 4, 5, 6, 8, 10, 20}}, "1a1a1a1a": {2, []int{2, 32}}}, &mu, nil}, // the histograms in p
			args{[]string{"1a1a1a1a"}, []string{"1a1a"}},                                                            // the function arguments
			map[string]Pat{"1a1a": {6, []int{0, 5, 6, 8, 10, 20}}, "1a1a1a1a": {2, []int{2, 32}}},                         // the expected result in p
		},
		{
			"",
			fields{map[string]Pat{"1122": {2, []int{8,32}}, "112233": {1, []int{8}}}, &mu, nil},
			args{[]string{"112233"}, []string{"1122"}},
			map[string]Pat{"1122": {1, []int{32}}, "112233": {1, []int{8}}},
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
			p.SortPositions()
			e := tt.exp
			a := p.Hist
			fmt.Println("exp:", reflect.ValueOf(e).Type(), e)
			fmt.Println("act:", reflect.ValueOf(a).Type(), a)
			result := reflect.DeepEqual(e, a)
			assert.True(t, result)

			assert.Equal(t, tt.exp, p.Hist)
		})
	}
}

// generated: ////////////////////////////////

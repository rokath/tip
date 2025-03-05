package pattern

import (
	"fmt"
	"sync"
	"testing"

	"github.com/tj/assert"
)

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
		{"", fields{map[string]Pat{s2h("a"): {4, []int{0, 1, 2, 3}}}, &m, nil}, args{4},
			/*   */ map[string]Pat{s2h("a"): {4, []int{0, 1, 2, 3}}}},

		{"", fields{map[string]Pat{s2h("aa"): {3, []int{0, 1, 2}}}, &m, nil}, args{4},
			/*   */ map[string]Pat{s2h("aa"): {4, []int{0, 1, 2}}}},

		{"", fields{map[string]Pat{s2h("aaa"): {2, []int{0, 1}}}, &m, nil}, args{4},
			/*   */ map[string]Pat{s2h("aaa"): {4, []int{0, 1}}}},

		{"", fields{map[string]Pat{s2h("aaaa"): {1, []int{0}}}, &m, nil}, args{4},
			/*   */ map[string]Pat{s2h("aaaa"): {4, []int{0}}}},

		{"", fields{map[string]Pat{s2h("ab"): {10, []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}}}, &m, nil}, args{8},
			/*   */ map[string]Pat{s2h("ab"): {float64(10*8) / 7, []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}}}},

		{"", fields{map[string]Pat{s2h("ab"): {3, []int{0, 3, 6}}}, &m, nil}, args{8},
			/*   */ map[string]Pat{s2h("ab"): {float64(3*8) / 7, []int{0, 3, 6}}}},

		{"", fields{map[string]Pat{s2h("ab"): {2, []int{0, 3}}, s2h("bc"): {1.0, []int{1}}, s2h("ca"): {1.0, []int{2}}, s2h("abc"): {1, []int{0}}, s2h("bca"): {1, []int{1}}, s2h("cab"): {1, []int{2}}, s2h("aba"): {1, []int{3}}, s2h("bab"): {1, []int{4}}}, &m, nil}, args{3},
			/*   */ map[string]Pat{s2h("ab"): {3, []int{0, 3}}, s2h("bc"): {1.5, []int{1}}, s2h("ca"): {1.5, []int{2}}, s2h("abc"): {3, []int{0}}, s2h("bca"): {3, []int{1}}, s2h("cab"): {3, []int{2}}, s2h("aba"): {3, []int{3}}, s2h("bab"): {3, []int{4}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Key,
			}
			p.BalanceByteUsage(tt.args.maxPatternSize)
			for k, v := range p.Hist {
				a := v.Weight
				e := tt.exp[k].Weight
				ok := withinTolerance(a, e, 0.001)
				if !ok {
					fmt.Println(k, e, a)
				}
				assert.True(t, ok)
			}
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

func TestHistogram_Reduce_0(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   int
		exp    map[string]Pat
	}{ // test cases:
		{
			"", // abcab: TODO Issue ca occurs 1 times but is inside bca and cab!
			//     01234
			fields{map[string]Pat{
				s2h("ab"):  {2, []int{0, 3}}, // in abc and cab and in (aba) and (bab) formally
				s2h("bc"):  {1, []int{1}},    // in abc and bca
				s2h("ca"):  {1, []int{2}},    // in bca and cab
				s2h("abc"): {1, []int{0}},
				s2h("bca"): {1, []int{1}},
				s2h("cab"): {1, []int{2}},
				s2h("aba"): {1, []int{3}},
				s2h("bab"): {1, []int{4}}}, &m, nil}, 3,
			map[string]Pat{
				s2h("abc"): {3, []int{0}},
				s2h("bca"): {3, []int{1}},
				s2h("cab"): {3, []int{2}},
				s2h("aba"): {3, []int{3}},
				s2h("bab"): {3, []int{4}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.BalanceByteUsage(tt.args)
			p.Reduce()
			p.DeleteEmptyKeys()
			assert.Equal(t, tt.exp, tt.fields.Hist)
		})
	}
}

func TestHistogram_Reduce_1_ok(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		data   []byte
		maxLen int
		ring   bool
		reduce bool
		exp    map[string]Pat
	}{ // test cases:
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("aaaaa"), 3, true, true,
			map[string]Pat{
				//s2h("aa"):  {0, []int{}},
				s2h("aaa"): {5, []int{0, 1, 2, 3, 4}},
			},
		},
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("aaaaa"), 3, true, false,
			map[string]Pat{
				s2h("aa"):  {5, []int{0, 1, 2, 3, 4}},
				s2h("aaa"): {5, []int{0, 1, 2, 3, 4}},
			},
		},
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("ABABCAB"), 3, true, true,
			map[string]Pat{
				s2h("ABA"): {2, []int{0, 5}}, // "414241":pattern.Pat{Weight:2, Pos:[]int{0, 5}},
				s2h("BAB"): {2, []int{1, 6}}, // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
				s2h("ABC"): {1, []int{2}},    // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
				s2h("BCA"): {1, []int{3}},    // "424341":pattern.Pat{Weight:1, Pos:[]int{3}},
				s2h("CAB"): {1, []int{4}},    // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
			},
		},
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("ABABCAB"), 3, false, true,
			map[string]Pat{
				s2h("ABA"): {1, []int{0}}, // "414241":pattern.Pat{Weight:2, Pos:[]int{0, 5}},
				s2h("BAB"): {1, []int{1}}, // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
				s2h("ABC"): {1, []int{2}}, // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
				s2h("BCA"): {1, []int{3}}, // "424341":pattern.Pat{Weight:1, Pos:[]int{3}},
				s2h("CAB"): {1, []int{4}}, // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
			},
		},
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("ABCABC"), 3, false, true,
			map[string]Pat{
				s2h("ABC"): {2, []int{0, 3}}, // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
				s2h("BCA"): {1, []int{1}},    // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
				s2h("CAB"): {1, []int{2}},    // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
			},
		},
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("ABCABCABC"), 3, false, true,
			map[string]Pat{
				s2h("ABC"): {3, []int{0, 3, 6}}, // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
				s2h("BCA"): {2, []int{1, 4}},    // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
				s2h("CAB"): {2, []int{2, 5}},    // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.ScanData(tt.data, tt.maxLen, tt.ring)
			if tt.reduce {
				p.Reduce()
			}
			p.DeleteEmptyKeys()
			p.SortPositions()
			assert.Equal(t, tt.exp, tt.fields.Hist)
		})
	}
}

func TestHistogram_Reduce_1_devel(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		data   []byte
		maxLen int
		ring   bool
		reduce bool
		exp    map[string]Pat
	}{ // test cases:
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("ABCABCABCABC"), 4, false, false, // max, ring reduce
			//map[string]Pat{ //////////////////////////// 0123456789ab
			map[string]Pat{
				s2h("AB"):   {4, []int{0, 3, 6, 9}},  //"4142":     {Weight: 4, Pos: []int{0, 3, 6, 9}},
				s2h("BC"):   {4, []int{1, 4, 7, 10}}, //"4243":     {Weight: 4, Pos: []int{1, 4, 7, 10}},
				s2h("CA"):   {3, []int{2, 5, 8}},     //"4341":     {Weight: 3, Pos: []int{2, 5, 8}},
				s2h("ABC"):  {4, []int{0, 3, 6, 9}},  //"414243":   {Weight: 4, Pos: []int{0, 3, 6, 9}},
				s2h("BCA"):  {3, []int{1, 4, 7}},     //"424341":   {Weight: 3, Pos: []int{1, 4, 7}},
				s2h("CAB"):  {3, []int{2, 5, 8}},     //"434142":   {Weight: 3, Pos: []int{2, 5, 8}},
				s2h("ABCA"): {3, []int{0, 3, 6}},     //"41424341": {Weight: 3, Pos: []int{0, 3, 6}},
				s2h("BCAB"): {3, []int{1, 4, 7}},     //"42434142": {Weight: 3, Pos: []int{1, 4, 7}},
				s2h("CABC"): {3, []int{2, 5, 8}},     //"43414243": {Weight: 3, Pos: []int{2, 5, 8}},
			},
		},
		{
			"", fields{map[string]Pat{}, &m, nil}, []byte("ABCABCABCABC"), 4, false, true,
			map[string]Pat{ ////////////////////////////// 0123456789ab
				////////////////////////////////////////// ABCA  ABCA   <- 2 pattern
				////////////////////////////////////////// ABCABCABCABC <- 3 pattern
				s2h("ABCA"): {3, []int{0, 3, 6}}, //"41424341":pattern.Pat{Weight:3, Pos:[]int{0, 3, 6}},
				s2h("BCAB"): {3, []int{1, 4, 7}}, //"42434142":pattern.Pat{Weight:3, Pos:[]int{1, 4, 7}},
				s2h("CABC"): {3, []int{2, 5, 8}}, //"43414243":pattern.Pat{Weight:3, Pos:[]int{2, 5, 8}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.ScanData(tt.data, tt.maxLen, tt.ring)
			if tt.reduce {
				p.Reduce()
			}
			p.DeleteEmptyKeys()
			p.SortPositions()
			assert.Equal(t, tt.exp, tt.fields.Hist)
		})
	}
}

func TestHistogram_AddWeigths(t *testing.T) {
	var mu sync.Mutex
	tests := []struct {
		name string
		p    *Histogram
		exp  map[string]float64
	}{
		// test cases:
		{"", &Histogram{map[string]Pat{"1122": {2, []int{8, 32}}, "112233": {1, []int{8}}}, &mu, nil}, map[string]float64{"1122": 4, "112233": 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.AddWeigths()
			for k, v := range tt.p.Hist {
				a := v.Weight
				e := tt.exp[k]
				fmt.Println(k, a, e)
				assert.True(t, withinTolerance(e, a, 0.001))
			}
		})
	}
}

func TestHistogram_ReduceSubKey(t *testing.T) {
	var mu sync.Mutex
	type args struct {
		bkey   string
		subKey string
	}
	tests := []struct {
		name string
		p    *Histogram
		args args
		exp  *Histogram
	}{
		// test cases:
		{
			"",
			&Histogram{map[string]Pat{s2h("ab"): {2, []int{8, 32}}, s2h("112233"): {1, []int{8}}}, &mu, nil},
			args{s2h("abc"), s2h("ab")},
			&Histogram{map[string]Pat{s2h("ab"): {2, []int{8, 32}}, s2h("112233"): {1, []int{8}}}, &mu, nil},
		},
		{
			"",
			&Histogram{map[string]Pat{s2h("ab"): {2, []int{8, 32}}, s2h("abc"): {1, []int{8}}}, &mu, nil},
			args{s2h("abc"), s2h("ab")},
			&Histogram{map[string]Pat{s2h("ab"): {1, []int{32}}, s2h("abc"): {1, []int{8}}}, &mu, nil},
		},
		{
			"",
			&Histogram{map[string]Pat{s2h("ab"): {2, []int{8, 32}}, s2h("bc"): {2, []int{9, 44}}, s2h("abc"): {1, []int{8}}}, &mu, nil},
			args{s2h("abc"), s2h("bc")},
			&Histogram{map[string]Pat{s2h("ab"): {2, []int{8, 32}}, s2h("bc"): {1, []int{44}}, s2h("abc"): {1, []int{8}}}, &mu, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ReduceSubKey(tt.args.bkey, tt.args.subKey)
			assert.Equal(t, tt.exp, tt.p)
		})
	}
}

// generated: ////////////////////////////////

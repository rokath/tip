package pattern

import (
	"sync"
	"testing"

	"github.com/tj/assert"
)

//  func TestHistogram_BalanceByteUsage(t *testing.T) {
//  	var m sync.Mutex
//  	type fields struct {
//  		Hist map[string]Pattern
//  		mu   *sync.Mutex
//  		Key  []string
//  	}
//  	type args struct {
//  		maxPatternSize int
//  	}
//  	tests := []struct {
//  		name   string
//  		fields fields
//  		args   args
//  		exp    map[string]Pattern
//  	}{ // test cases:
//  		{"", fields{map[string]Pattern{s2h("a"): {[]byte{'a'}, []int{0, 1, 2, 3}}}, &m, nil}, args{4},
//  			/*   */ map[string]Pattern{s2h("a"): {[]byte{'a'}, []int{0, 1, 2, 3}}}},
//
//  		{"", fields{map[string]Pattern{s2h("aa"): {[]byte{'a', 'a'}, []int{0, 1, 2}}}, &m, nil}, args{4},
//  			/*   */ map[string]Pattern{s2h("aa"): {[]byte{'a', 'a'}, []int{0, 1, 2}}}},
//
//  		{"", fields{map[string]Pattern{s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1}}}, &m, nil}, args{4},
//  			/*   */ map[string]Pattern{s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1}}}},
//
//  		{"", fields{map[string]Pattern{s2h("aaaa"): {[]byte{'a', 'a', 'a', 'a'}, []int{0}}}, &m, nil}, args{4},
//  			/*   */ map[string]Pattern{s2h("aaaa"): {[]byte{'a', 'a', 'a', 'a'}, []int{0}}}},
//
//  		{"", fields{map[string]Pattern{s2h("ab"): {[]byte{'a', 'b'}, []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}}}, &m, nil}, args{8},
//  			/*   */ map[string]Pattern{s2h("ab"): {[]byte{'a', 'b'}, []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}}}},
//
//  		{"", fields{map[string]Pattern{s2h("ab"): {[]byte{'a', 'b'}, []int{0, 3, 6}}}, &m, nil}, args{8},
//  			/*   */ map[string]Pattern{s2h("ab"): {[]byte{'a', 'b'}, []int{0, 3, 6}}}},
//
//  		{"", fields{map[string]Pattern{
//  			s2h("ab"):  {[]byte{'a', 'b'}, []int{0, 3}},
//  			s2h("bc"):  {[]byte{'b', 'c'}, []int{1}},
//  			s2h("ca"):  {[]byte{'c', 'a'}, []int{2}},
//  			s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{0}},
//  			s2h("bca"): {[]byte{'b', 'c', 'a'}, []int{1}},
//  			s2h("cab"): {[]byte{'c', 'a', 'b'}, []int{2}},
//  			s2h("aba"): {[]byte{'a', 'b', 'a'}, []int{3}},
//  			s2h("bab"): {[]byte{'b', 'a', 'b'}, []int{4}},
//  		}, &m, nil}, args{3},
//  			map[string]Pattern{
//  				s2h("ab"):  {[]byte{'a', 'b'}, []int{0, 3}},
//  				s2h("bc"):  {[]byte{'b', 'c'}, []int{1}},
//  				s2h("ca"):  {[]byte{'c', 'a'}, []int{2}},
//  				s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{0}},
//  				s2h("bca"): {[]byte{'b', 'c', 'a'}, []int{1}},
//  				s2h("cab"): {[]byte{'c', 'a', 'b'}, []int{2}},
//  				s2h("aba"): {[]byte{'a', 'b', 'a'}, []int{3}},
//  				s2h("bab"): {[]byte{'b', 'a', 'b'}, []int{4}},
//  			},
//  		},
//  	}
//  	for _, tt := range tests {
//  		t.Run(tt.name, func(t *testing.T) {
//  			p := &Histogram{
//  				Hist: tt.fields.Hist,
//  				mu:   tt.fields.mu,
//  				Keys: tt.fields.Key,
//  			}
//  			p.ComputeBalance(tt.args.maxPatternSize)
//  			for k, v := range p.Hist {
//  				a := v.Weight
//  				e := tt.exp[k].Weight
//  				ok := withinTolerance(a, e, 0.001)
//  				if !ok {
//  					fmt.Println(k, e, a)
//  				}
//  				assert.True(t, ok)
//  			}
//  		})
//  	}
//

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
		Hist map[string]Pattern
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   int
		exp    map[string]Pattern
	}{ // test cases:
		{
			"", // abcab: TODO Issue ca occurs 1 times but is inside bca and cab!
			//     01234
			fields{
				map[string]Pattern{
					s2h("ab"):  {[]byte{'a', 'b'}, []int{0, 3}, []int{}}, // in abc and cab and in (aba) and (bab) formally
					s2h("bc"):  {[]byte{'b', 'c'}, []int{1}, []int{}},    // in abc and bca
					s2h("ca"):  {[]byte{'c', 'a'}, []int{2}, []int{}},    // in bca and cab
					s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{0}, []int{}},
					s2h("bca"): {[]byte{'b', 'c', 'a'}, []int{1}, []int{}},
					s2h("cab"): {[]byte{'c', 'a', 'b'}, []int{2}, []int{}},
					s2h("aba"): {[]byte{'a', 'b', 'a'}, []int{3}, []int{}},
					s2h("bab"): {[]byte{'b', 'a', 'b'}, []int{4}, []int{}},
				},
				&m, nil,
			},
			3,
			map[string]Pattern{
				s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{0}, []int{}},
				s2h("bca"): {[]byte{'b', 'c', 'a'}, []int{1}, []int{}},
				s2h("cab"): {[]byte{'c', 'a', 'b'}, []int{2}, []int{}},
				s2h("aba"): {[]byte{'a', 'b', 'a'}, []int{3}, []int{}},
				s2h("bab"): {[]byte{'b', 'a', 'b'}, []int{4}, []int{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Keys: tt.fields.Keys,
			}
			//p.ComputeBalance(tt.args)
			p.ReduceFromSmallerSide()
			p.DeleteEmptyKeys()
			assert.Equal(t, tt.exp, tt.fields.Hist)
		})
	}
}

func TestHistogram_Reduce_1_ok(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pattern
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
		exp    map[string]Pattern
	}{ // test cases:
		{ // name fields                               data    01234  maxLen ring reduce
			"", fields{map[string]Pattern{}, &m, nil}, []byte("aaaaa"), 3, true, true,
			map[string]Pattern{ // expected (with nil instead of []int{} )
				s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1, 2, 3, 4}, nil},
			},
		},
		//  { // name fields                               data    01234  maxLen ring reduce
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("aaaaa"), 3, true, true,
		//  	map[string]Pattern{
		//  		s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1, 2, 3, 4}, []int{}},
		//  	},
		//  },
		//  { // name fields                               data    01234567  maxLen ring reduce
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("aaaabaac"), 4, false, false,
		//  	map[string]Pattern{
		//  		"6161":     {Bytes: []uint8{0x61, 0x61},    Pos: []int{0, 1, 2, 5}, DeletedPos:[]int{}},
		//  		"616161":   {Bytes: []uint8{0x61, 0x61, 0x61},    Pos: []int{0, 1}, DeletedPos:[]int{}},
		//  		"61616161": {Bytes: []uint8{0x61, 0x61, 0x61, 0x61}, Pos: []int{0}, DeletedPos:[]int{}},
		//  		"61616162": {Bytes: []uint8{0x61, 0x61, 0x61, 0x62}, Pos: []int{1}, DeletedPos:[]int{}},
		//  		"616162":   {Bytes: []uint8{0x61, 0x61, 0x62},       Pos: []int{2}, DeletedPos:[]int{}},
		//  		"61616261": {Bytes: []uint8{0x61, 0x61, 0x62, 0x61}, Pos: []int{2}, DeletedPos:[]int{}},
		//  		"616163":   {Bytes: []uint8{0x61, 0x61, 0x63},       Pos: []int{5}, DeletedPos:[]int{}},
		//  		"6162":     {Bytes: []uint8{0x61, 0x62},             Pos: []int{3}, DeletedPos:[]int{}},
		//  		"616261":   {Bytes: []uint8{0x61, 0x62, 0x61},       Pos: []int{3}, DeletedPos:[]int{}},
		//  		"61626161": {Bytes: []uint8{0x61, 0x62, 0x61, 0x61}, Pos: []int{3}, DeletedPos:[]int{}},
		//  		"6163":     {Bytes: []uint8{0x61, 0x63},             Pos: []int{6}, DeletedPos:[]int{}},
		//  		"6261":     {Bytes: []uint8{0x62, 0x61},             Pos: []int{4}, DeletedPos:[]int{}},
		//  		"626161":   {Bytes: []uint8{0x62, 0x61, 0x61},       Pos: []int{4}, DeletedPos:[]int{}},
		//  		"62616163": {Bytes: []uint8{0x62, 0x61, 0x61, 0x63}, Pos: []int{4}, DeletedPos:[]int{}},
		//  	},
		//  },
		//  { // name fields                               data    01234567  maxLen ring reduce
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("aaaabaac"), 4, false, true,
		//  	map[string]Pattern{
		//  		s2h("aaaa"): {[]byte{'a', 'a', 'a', 'a'}, []int{0}, []int{}},
		//  		s2h("aaab"): {[]byte{'a', 'a', 'a', 'b'}, []int{1}, []int{}},
		//  		s2h("aaba"): {[]byte{'a', 'a', 'b', 'a'}, []int{2}, []int{}},
		//  		s2h("abaa"): {[]byte{'a', 'b', 'a', 'a'}, []int{3}, []int{}},
		//  		s2h("baac"): {[]byte{'b', 'a', 'a', 'c'}, []int{4}, []int{}},
		//  	},
		//  },
		//  { // name fields                               data    01234  maxLen ring reduce
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("aaaaa"), 3, true, true,
		//  	map[string]Pattern{
		//  		s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1, 2, 3, 4}, []int{}},
		//  	},
		//  },
		//  { // name fields                               data    01234  maxLen ring reduce
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("aaaaa"), 3, true, false,
		//  	map[string]Pattern{
		//  		s2h("aa"):  {[]byte{'a', 'a'}, []int{0, 1, 2, 3, 4}, []int{}},
		//  		s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1, 2, 3, 4}, []int{}},
		//  	},
		//  },
		//  {
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("ABABCAB"), 3, true, true,
		//  	map[string]Pattern{
		//  		s2h("ABA"): {[]byte{'A', 'B', 'A'}, []int{0, 5}, []int{}}, // "414241":pattern.Pat{Weight:2, Pos:[]int{0, 5}},
		//  		s2h("BAB"): {[]byte{'B', 'A', 'B'}, []int{1, 6}, []int{}}, // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
		//  		s2h("ABC"): {[]byte{'A', 'B', 'C'}, []int{2}, []int{}},    // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
		//  		s2h("BCA"): {[]byte{'B', 'C', 'A'}, []int{3}, []int{}},    // "424341":pattern.Pat{Weight:1, Pos:[]int{3}},
		//  		s2h("CAB"): {[]byte{'C', 'A', 'B'}, []int{4}, []int{}},    // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
		//  	},
		//  },
		//  {
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("ABABCAB"), 3, false, true,
		//  	map[string]Pattern{
		//  		s2h("ABA"): {[]byte{'A', 'B', 'A'}, []int{0}, []int{}}, // "414241":pattern.Pat{Weight:2, Pos:[]int{0, 5}},
		//  		s2h("BAB"): {[]byte{'B', 'A', 'B'}, []int{1}, []int{}}, // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
		//  		s2h("ABC"): {[]byte{'A', 'B', 'C'}, []int{2}, []int{}}, // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
		//  		s2h("BCA"): {[]byte{'B', 'C', 'A'}, []int{3}, []int{}}, // "424341":pattern.Pat{Weight:1, Pos:[]int{3}},
		//  		s2h("CAB"): {[]byte{'C', 'A', 'B'}, []int{4}, []int{}}, // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
		//  	},
		//  },
		//  {
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("ABCABC"), 3, false, true,
		//  	map[string]Pattern{
		//  		s2h("ABC"): {[]byte{'A', 'B', 'C'}, []int{0, 3}, []int{}}, // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
		//  		s2h("BCA"): {[]byte{'B', 'C', 'A'}, []int{1}, []int{}},    // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
		//  		s2h("CAB"): {[]byte{'C', 'A', 'B'}, []int{2}, []int{}},    // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
		//  	},
		//  },
		//  {
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("ABCABCABC"), 3, false, true,
		//  	map[string]Pattern{
		//  		s2h("ABC"): {[]byte{'A', 'B', 'C'}, []int{0, 3, 6}, []int{}}, // "414243":pattern.Pat{Weight:1, Pos:[]int{2}},
		//  		s2h("BCA"): {[]byte{'B', 'C', 'A'}, []int{1, 4}, []int{}},    // "424142":pattern.Pat{Weight:2, Pos:[]int{1, 6}},
		//  		s2h("CAB"): {[]byte{'C', 'A', 'B'}, []int{2, 5}, []int{}},    // "434142":pattern.Pat{Weight:1, Pos:[]int{4}}}
		//  	},
		//  },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Keys: tt.fields.Keys,
			}
			p.ScanData(tt.data, tt.maxLen, tt.ring)
			if tt.reduce {
				p.ReduceFromSmallerSide()
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
		Hist map[string]Pattern
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
		exp    map[string]Pattern
	}{ // test cases:
		{ // name fields                               data                  max, ring reduce
			"", fields{map[string]Pattern{}, &m, nil}, []byte("ABCABCABCABC"), 4, false, false,
			//map[string]Pat{ //////////////////////////////// 0123456789ab
			map[string]Pattern{
				s2h("AB"):   {Bytes: []byte{'A', 'B'}, Pos: []int{0, 3, 6, 9}},        //"4142":     {Weight: 4, Pos: []int{0, 3, 6, 9}},
				s2h("BC"):   {Bytes: []byte{'B', 'C'}, Pos: []int{1, 4, 7, 10}},       //"4243":     {Weight: 4, Pos: []int{1, 4, 7, 10}},
				s2h("CA"):   {Bytes: []byte{'C', 'A'}, Pos: []int{2, 5, 8}},           //"4341":     {Weight: 3, Pos: []int{2, 5, 8}},
				s2h("ABC"):  {Bytes: []byte{'A', 'B', 'C'}, Pos: []int{0, 3, 6, 9}},   //"414243":   {Weight: 4, Pos: []int{0, 3, 6, 9}},
				s2h("BCA"):  {Bytes: []byte{'B', 'C', 'A'}, Pos: []int{1, 4, 7}},      //"424341":   {Weight: 3, Pos: []int{1, 4, 7}},
				s2h("CAB"):  {Bytes: []byte{'C', 'A', 'B'}, Pos: []int{2, 5, 8}},      //"434142":   {Weight: 3, Pos: []int{2, 5, 8}},
				s2h("ABCA"): {Bytes: []byte{'A', 'B', 'C', 'A'}, Pos: []int{0, 3, 6}}, //"41424341": {Weight: 3, Pos: []int{0, 3, 6}},
				s2h("BCAB"): {Bytes: []byte{'B', 'C', 'A', 'B'}, Pos: []int{1, 4, 7}}, //"42434142": {Weight: 3, Pos: []int{1, 4, 7}},
				s2h("CABC"): {Bytes: []byte{'C', 'A', 'B', 'C'}, Pos: []int{2, 5, 8}}, //"43414243": {Weight: 3, Pos: []int{2, 5, 8}},
			},
		},
		{ // name fields                               data                  max, ring reduce
			"", fields{map[string]Pattern{}, &m, nil}, []byte("ABCABCABCABC"), 4, false, true,
			map[string]Pattern{ ////////////////////////////// 0123456789ab
				////////////////////////////////////////////// ABCA  ABCA   <- 2 pattern
				////////////////////////////////////////////// ABCABCABCABC <- 3 pattern
				s2h("ABCA"): {Bytes: []byte{'A', 'B', 'C', 'A'}, Pos: []int{0, 3, 6}}, //"41424341":pattern.Pat{Weight:3, Pos:[]int{0, 3, 6}},
				s2h("BCAB"): {Bytes: []byte{'B', 'C', 'A', 'B'}, Pos: []int{1, 4, 7}}, //"42434142":pattern.Pat{Weight:3, Pos:[]int{1, 4, 7}},
				s2h("CABC"): {Bytes: []byte{'C', 'A', 'B', 'C'}, Pos: []int{2, 5, 8}}, //"43414243":pattern.Pat{Weight:3, Pos:[]int{2, 5, 8}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Keys: tt.fields.Keys,
			}
			p.ScanData(tt.data, tt.maxLen, tt.ring)
			if tt.reduce {
				p.ReduceFromSmallerSide()
			}
			p.DeleteEmptyKeys()
			p.SortPositions()
			assert.Equal(t, tt.exp, tt.fields.Hist)
		})
	}
}

/*
func TestHistogram_ComputeValues(t *testing.T) {
	var mu sync.Mutex
	tests := []struct {
		name string
		p    *Histogram
		exp  map[string]float64
	}{ // test cases:
		{"", &Histogram{
			map[string]Pattern{ // data
				"1122":   {[]byte{0x11, 0x22}, []int{8, 32}}, // count of positions is 2
				"112233": {[]byte{0x11, 0x22, 0x33}, []int{8}}, // count of positions is 1
			},
			&mu, nil,
		},
			map[string]float64{ // exp
				"1122":   4,
				"112233": 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ComputeValues(8)
			for k, v := range tt.p.Hist {
				a := v.Weight
				e := tt.exp[k]
				fmt.Println(k, a, e)
				assert.True(t, withinTolerance(e, a, 0.001))
			}
		})
	}
}
*/

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
	}{ // test cases:
		{"",
			&Histogram{ // data
				map[string]Pattern{
					s2h("ab"):     {Bytes: []byte{'a', 'b'}, Pos: []int{8, 32}},
					s2h("112233"): {Bytes: []byte{0x11, 0x22, 0x33}, Pos: []int{8}},
				},
				&mu, nil,
			},
			args{ // param
				s2h("abc"), // key
				s2h("ab"),  // subkey
			},
			&Histogram{map[string]Pattern{ // expected data (unchanged)
				s2h("ab"):     {Bytes: []byte{'a', 'b'}, Pos: []int{8, 32}},
				s2h("112233"): {Bytes: []byte{0x11, 0x22, 0x33}, Pos: []int{8}},
			}, &mu, nil},
		},
		{"",
			&Histogram{map[string]Pattern{ // data
				s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{8, 32}},
				s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},
			}, &mu, nil},
			args{s2h("abc"), s2h("ab")}, // param key & subkey
			&Histogram{map[string]Pattern{ // expected data (removed position 8 in ab)
				s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{32}},
				s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},
			}, &mu, nil}, // exp
		},
		{"",
			&Histogram{map[string]Pattern{ // data
				s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{8, 32}},
				s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{9, 44}},
				s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},
			}, &mu, nil},
			args{s2h("abc"), s2h("bc")}, // param key & subkey
			&Histogram{map[string]Pattern{ // expected data (removed position 9 in bc)
				s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{8, 32}},
				s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{44}},
				s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},
			}, &mu, nil},
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

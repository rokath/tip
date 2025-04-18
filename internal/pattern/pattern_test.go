package pattern

import (
	"fmt"
	"reflect"
	"slices"
	"sync"
	"testing"

	"github.com/tj/assert"
)

func TestHistogramGeneration(t *testing.T) {
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
		exp    map[string]Pattern
	}{ // test cases:
		{ // name fields                               data    01234567  max, ring,
			"", fields{map[string]Pattern{}, &m, nil}, []byte("xxAAAAyy"), 4, false,
			map[string]Pattern{ // expected
				s2h("xxAA"): {Pos: []int{0}},    //// xxAA     @ 0
				s2h("xAAA"): {Pos: []int{1}},    ////  xAAA    @ 1
				s2h("AAAA"): {Pos: []int{2}},    ////   AAAA   @ 2
				s2h("AAAy"): {Pos: []int{3}},    ////    AAAy  @ 3
				s2h("AAyy"): {Pos: []int{4}},    ////     AAyy @ 4
				s2h("xxA"):  {Pos: []int{0}},    //// xxA      @ 0
				s2h("xAA"):  {Pos: []int{1}},    ////  xAA     @ 1
				s2h("AAA"):  {Pos: []int{2, 3}}, ////   AAA    @ 2 = 2 3
				//                                 //    AAA   @ 3 = 2 3
				s2h("AAy"): {Pos: []int{4}},       //     AAy  @ 4
				s2h("Ayy"): {Pos: []int{5}},       //      Ayy @ 5
				s2h("xx"):  {Pos: []int{0}},       // xx       @ 0
				s2h("xA"):  {Pos: []int{1}},       //  xA      @ 1
				s2h("AA"):  {Pos: []int{2, 3, 4}}, //   AA     @ 2 = 2 3 4
				s2h("Ay"):  {Pos: []int{5}},       //      Ay  @ 5
				s2h("yy"):  {Pos: []int{6}},       //       yy @ 6
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

			for k, v := range tt.exp {
				fmt.Println(k, v.Pos, tt.fields.Hist[k].Pos)
			}
			for k, v := range tt.exp {
				act := tt.fields.Hist[k].Pos
				slices.Sort(act)
				assert.Equal(t, v.Pos, act)
			}
		})
	}
}

// GetKeys extracts all p.Hist keys into p.Keys.
func (p *Histogram) GetKeys() {
	p.mu.Lock()
	for key := range p.Hist {
		p.Keys = append(p.Keys, key)
	}
	p.mu.Unlock()
}

func TestHistogram_scanForRepetitions(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pattern
		mu   *sync.Mutex
	}
	type args struct {
		data  []byte
		ptLen int
		ring  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    map[string]Pattern
	}{ // test cases:
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("ab"), 2, false},
			map[string]Pattern{s2h("ab"): {Bytes: []byte{'a', 'b'}, Pos: []int{0}}},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("ab"), 2, true},
			map[string]Pattern{
				s2h("ab"): {Bytes: []byte{'a', 'b'}, Pos: []int{0}},
				s2h("ba"): {Bytes: []byte{'b', 'a'}, Pos: []int{1}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 2, false},
			map[string]Pattern{
				s2h("ab"): {Bytes: []byte{'a', 'b'}, Pos: []int{0}},
				s2h("bc"): {Bytes: []byte{'b', 'c'}, Pos: []int{1}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 2, true},
			map[string]Pattern{
				s2h("ab"): {Bytes: []byte{'a', 'b'}, Pos: []int{0}},
				s2h("bc"): {Bytes: []byte{'b', 'c'}, Pos: []int{1}},
				s2h("ca"): {Bytes: []byte{'c', 'a'}, Pos: []int{2}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaa"), 2, false},
			map[string]Pattern{
				s2h("aa"): {Bytes: []byte{'a', 'a'}, Pos: []int{0, 1}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaa"), 2, true},
			map[string]Pattern{
				s2h("aa"): {Bytes: []byte{'a', 'a'}, Pos: []int{0, 1, 2}},
			},
		},
	}
	for _, tt := range tests {
		p := &Histogram{
			Hist: tt.fields.Hist,
			mu:   tt.fields.mu,
		}
		p.scanForRepetitions(tt.args.data, tt.args.ptLen, tt.args.ring)
		p.SortPositions()
		assert.Equal(t, tt.exp, p.Hist)
	}
}

func TestHistogram_ScanData(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pattern
		mu   *sync.Mutex
	}
	type args struct {
		data           []byte
		maxPatternSize int
		ring           bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    map[string]Pattern
	}{ // test cases:
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 3, false},
			map[string]Pattern{
				"22aaaa": {Bytes: []byte{0x22, 0xaa, 0xaa}, Pos: []int{4}},
				"aa22aa": {Bytes: []byte{0xaa, 0x22, 0xaa}, Pos: []int{3}},
				"aaaa22": {Bytes: []byte{0xaa, 0xaa, 0x22}, Pos: []int{2}},
				"aaaaaa": {Bytes: []byte{0xaa, 0xaa, 0xaa}, Pos: []int{0, 1, 5}},
				"22aa":   {Bytes: []byte{0x22, 0xaa}, Pos: []int{4}},
				"aa22":   {Bytes: []byte{0xaa, 0x22}, Pos: []int{3}},
				"aaaa":   {Bytes: []byte{0xaa, 0xaa}, Pos: []int{0, 1, 2, 5, 6}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}, 4, false},
			map[string]Pattern{
				"0000":     {Bytes: []byte{0x00, 0x00}, Pos: []int{4, 10}},
				"0000ff":   {Bytes: []byte{0x00, 0x00, 0xff}, Pos: []int{4, 10}},
				"0000ffff": {Bytes: []byte{0x00, 0x00, 0xff, 0xff}, Pos: []int{4, 10}},
				"00ff":     {Bytes: []byte{0x00, 0xff}, Pos: []int{5, 11}},
				"00ffff":   {Bytes: []byte{0x00, 0xff, 0xff}, Pos: []int{5, 11}},
				"00ffffff": {Bytes: []byte{0x00, 0xff, 0xff, 0xff}, Pos: []int{5, 11}},
				"ff00":     {Bytes: []byte{0xff, 0x00}, Pos: []int{3, 9}},
				"ff0000":   {Bytes: []byte{0xff, 0x00, 0x00}, Pos: []int{3, 9}},
				"ff0000ff": {Bytes: []byte{0xff, 0x00, 0x00, 0xff}, Pos: []int{3, 9}},
				"ffff":     {Bytes: []byte{0xff, 0xff}, Pos: []int{0, 1, 2, 6, 7, 8, 12, 13, 14}},
				"ffff00":   {Bytes: []byte{0xff, 0xff, 0x00}, Pos: []int{2, 8}},
				"ffff0000": {Bytes: []byte{0xff, 0xff, 0x00, 0x00}, Pos: []int{2, 8}},
				"ffffff":   {Bytes: []byte{0xff, 0xff, 0xff}, Pos: []int{0, 1, 6, 7, 12, 13}},
				"ffffff00": {Bytes: []byte{0xff, 0xff, 0xff, 0x00}, Pos: []int{1, 7}},
				"ffffffff": {Bytes: []byte{0xff, 0xff, 0xff, 0xff}, Pos: []int{0, 6, 12}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff}, 3, false},
			map[string]Pattern{
				"ffff":   {Bytes: []byte{0xff, 0xff}, Pos: []int{0, 1, 2}},
				"ffffff": {Bytes: []byte{0xff, 0xff, 0xff}, Pos: []int{0, 1}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 3, false},
			map[string]Pattern{
				"1122":   {Bytes: []byte{0x11, 0x22}, Pos: []int{0}},
				"112233": {Bytes: []byte{0x11, 0x22, 0x33}, Pos: []int{0}},
				"2233":   {Bytes: []byte{0x22, 0x33}, Pos: []int{1, 4}},
				"2233aa": {Bytes: []byte{0x22, 0x33, 0xaa}, Pos: []int{1}},
				"2233bb": {Bytes: []byte{0x22, 0x33, 0xbb}, Pos: []int{4}},
				"33aa":   {Bytes: []byte{0x33, 0xaa}, Pos: []int{2}},
				"33aa22": {Bytes: []byte{0x33, 0xaa, 0x22}, Pos: []int{2}},
				"33bb":   {Bytes: []byte{0x33, 0xbb}, Pos: []int{5}},
				"aa22":   {Bytes: []byte{0xaa, 0x22}, Pos: []int{3}},
				"aa2233": {Bytes: []byte{0xaa, 0x22, 0x33}, Pos: []int{3}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0x11, 0x22, 0x33}, 2, false},
			map[string]Pattern{
				"1122": {Bytes: []byte{0x11, 0x22}, Pos: []int{0}},
				"2233": {Bytes: []byte{0x22, 0x33}, Pos: []int{1}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 2, false},
			map[string]Pattern{
				s2h("ab"): {Bytes: []byte{'a', 'b'}, Pos: []int{0}},
				s2h("bc"): {Bytes: []byte{'b', 'c'}, Pos: []int{1}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 3, false},
			map[string]Pattern{
				s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{0}},
				s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{1}},
				s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{0}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 3, true},
			map[string]Pattern{
				s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{0}},
				s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{1}},
				s2h("ca"):  {Bytes: []byte{'c', 'a'}, Pos: []int{2}},
				s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{0}},
				s2h("bca"): {Bytes: []byte{'b', 'c', 'a'}, Pos: []int{1}},
				s2h("cab"): {Bytes: []byte{'c', 'a', 'b'}, Pos: []int{2}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaa"), 3, true},
			map[string]Pattern{
				s2h("aa"):  {Bytes: []byte{'a', 'a'}, Pos: []int{0, 1, 2}},
				s2h("aaa"): {Bytes: []byte{'a', 'a', 'a'}, Pos: []int{0, 1, 2}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaaaa"), 3, false},
			map[string]Pattern{
				s2h("aa"):  {Bytes: []byte{'a', 'a'}, Pos: []int{0, 1, 2, 3}},
				s2h("aaa"): {Bytes: []byte{'a', 'a', 'a'}, Pos: []int{0, 1, 2}},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaaaa"), 3, true},
			map[string]Pattern{
				s2h("aa"):  {Bytes: []byte{'a', 'a'}, Pos: []int{0, 1, 2, 3, 4}},
				s2h("aaa"): {Bytes: []byte{'a', 'a', 'a'}, Pos: []int{0, 1, 2, 3, 4}},
			},
		},
	}
	for _, tt := range tests {
		p := &Histogram{
			Hist: tt.fields.Hist,
			mu:   tt.fields.mu,
		}
		p.ScanData(tt.args.data, tt.args.maxPatternSize, tt.args.ring)
		p.SortPositions()
		assert.Equal(t, tt.exp, p.Hist)
	}
}

func TestHistogram_DiscardSeldomPattern(t *testing.T) {
	var m sync.Mutex
	type args struct {
		discardSize int
	}
	tests := []struct {
		name string
		p    *Histogram
		args args
		exp  *Histogram
	}{ // test cases:
		{"",
			&Histogram{
				map[string]Pattern{ // data
					s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{8, 16, 24, 32}}, // 4 positions
					s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{44, 66, 88}},    // 3 positions
					s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},        // 1 position
					s2h("xyz"): {Bytes: []byte{'x', 'y', 'z'}, Pos: []int{}},         // 0 position
				},
				&m, nil,
			},
			args{3}, // discard 0, 1, 2, 3
			&Histogram{
				map[string]Pattern{ // expected
					s2h("ab"): {Bytes: []byte{'a', 'b'}, Pos: []int{8, 16, 24, 32}}, // 4 positions
				},
				&m, nil,
			},
		},
		{"",
			&Histogram{
				map[string]Pattern{ // data
					s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{8, 16, 24, 32}}, // 4 positions
					s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{44, 66, 88}},    // 3 positions
					s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},        // 1 position
					s2h("xyz"): {Bytes: []byte{'x', 'y', 'z'}, Pos: []int{}},         // 0 position
				},
				&m, nil,
			}, args{1}, // discard 0, 1
			&Histogram{
				map[string]Pattern{ // expected
					s2h("ab"): {Bytes: []byte{'a', 'b'}, Pos: []int{8, 16, 24, 32}}, // 4 positions
					s2h("bc"): {Bytes: []byte{'b', 'c'}, Pos: []int{44, 66, 88}},    // 3 positions
				},
				&m, nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.DiscardSeldomPattern(tt.args.discardSize)
			assert.Equal(t, tt.exp, tt.p)
		})
	}
}

func TestHistogram_GetKeys(t *testing.T) {
	var m sync.Mutex
	tests := []struct {
		name string
		p    *Histogram
		exp  []string
	}{ // test cases:
		{"",
			&Histogram{
				map[string]Pattern{ // data
					s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{8, 16, 24, 32}},
					s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{44}},
					s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},
				},
				&m, nil,
			},
			[]string{"ab", "bc", "abc"}, // expected
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.GetKeys()
			assert.Equal(t, len(tt.p.Hist), len(tt.exp))
			for _, k := range tt.exp {
				_, ok := tt.p.Hist[s2h(k)]
				assert.True(t, ok)
			}
		})
	}
}

func TestHistogram_ExportAsList(t *testing.T) {
	var m sync.Mutex
	tests := []struct {
		name     string
		p        *Histogram
		wantList []Pattern
	}{
		// test cases:
		{"",
			&Histogram{
				map[string]Pattern{ // data
					s2h("ab"):  {Bytes: []byte{'a', 'b'}, Pos: []int{8, 16, 24, 32}},
					s2h("bc"):  {Bytes: []byte{'b', 'c'}, Pos: []int{44}},
					s2h("abc"): {Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},
				},
				&m, nil,
			},
			[]Pattern{ // expected
				{Bytes: []byte{'a', 'b'}, Pos: []int{8, 16, 24, 32}},
				{Bytes: []byte{'b', 'c'}, Pos: []int{44}},
				{Bytes: []byte{'a', 'b', 'c'}, Pos: []int{8}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList := tt.p.ExportAsList()
			actList := SortByDescCount(gotList)
			expList := SortByDescCount(tt.wantList)
			if !reflect.DeepEqual(actList, expList) {
				t.Errorf("Histogram.ExportAsList() = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

package pattern

import (
	"reflect"
	"sync"
	"testing"

	"github.com/tj/assert"
)

func TestHistogram_scanForRepetitions(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pat
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
		exp    map[string]Pat
	}{ // test cases:
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("ab"), 2, false},
			map[string]Pat{s2h("ab"): {1, []int{0}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("ab"), 2, true},
			map[string]Pat{s2h("ab"): {1, []int{0}},s2h("ba"): {1, []int{1}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("abc"), 2, false},
			map[string]Pat{s2h("ab"): {1, []int{0}},s2h("bc"): {1, []int{1}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("abc"), 2, true},
			map[string]Pat{s2h("ab"): {1, []int{0}},s2h("bc"): {1, []int{1}},s2h("ca"): {1, []int{2}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("aaa"), 2, false},
			map[string]Pat{s2h("aa"): {2, []int{0,1}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("aaa"), 2, true},
			map[string]Pat{s2h("aa"): {3, []int{0,1,2}}},
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
		Hist map[string]Pat
		mu   *sync.Mutex
	}
	type args struct {
		data           []byte
		maxPatternSize int
		ring bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    map[string]Pat
	}{ // test cases:
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 3, false},
			map[string]Pat{"22aaaa": {1, []int{4}}, "aa22aa": {1, []int{3}}, "aaaa22": {1, []int{2}}, "aaaaaa": {3, []int{0, 1, 5}}, "22aa": {1, []int{4}}, "aa22": {1, []int{3}}, "aaaa": {5, []int{0, 1, 2, 5, 6}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}, 4, false},
			map[string]Pat{
				"0000":     {2, []int{4, 10}},
				"0000ff":   {2, []int{4, 10}},
				"0000ffff": {2, []int{4, 10}},
				"00ff":     {2, []int{5, 11}},
				"00ffff":   {2, []int{5, 11}},
				"00ffffff": {2, []int{5, 11}},
				"ff00":     {2, []int{3, 9}},
				"ff0000":   {2, []int{3, 9}},
				"ff0000ff": {2, []int{3, 9}},
				"ffff":     {9, []int{0, 1, 2, 6, 7, 8, 12, 13, 14}},
				"ffff00":   {2, []int{2, 8}},
				"ffff0000": {2, []int{2, 8}},
				"ffffff":   {6, []int{0, 1, 6, 7, 12, 13}},
				"ffffff00": {2, []int{1, 7}},
				"ffffffff": {3, []int{0, 6, 12}},
			},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff}, 3, false},
			map[string]Pat{"ffff": {3, []int{0, 1, 2}}, "ffffff": {2, []int{0, 1}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 3, false},
			map[string]Pat{"1122": {1, []int{0}}, "112233": {1, []int{0}}, "2233": {2, []int{1, 4}}, "2233aa": {1, []int{1}}, "2233bb": {1, []int{4}}, "33aa": {1, []int{2}}, "33aa22": {1, []int{2}}, "33bb": {1, []int{5}}, "aa22": {1, []int{3}}, "aa2233": {1, []int{3}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33}, 2, false},
			map[string]Pat{"1122": {1, []int{0}}, "2233": {1, []int{1}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("abc"), 2, false},
			map[string]Pat{s2h("ab"): {1, []int{0}}, s2h("bc"): {1, []int{1}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("abc"), 3, false},
			map[string]Pat{s2h("ab"): {1, []int{0}}, s2h("bc"): {1, []int{1}}, s2h("abc"): {1, []int{0}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("abc"), 3, true},
			map[string]Pat{s2h("ab"): {1, []int{0}}, s2h("bc"): {1, []int{1}}, s2h("ca"): {1, []int{2}}, s2h("abc"): {1, []int{0}}, s2h("bca"): {1, []int{1}}, s2h("cab"): {1, []int{2}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("aaa"), 3, true},
			map[string]Pat{s2h("aa"): {3, []int{0,1,2}}, s2h("aaa"): {3, []int{0,1,2}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("aaaaa"), 3, false},
			map[string]Pat{s2h("aa"): {4, []int{0,1,2,3}}, s2h("aaa"): {3, []int{0,1,2}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte("aaaaa"), 3, true},
			map[string]Pat{s2h("aa"): {5, []int{0,1,2,3,4}}, s2h("aaa"): {5, []int{0,1,2,3,4}}},
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
		discardSize float64
	}
	tests := []struct {
		name string
		p    *Histogram
		args args
		exp  *Histogram
	}{
		// test cases:
		{"",
			&Histogram{map[string]Pat{s2h("ab"): {4, []int{8, 16, 24, 32}}, s2h("bc"): {1, []int{44}}, s2h("abc"): {1, []int{8}}}, &m, nil}, args{3},
			&Histogram{map[string]Pat{s2h("ab"): {4, []int{8, 16, 24, 32}}}, &m, nil},
		},
		{"",
			&Histogram{map[string]Pat{s2h("ab"): {4, []int{8, 16, 24, 32}}, s2h("bc"): {2, []int{44}}, s2h("abc"): {1, []int{8}}}, &m, nil}, args{1},
			&Histogram{map[string]Pat{s2h("ab"): {4, []int{8, 16, 24, 32}}, s2h("bc"): {2, []int{44}}}, &m, nil},
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
	}{
		// test cases:
		{"",
			&Histogram{map[string]Pat{s2h("ab"): {4, []int{8, 16, 24, 32}}, s2h("bc"): {1, []int{44}}, s2h("abc"): {1, []int{8}}}, &m, nil},
			[]string{"ab", "bc", "abc"},
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
		wantList []Patt
	}{
		// test cases:
		{"",
			&Histogram{map[string]Pat{s2h("ab"): {4, []int{8, 16, 24, 32}}, s2h("bc"): {1, []int{44}}, s2h("abc"): {1, []int{8}}}, &m, nil},
			[]Patt{
				{4, []byte{0x61, 0x62}, s2h("ab")},
				{1, []byte{0x62, 0x63}, s2h("bc")},
				{1, []byte{0x61, 0x62, 0x63}, s2h("abc")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList := tt.p.ExportAsList()
			actList := SortByDescCountDescLength(gotList)
			expList := SortByDescCountDescLength(tt.wantList)
			if !reflect.DeepEqual(actList, expList) {
				t.Errorf("Histogram.ExportAsList() = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

// generated: ////////////////////////////////

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
			map[string]Pattern{s2h("ab"): {[]byte{'a', 'b'}, []int{0}, 0, 0, 0, 0}},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("ab"), 2, true},
			map[string]Pattern{
				s2h("ab"): {[]byte{'a', 'b'}, []int{0}, 0, 0, 0, 0},
				s2h("ba"): {[]byte{'b', 'a'}, []int{1}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 2, false},
			map[string]Pattern{
				s2h("ab"): {[]byte{'a', 'b'}, []int{0}, 0, 0, 0, 0},
				s2h("bc"): {[]byte{'b', 'c'}, []int{1}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 2, true},
			map[string]Pattern{
				s2h("ab"): {[]byte{'a', 'b'}, []int{0}, 0, 0, 0, 0},
				s2h("bc"): {[]byte{'b', 'c'}, []int{1}, 0, 0, 0, 0},
				s2h("ca"): {[]byte{'c', 'a'}, []int{2}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaa"), 2, false},
			map[string]Pattern{
				s2h("aa"): {[]byte{'a', 'a'}, []int{0, 1}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaa"), 2, true},
			map[string]Pattern{
				s2h("aa"): {[]byte{'a', 'a'}, []int{0, 1, 2}, 0, 0, 0, 0},
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
				"22aaaa": {[]byte{0x22, 0xaa, 0xaa}, []int{4}, 0, 0, 0, 0},
				"aa22aa": {[]byte{0xaa, 0x22, 0xaa}, []int{3}, 0, 0, 0, 0},
				"aaaa22": {[]byte{0xaa, 0xaa, 0x22}, []int{2}, 0, 0, 0, 0},
				"aaaaaa": {[]byte{0xaa, 0xaa, 0xaa}, []int{0, 1, 5}, 0, 0, 0, 0},
				"22aa":   {[]byte{0x22, 0xaa}, []int{4}, 0, 0, 0, 0},
				"aa22":   {[]byte{0xaa, 0x22}, []int{3}, 0, 0, 0, 0},
				"aaaa":   {[]byte{0xaa, 0xaa}, []int{0, 1, 2, 5, 6}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}, 4, false},
			map[string]Pattern{
				"0000":     {[]byte{0x00, 0x00}, []int{4, 10}, 0, 0, 0, 0},
				"0000ff":   {[]byte{0x00, 0x00, 0xff}, []int{4, 10}, 0, 0, 0, 0},
				"0000ffff": {[]byte{0x00, 0x00, 0xff, 0xff}, []int{4, 10}, 0, 0, 0, 0},
				"00ff":     {[]byte{0x00, 0xff}, []int{5, 11}, 0, 0, 0, 0},
				"00ffff":   {[]byte{0x00, 0xff, 0xff}, []int{5, 11}, 0, 0, 0, 0},
				"00ffffff": {[]byte{0x00, 0xff, 0xff, 0xff}, []int{5, 11}, 0, 0, 0, 0},
				"ff00":     {[]byte{0xff, 0x00}, []int{3, 9}, 0, 0, 0, 0},
				"ff0000":   {[]byte{0xff, 0x00, 0x00}, []int{3, 9}, 0, 0, 0, 0},
				"ff0000ff": {[]byte{0xff, 0x00, 0x00, 0xff}, []int{3, 9}, 0, 0, 0, 0},
				"ffff":     {[]byte{0xff, 0xff}, []int{0, 1, 2, 6, 7, 8, 12, 13, 14}, 0, 0, 0, 0},
				"ffff00":   {[]byte{0xff, 0xff, 0x00}, []int{2, 8}, 0, 0, 0, 0},
				"ffff0000": {[]byte{0xff, 0xff, 0x00, 0x00}, []int{2, 8}, 0, 0, 0, 0},
				"ffffff":   {[]byte{0xff, 0xff, 0xff}, []int{0, 1, 6, 7, 12, 13}, 0, 0, 0, 0},
				"ffffff00": {[]byte{0xff, 0xff, 0xff, 0x00}, []int{1, 7}, 0, 0, 0, 0},
				"ffffffff": {[]byte{0xff, 0xff, 0xff, 0xff}, []int{0, 6, 12}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff}, 3, false},
			map[string]Pattern{
				"ffff":   {[]byte{0xff, 0xff}, []int{0, 1, 2}, 0, 0, 0, 0},
				"ffffff": {[]byte{0xff, 0xff, 0xff}, []int{0, 1}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 3, false},
			map[string]Pattern{
				"1122":   {[]byte{0x11, 0x22}, []int{0}, 0, 0, 0, 0},
				"112233": {[]byte{0x11, 0x22, 0x33}, []int{0}, 0, 0, 0, 0},
				"2233":   {[]byte{0x22, 0x33}, []int{1, 4}, 0, 0, 0, 0},
				"2233aa": {[]byte{0x22, 0x33, 0xaa}, []int{1}, 0, 0, 0, 0},
				"2233bb": {[]byte{0x22, 0x33, 0xbb}, []int{4}, 0, 0, 0, 0},
				"33aa":   {[]byte{0x33, 0xaa}, []int{2}, 0, 0, 0, 0},
				"33aa22": {[]byte{0x33, 0xaa, 0x22}, []int{2}, 0, 0, 0, 0},
				"33bb":   {[]byte{0x33, 0xbb}, []int{5}, 0, 0, 0, 0},
				"aa22":   {[]byte{0xaa, 0x22}, []int{3}, 0, 0, 0, 0},
				"aa2233": {[]byte{0xaa, 0x22, 0x33}, []int{3}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte{0x11, 0x22, 0x33}, 2, false},
			map[string]Pattern{
				"1122": {[]byte{0x11, 0x22}, []int{0}, 0, 0, 0, 0},
				"2233": {[]byte{0x22, 0x33}, []int{1}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 2, false},
			map[string]Pattern{
				s2h("ab"): {[]byte{'a', 'b'}, []int{0}, 0, 0, 0, 0},
				s2h("bc"): {[]byte{'b', 'c'}, []int{1}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 3, false},
			map[string]Pattern{
				s2h("ab"):  {[]byte{'a', 'b'}, []int{0}, 0, 0, 0, 0},
				s2h("bc"):  {[]byte{'b', 'c'}, []int{1}, 0, 0, 0, 0},
				s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{0}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("abc"), 3, true},
			map[string]Pattern{
				s2h("ab"):  {[]byte{'a', 'b'}, []int{0}, 0, 0, 0, 0},
				s2h("bc"):  {[]byte{'b', 'c'}, []int{1}, 0, 0, 0, 0},
				s2h("ca"):  {[]byte{'c', 'a'}, []int{2}, 0, 0, 0, 0},
				s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{0}, 0, 0, 0, 0},
				s2h("bca"): {[]byte{'b', 'c', 'a'}, []int{1}, 0, 0, 0, 0},
				s2h("cab"): {[]byte{'c', 'a', 'b'}, []int{2}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaa"), 3, true},
			map[string]Pattern{
				s2h("aa"):  {[]byte{'a', 'a'}, []int{0, 1, 2}, 0, 0, 0, 0},
				s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1, 2}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaaaa"), 3, false},
			map[string]Pattern{
				s2h("aa"):  {[]byte{'a', 'a'}, []int{0, 1, 2, 3}, 0, 0, 0, 0},
				s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1, 2}, 0, 0, 0, 0},
			},
		},
		{
			"", // name
			fields{map[string]Pattern{}, &m},
			args{[]byte("aaaaa"), 3, true},
			map[string]Pattern{
				s2h("aa"):  {[]byte{'a', 'a'}, []int{0, 1, 2, 3, 4}, 0, 0, 0, 0},
				s2h("aaa"): {[]byte{'a', 'a', 'a'}, []int{0, 1, 2, 3, 4}, 0, 0, 0, 0},
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
					s2h("ab"):  {[]byte{'a', 'b'}, []int{8, 16, 24, 32}, 0, 0, 0, 0}, // 4 positions
					s2h("bc"):  {[]byte{'b', 'c'}, []int{44, 66, 88}, 0, 0, 0, 0},    // 3 positions
					s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{8}, 0, 0, 0, 0},        // 1 position
					s2h("xyz"): {[]byte{'x', 'y', 'z'}, []int{}, 0, 0, 0, 0},         // 0 position
				},
				&m, nil,
			},
			args{3}, // discard 0, 1, 2, 3
			&Histogram{
				map[string]Pattern{ // expected
					s2h("ab"): {[]byte{'a', 'b'}, []int{8, 16, 24, 32}, 0, 0, 0, 0}, // 4 positions
				},
				&m, nil,
			},
		},
		{"",
			&Histogram{
				map[string]Pattern{ // data
					s2h("ab"):  {[]byte{'a', 'b'}, []int{8, 16, 24, 32}, 0, 0, 0, 0}, // 4 positions
					s2h("bc"):  {[]byte{'b', 'c'}, []int{44, 66, 88}, 0, 0, 0, 0},    // 3 positions
					s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{8}, 0, 0, 0, 0},        // 1 position
					s2h("xyz"): {[]byte{'x', 'y', 'z'}, []int{}, 0, 0, 0, 0},         // 0 position
				},
				&m, nil,
			}, args{1}, // discard 0, 1
			&Histogram{
				map[string]Pattern{ // expected
					s2h("ab"): {[]byte{'a', 'b'}, []int{8, 16, 24, 32}, 0, 0, 0, 0}, // 4 positions
					s2h("bc"): {[]byte{'b', 'c'}, []int{44, 66, 88}, 0, 0, 0, 0},    // 3 positions
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
					s2h("ab"):  {[]byte{'a', 'b'}, []int{8, 16, 24, 32}, 0, 0, 0, 0},
					s2h("bc"):  {[]byte{'b', 'c'}, []int{44}, 0, 0, 0, 0},
					s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{8}, 0, 0, 0, 0},
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
					s2h("ab"):  {[]byte{'a', 'b'}, []int{8, 16, 24, 32}, 0, 0, 0, 0},
					s2h("bc"):  {[]byte{'b', 'c'}, []int{44}, 0, 0, 0, 0},
					s2h("abc"): {[]byte{'a', 'b', 'c'}, []int{8}, 0, 0, 0, 0},
				},
				&m, nil,
			},
			[]Pattern{ // expected
				{[]byte{'a', 'b'}, []int{8, 16, 24, 32}, 0, 0, 0, 0},
				{[]byte{'b', 'c'}, []int{44}, 0, 0, 0, 0},
				{[]byte{'a', 'b', 'c'}, []int{8}, 0, 0, 0, 0},
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

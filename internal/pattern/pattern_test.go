package pattern

import (
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
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 2},
			map[string]Pat{"22aa": {1, []int{4}}, "aa22": {1, []int{3}}, "aaaa": {5, []int{0, 1, 2, 5, 6}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 3},
			map[string]Pat{"22aaaa": {1, []int{4}}, "aa22aa": {1, []int{3}}, "aaaa22": {1, []int{2}}, "aaaaaa": {3, []int{0, 1, 5}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 3},
			map[string]Pat{"112233": {1, []int{0}}, "2233aa": {1, []int{1}}, "2233bb": {1, []int{4}}, "33aa22": {1, []int{2}}, "aa2233": {1, []int{3}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 2},
			map[string]Pat{"1122": {1, []int{0}}, "2233": {2, []int{1, 4}}, "33aa": {1, []int{2}}, "33bb": {1, []int{5}}, "aa22": {1, []int{3}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0x22, 0x33}, 2},
			map[string]Pat{"1122": {1, []int{0}}, "2233": {2, []int{1, 3}}, "3322": {1, []int{2}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33}, 2},
			map[string]Pat{"1122": {1, []int{0}}, "2233": {1, []int{1}}},
		},
	}
	for _, tt := range tests {
		p := &Histogram{
			Hist: tt.fields.Hist,
			mu:   tt.fields.mu,
		}
		p.scanForRepetitions(tt.args.data, tt.args.ptLen)
		p.SortPositions()
		assert.Equal(t, tt.exp, p.Hist)
	}
}

func TestHistogram_Extend(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
	}
	type args struct {
		data           []byte
		maxPatternSize int
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
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 3},
			map[string]Pat{"22aaaa": {1, []int{4}}, "aa22aa": {1, []int{3}}, "aaaa22": {1, []int{2}}, "aaaaaa": {3, []int{0, 1, 5}}, "22aa": {1, []int{4}}, "aa22": {1, []int{3}}, "aaaa": {5, []int{0, 1, 2, 5, 6}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}, 4},
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
			args{[]byte{0xff, 0xff, 0xff, 0xff}, 3},
			map[string]Pat{"ffff": {3, []int{0, 1, 2}}, "ffffff": {2, []int{0, 1}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 3},
			map[string]Pat{"1122": {1, []int{0}}, "112233": {1, []int{0}}, "2233": {2, []int{1, 4}}, "2233aa": {1, []int{1}}, "2233bb": {1, []int{4}}, "33aa": {1, []int{2}}, "33aa22": {1, []int{2}}, "33bb": {1, []int{5}}, "aa22": {1, []int{3}}, "aa2233": {1, []int{3}}},
		},
		{
			"", // name
			fields{map[string]Pat{}, &m},
			args{[]byte{0x11, 0x22, 0x33}, 2},
			map[string]Pat{"1122": {1, []int{0}}, "2233": {1, []int{1}}},
		},
	}
	for _, tt := range tests {
		p := &Histogram{
			Hist: tt.fields.Hist,
			mu:   tt.fields.mu,
		}
		p.Extend(tt.args.data, tt.args.maxPatternSize)
		assert.Equal(t, tt.exp, p.Hist)
	}
}

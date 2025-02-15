package pattern

import (
	"sync"
	"testing"

	"github.com/tj/assert"
)

// generated: ////////////////////////////////

func TestHistogram_scanForRepetitions(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]int
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
		exp    map[string]int
	}{
		// TODO: Add test cases.
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 2},
			map[string]int{"22aa": 1, "aa22": 1, "aaaa": 5},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 3},
			map[string]int{"22aaaa": 1, "aa22aa": 1, "aaaa22": 1, "aaaaaa": 3},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 3},
			map[string]int{"112233": 1, "2233aa": 1, "2233bb": 1, "33aa22": 1, "aa2233": 1},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 2},
			map[string]int{"1122": 1, "2233": 2, "33aa": 1, "33bb": 1, "aa22": 1},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0x22, 0x33}, 2},
			map[string]int{"1122": 1, "2233": 2, "3322": 1},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0x11, 0x22, 0x33}, 2},
			map[string]int{"1122": 1, "2233": 1},
		},
	}
	for _, tt := range tests {
		p := &Histogram{
			Hist: tt.fields.Hist,
			mu:   tt.fields.mu,
		}
		p.scanForRepetitions(tt.args.data, tt.args.ptLen)
		assert.Equal(t, tt.exp, p.Hist)
	}
}

func TestHistogram_Extend(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]int
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
		exp    map[string]int
	}{
		// TODO: Add test cases.
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0x22, 0xaa, 0xaa, 0xaa}, 3},
			map[string]int{"22aaaa": 1, "aa22aa": 1, "aaaa22": 1, "aaaaaa": 3, "22aa": 1, "aa22": 1, "aaaa": 5},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}, 4},
			map[string]int{"0000": 2, "0000ff": 2, "0000ffff": 2, "00ff": 2, "00ffff": 2, "00ffffff": 2, "ff00": 2, "ff0000": 2, "ff0000ff": 2, "ffff": 9, "ffff00": 2, "ffff0000": 2, "ffffff": 6, "ffffff00": 2, "ffffffff": 3},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0xff, 0xff, 0xff, 0xff}, 3},
			map[string]int{"ffff": 3, "ffffff": 2},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0x11, 0x22, 0x33, 0xaa, 0x22, 0x33, 0xbb}, 3},
			map[string]int{"1122": 1, "112233": 1, "2233": 2, "2233aa": 1, "2233bb": 1, "33aa": 1, "33aa22": 1, "33bb": 1, "aa22": 1, "aa2233": 1},
		},
		{
			"", // name
			fields{map[string]int{}, &m},
			args{[]byte{0x11, 0x22, 0x33}, 2},
			map[string]int{"1122": 1, "2233": 1},
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

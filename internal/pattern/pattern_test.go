package pattern

import (
	"testing"


	"github.com/tj/assert"
)

func TestScanForRepetitions(t *testing.T) {
	type e struct {
		d []byte         // data
		l int            // ptLen
		m map[string]int // expected map
	}
	tt := []e{
		{[]byte{1, 2, 3, 1, 2, 3}, 2, map[string]int{"0102": 2, "0203": 2, "0301": 1}}, // intended
		{[]byte{1, 2, 3, 1, 2, 3}, 3, map[string]int{"010203": 2, "020301": 1, "030102": 1}},
		{[]byte{1, 2, 3, 1, 2, 3}, 4, map[string]int{"01020301": 1, "02030102": 1, "03010203": 1}},

		{[]byte{1, 2, 3, 4, 1, 2, 3, 4}, 2, map[string]int{"0102": 2, "0203": 2, "0304": 2, "0401": 1}},
		{[]byte{1, 2, 3, 4, 1, 2, 3, 4}, 3, map[string]int{"010203": 2, "020304": 2, "030401": 1, "040102": 1}},
		{[]byte{1, 2, 3, 4, 1, 2, 3, 4}, 4, map[string]int{"01020304": 2, "02030401": 1, "03040102": 1, "04010203": 1}},

		{[]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}, 2, map[string]int{"0102": 3, "0203": 3, "0304": 3, "0401": 2}},
		{[]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}, 3, map[string]int{"010203": 3, "020304": 3, "030401": 2, "040102": 2}},
		{[]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}, 4, map[string]int{"01020304": 3, "02030401": 2, "03040102": 2, "04010203": 2}},
	}
	for _, x := range tt {
		m := ScanForRepetitions(x.d, x.l)
		assert.Equal(t, x.m, m)
	}
}

/*

func _Test_reduceSubPatternCounts(t *testing.T) {
	ps := []pat_t{
		{900, []byte{1, 2}, ""},
		{300, []byte{1, 2, 3}, ""},          // 300*{1,2}
		{100, []byte{1, 2, 3, 1, 2, 3}, ""}, // 200*{1,2}, 200*{1,2,3}
	}
	exp := []pat_t{
		{400, []byte{1, 2}, ""},    // 900-300-200
		{100, []byte{1, 2, 3}, ""}, // 300-200
		{100, []byte{1, 2, 3, 1, 2, 3}, ""},
	}
	act := reduceSubPatternCounts(ps)
	assert.Equal(t, exp, act)
}

*/

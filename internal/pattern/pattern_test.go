package pattern

import (
	"fmt"
	"maps"
	"testing"

	"github.com/tj/assert"
)

func Test_buildHistogram(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	tt := []struct {
		data []byte         // data
		max  int            // max pattern length
		exp  map[string]int // expected map
	}{
		{[]byte{1, 2, 3, 1, 2, 3}, 2, map[string]int{"0102": 2, "0203": 2, "0301": 1}},
		{[]byte{1, 2, 3, 1, 2, 3}, 3, map[string]int{
			"0102": 2, "0203": 2, "0301": 1,
			"010203": 2, "020301": 1, "030102": 1,
		}},
		{[]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}, 4, map[string]int{
			"0102": 3, "0203": 3, "0304": 3, "0401": 2,
			"010203": 3, "020304": 3, "030401": 2, "040102": 2,
			"01020304": 3, "02030401": 2, "03040102": 2, "04010203": 2,
		}},
	}
	for _, x := range tt {
		m := BuildHistogram(x.data, x.max)
		assert.Equal(t, x.exp, m)
	}
}

func Test_reduceSubCounts(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	ps := []Patt{
		{9, []byte{1, 2}, "0102"},      // {1, 2} is 1 times in each of 3 {1, 2, 3}, and 2 times in one {1, 2, 3, 1, 2, 3}
		{3, []byte{1, 2, 3}, "010203"}, // {1, 2, 3} is 2 times in one {1, 2, 3, 1, 2, 3}
		{1, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
	}
	exp := []Patt{
		{4, []byte{1, 2}, "0102"},      // 9-3-2 = 4
		{1, []byte{1, 2, 3}, "010203"}, // 3-2 = 1
		{1, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
	}
	act := reduceSubCounts(ps)
	assert.Equal(t, exp, act)
}

func Test_scanForRepetitions(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	tt := []struct {
		d []byte         // data
		l int            // ptLen
		m map[string]int // expected map
	}{
		{[]byte{1, 2, 3, 1, 2, 3}, 2, map[string]int{"0102": 2, "0203": 2, "0301": 1}},
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
		m := scanForRepetitions(x.d, x.l)
		assert.Equal(t, x.m, m)
	}
}

func Test_histogramToList(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	tt := []struct {
		m   map[string]int // given map
		exp []Patt         // expected list
	}{
		{map[string]int{"0102": 2, "0203": 2, "0301": 1}, []Patt{{2, []byte{1, 2}, "0102"}, {2, []byte{2, 3}, "0203"}, {1, []byte{3, 1}, "0301"}}},
		{map[string]int{"0102": 4, "0808": 7}, []Patt{{7, []byte{8, 8}, "0808"}, {4, []byte{1, 2}, "0102"}}},
	}

	for _, x := range tt {
		result := histogramToList(x.m)
		act := SortByDescentingCountAndLengthAndAphabetical(result)
		assert.Equal(t, x.exp, act)
	}
}

func Test_sortByIncreasingLength(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	pat := []Patt{
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
		{900, []byte{1, 2}, "0102"},
		{300, []byte{1, 2, 3}, "010203"},
	}
	exp := []Patt{
		{900, []byte{1, 2}, "0102"},
		{300, []byte{1, 2, 3}, "010203"},
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
	}
	act := SortByIncreasingLengthAndAlphabetical(pat)
	assert.Equal(t, exp, act)
}

func Test_SortByDescentingCountAndLengthAndAphabetical(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	pat := []Patt{
		{100, []byte{1, 2, 3, 1, 2, 3, 4}, "01020301020304"},
		{100, []byte{1, 2, 3, 4}, "01020304"},
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
		{900, []byte{1, 2}, "0102"},
		{100, []byte{8, 2, 3, 1, 2, 3}, "080203010203"},
		{300, []byte{1, 2, 3}, "010203"},
	}
	exp := []Patt{
		{900, []byte{1, 2}, "0102"},
		{300, []byte{1, 2, 3}, "010203"},
		{100, []byte{1, 2, 3, 1, 2, 3, 4}, "01020301020304"},
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
		{100, []byte{8, 2, 3, 1, 2, 3}, "080203010203"},
		{100, []byte{1, 2, 3, 4}, "01020304"},
	}
	act := SortByDescentingCountAndLengthAndAphabetical(pat)
	assert.Equal(t, exp, act)
}

func Test_extendHistorgamMap(t *testing.T) {
	tt := []struct {
		dst, src, exp map[string]int // expected map
	}{
		{
			map[string]int{"0102": 5},
			map[string]int{"0301": 2},
			map[string]int{"0102": 5, "0301": 2},
		},
		{
			map[string]int{"0102": 5},
			map[string]int{"0102": 2},
			map[string]int{"0102": 7},
		},
		{
			map[string]int{},
			map[string]int{"0102": 2},
			map[string]int{"0102": 2},
		},
		{
			map[string]int{"0102": 2},
			map[string]int{},
			map[string]int{"0102": 2},
		},
		{
			map[string]int{},
			map[string]int{},
			map[string]int{},
		},
	}

	for i, x := range tt {
		act := make(map[string]int)
		maps.Copy(act, x.dst)
		extendHistorgamMap(act, x.src)

		assert.Equal(t, tt[i].exp, act)
	}

	fmt.Println(tt)
}

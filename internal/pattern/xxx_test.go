package pattern


/*
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
*/

/*

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
			//fmt.Println("exp:", reflect.ValueOf(e).Type(), e)
			//fmt.Println("act:", reflect.ValueOf(a).Type(), a)
			result := reflect.DeepEqual(e, a)
			assert.True(t, result)
		})
	}
}

func TestHistogram_SortKeysByDescSize(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]int
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		exp    fields
	}{
		// TODO: Add test cases.
		{
			"", // name
			fields{map[string]int{}, &mu, []string{"bb11", "112233", "aa22"}},
			fields{map[string]int{}, &mu, []string{"112233", "aa22", "bb11"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.SortKeysByDescSize()
			for i := range p.Key {
				assert.Equal(t, tt.exp.Keys[i], tt.fields.Keys[i])
			}
		})
	}
}

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


func _XXX_Test_reduceSubCounts(t *testing.T) {
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

func _XXX_Test_scan_ForRepetitions(t *testing.T) {
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
		m := scan_ForRepetitions(x.d, x.l)
		assert.Equal(t, x.m, m)
	}
}

func _XXX_Test_histogramToList(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	tt := []struct {
		m   map[string]int // given map
		exp []Patt         // expected list
	}{
		{map[string]int{"0102": 2, "0203": 2, "0301": 1}, []Patt{{2, []byte{1, 2}, "0102"}, {2, []byte{2, 3}, "0203"}, {1, []byte{3, 1}, "0301"}}},
		{map[string]int{"0102": 4, "0808": 7}, []Patt{{7, []byte{8, 8}, "0808"}, {4, []byte{1, 2}, "0102"}}},
	}

	for _, x := range tt {
		result := HistogramToList(x.m)
		act := SortByDescCountDescLength(result)
		assert.Equal(t, x.exp, act)
	}
}

func _XXX_Test_sortByIncreasingLength(t *testing.T) {
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
	act := SortByIncLength(pat)
	assert.Equal(t, exp, act)
}

func _XXX_Test_extendHistorgamMap(t *testing.T) {
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

func TestScanData(Histogram(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	type args struct {
		hist map[string]int
		data []byte
		max  int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"1st", args{nil, nil, 4}},
		{"2nd", args{map[string]int{"0102": 3}, []byte{0x01, 0x02}, 4}},
	}
	for _, tt := range tests {
		ScanData(Histogram(tt.args.hist, tt.args.data, tt.args.max)
	}
}

*/

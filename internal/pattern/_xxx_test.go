package pattern

/*
func TestHistogram_Reduce_2(t *testing.T) {
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
		{ // name fields                               data    01234567  max, ring, reduce,
			"", fields{map[string]Pattern{}, &m, nil}, []byte("xxAAAAyy"), 4, false, false,
			map[string]Pattern{ // expected
				s2h("xxAA"): {Pos: []int{0}},    //// xxAA     @ 0
				s2h("xAAA"): {Pos: []int{1}},    ////  xAAA    @ 1       <- when xAAA is big key, xAA@1 is removed but AA@2 needs to be restored then
				s2h("AAAA"): {Pos: []int{2}},    ////   AAAA   @ 2
				s2h("AAAy"): {Pos: []int{3}},    ////    AAAy  @ 3       <- when AAAy is big key, AAy@4 is removed but AA@4 needs to be restored then
				s2h("AAyy"): {Pos: []int{4}},    ////     AAyy @ 4
				s2h("xxA"):  {Pos: []int{0}},    //// xxA      @ 0
				s2h("xAA"):  {Pos: []int{1}},    ////  xAA     @ 1        <- when xAA is big key AA@2 is removed
				s2h("AAA"):  {Pos: []int{2, 3}}, ////   AAA    @ 2 = 2 3
				//                                 //    AAA   @ 3 = 2 3
				s2h("AAy"): {Pos: []int{4}},       //     AAy  @ 4        <- when AAy is big key AA@4 is removed
				s2h("Ayy"): {Pos: []int{5}},       //      Ayy @ 5
				s2h("xx"):  {Pos: []int{0}},       // xx       @ 0
				s2h("xA"):  {Pos: []int{1}},       //  xA      @ 1
				s2h("AA"):  {Pos: []int{2, 3, 4}}, //   AA     @ 2 = 2 3 4
				s2h("Ay"):  {Pos: []int{5}},       //      Ay  @ 5
				s2h("yy"):  {Pos: []int{6}},       //       yy @ 6
			},
		},

		//  // WITHOUT RESTORE
		//  { // name fields                               data    01234567  max, ring, reduce,
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("xxAAAAyy"), 4, false, true,
		//  	map[string]Pattern{ // expected                              - removes -
		//  		s2h("xxAA"): {Pos: []int{0}}, // xxAA     @ 0            xxA@0 xAA@1
		//  		s2h("xAAA"): {Pos: []int{1}}, //  xAAA    @ 1            xAA@1 AAA@2                                 <- when xAAA is big key, xAA@1 is removed but AA@2 needs to be restored then
		//  		s2h("AAAA"): {Pos: []int{2}}, //   AAAA   @ 2            AAA@2 AAA@3
		//  		s2h("AAAy"): {Pos: []int{3}}, //    AAAy  @ 3            AAA@3 AAy@4                                 <- when AAAy is big key, AAy@4 is removed but AA@4 needs to be restored then
		//  		s2h("AAyy"): {Pos: []int{4}}, //     AAyy @ 4            AAy@4 Ayy@6
		//  		s2h("xxA"):  {Pos: []int{}},  // xxA      @ 0            xx@0  xA@1
		//  		s2h("xAA"):  {Pos: []int{}},  //  xAA     @ 1            xA@1  AA@2                                  <- when xAA is big key AA@2 is removed
		//  		s2h("AAA"):  {Pos: []int{}},  //   AAA    @ 2 = 2 3      AA@2  AA@3  AA@4
		//  		//                            //    AAA   @ 3 = 2 3
		//  		s2h("AAy"): {Pos: []int{}}, ////     AAy  @ 4            AA@4  Ay@5                                  <- when AAy is big key AA@4 is removed
		//  		s2h("Ayy"): {Pos: []int{}}, ////      Ayy @ 5            Ay@5  yy@6
		//  		s2h("xx"):  {Pos: []int{}}, //// xx       @ 0
		//  		s2h("xA"):  {Pos: []int{}}, ////  xA      @ 1
		//  		s2h("AA"):  {Pos: []int{}}, ////   AA     @ 2 = 2 3 4
		//  		s2h("Ay"):  {Pos: []int{}}, ////      Ay  @ 5
		//  		s2h("yy"):  {Pos: []int{}}, ////       yy @ 6
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
				//p.ReduceFromSmallerSide()
			}

			for k, v := range tt.exp {
				fmt.Println(k, v.Pos, tt.fields.Hist[k].Pos)
			}
			for k, v := range tt.exp {
				assert.Equal(t, v.Pos, tt.fields.Hist[k].Pos)
			}
		})
	}
}
*/

/*
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
*/
/*
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
*/
/*
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
*/
/*
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
*/

/*

func withinTolerance(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	return (d / math.Abs(b)) < epsilon
}

*/

/*
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
*/

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

/*

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

func TestHistogram_ReduceOverlappingKeys(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]Pat
		mu   *sync.Mutex
		Keys []string
	}
	type args struct {
		equalSize1stKey []string
		equalSize2ndKey []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    map[string]Pat
	}{ // test cases:
		{
			"",
			fields{map[string]Pat{"1a1a": {3, []int{0, 1, 5}}, "1a1a1a": {1, []int{0}}}, &mu, nil}, // the histogram in p
			args{[]string{"1a1a1a"}, []string{"1a1a"}},                                             // the function arguments
			map[string]Pat{"1a1a": {1, []int{5}}, "1a1a1a": {1, []int{0}}},                         // the expected result in p
		},
		{
			"", // case: |xx1a1a1a1axx...|
			fields{map[string]Pat{"1a1a": {9, []int{0, 2, 3, 4, 5, 6, 8, 10, 20}}, "1a1a1a1a": {2, []int{2, 32}}}, &mu, nil}, // the histograms in p
			args{[]string{"1a1a1a1a"}, []string{"1a1a"}},                                          // the function arguments
			map[string]Pat{"1a1a": {6, []int{0, 5, 6, 8, 10, 20}}, "1a1a1a1a": {2, []int{2, 32}}}, // the expected result in p
		},
		{
			"",
			fields{map[string]Pat{"1122": {2, []int{8, 32}}, "112233": {1, []int{8}}}, &mu, nil},
			args{[]string{"112233"}, []string{"1122"}},
			map[string]Pat{"1122": {1, []int{32}}, "112233": {1, []int{8}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.ReduceOverlappingKeys(tt.args.equalSize1stKey, tt.args.equalSize2ndKey)
			p.SortPositions()
			e := tt.exp
			a := p.Hist
			fmt.Println("exp:", reflect.ValueOf(e).Type(), e)
			fmt.Println("act:", reflect.ValueOf(a).Type(), a)
			result := reflect.DeepEqual(e, a)
			assert.True(t, result)
			assert.Equal(t, tt.exp, p.Hist)
		})
	}
}


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

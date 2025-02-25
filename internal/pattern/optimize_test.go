package pattern

/*
func _TestHistogram_Reduce(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]int
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		exp    map[string]int
	}{
		// TODO: Add test cases.
		{
			"",
			fields{map[string]int{"1122": 10, "112233": 1}, &mu, nil},
			map[string]int{"1122": 9, "112233": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Key:  tt.fields.Keys,
			}
			p.GetKeys()
			p.SortKeysByDescSize()
			p.Reduce()
			assert.Equal(t, tt.exp, tt.fields.Hist)
		})
	}
}
*/

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
	}{
		// TODO: Add test cases.
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

func TestHistogram_ReduceOverlappingKeys(t *testing.T) {
	var mu sync.Mutex
	type fields struct {
		Hist map[string]int
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
		exp    map[string]int
	}{
		// TODO: Add test cases.
		{
			"",
			fields{map[string]int{"1111": 20, "11111111": 2}, &mu, nil},
			args{[]string{"11111111"}, []string{"1111"}},
			map[string]int{"1111": 14, "11111111": 2},
		},
		{
			"",
			fields{map[string]int{"1122": 10, "112233": 1}, &mu, nil},
			args{[]string{"112233"}, []string{"1122"}},
			map[string]int{"1122": 9, "112233": 1},
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
			assert.Equal(t, tt.exp, p.Hist)
		})
	}
}
*/

// generated: ////////////////////////////////

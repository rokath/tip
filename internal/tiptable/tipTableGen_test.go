package tiptable

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

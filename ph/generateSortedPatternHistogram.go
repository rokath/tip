package main

import "slices"

// generateB1Histogram scans data for 1-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB1Histogram(data []byte) map[byte]int {
	m := make(map[byte]int)
	for _, x := range data {
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateB2Histogram scans data for 2-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB2Histogram(data []byte) map[[2]byte]int {
	m := make(map[[2]byte]int)
	for i := 0; i < len(data)-1; i++ {
		x := [2]byte{data[i], data[i+1]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateB3Histogram scans data for 3-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB3Histogram(data []byte) map[[3]byte]int {
	m := make(map[[3]byte]int)
	for i := 0; i < len(data)-2; i++ {
		x := [3]byte{data[i], data[i+1], data[i+2]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateB4Histogram scans data for 4-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB4Histogram(data []byte) map[[4]byte]int {
	m := make(map[[4]byte]int)
	for i := 0; i < len(data)-3; i++ {
		x := [4]byte{data[i], data[i+1], data[i+2], data[i+3]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateB5Histogram scans data for 5-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB5Histogram(data []byte) map[[5]byte]int {
	m := make(map[[5]byte]int)
	for i := 0; i < len(data)-4; i++ {
		x := [5]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateB6Histogram scans data for 6-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB6Histogram(data []byte) map[[6]byte]int {
	m := make(map[[6]byte]int)
	for i := 0; i < len(data)-5; i++ {
		x := [6]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateB7Histogram scans data for 7-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB7Histogram(data []byte) map[[7]byte]int {
	m := make(map[[7]byte]int)
	for i := 0; i < len(data)-6; i++ {
		x := [7]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5], data[i+6]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateB8Histogram scans data for 8-byte repetitions, stores them as keys in m with the occurances count as values.
func generateB8Histogram(data []byte) map[[8]byte]int {
	m := make(map[[8]byte]int)
	for i := 0; i < len(data)-7; i++ {
		x := [8]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5], data[i+6], data[i+7]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	return m
}

// generateSortedPatternHistogram searches data for any 1-8 bytes pattern
// and returns them with their count in a by count sorted list.
func generateSortedPatternHistogram(data []byte) []nPatt {
	m1 := generateB1Histogram(data)
	m2 := generateB2Histogram(data)
	m3 := generateB3Histogram(data)
	m4 := generateB4Histogram(data)
	m5 := generateB5Histogram(data)
	m6 := generateB6Histogram(data)
	m7 := generateB7Histogram(data)
	m8 := generateB8Histogram(data)
	// Maps m1...m8 contain the pattern histograms.

	pn := make([]nPatt, 0, 1024)

	for k, v := range m1 {
		pn = append(pn, nPatt{v, []byte{k}})
	}
	for k, v := range m2 {
		pn = append(pn, nPatt{v, k[:]})
	}
	for k, v := range m3 {
		pn = append(pn, nPatt{v, k[:]})
	}
	for k, v := range m4 {
		pn = append(pn, nPatt{v, k[:]})
	}
	for k, v := range m5 {
		pn = append(pn, nPatt{v, k[:]})
	}
	for k, v := range m6 {
		pn = append(pn, nPatt{v, k[:]})
	}
	for k, v := range m7 {
		pn = append(pn, nPatt{v, k[:]})
	}
	for k, v := range m8 {
		pn = append(pn, nPatt{v, k[:]})
	}

	// sort pn for count and pattern length
	compareFn := func(a, b nPatt) int {
		if a.n > b.n {
			return 1
		}
		if a.n < b.n {
			return -1
		}
		if len(a.pattern) > len(b.pattern) {
			return 1
		}
		if len(a.pattern) < len(b.pattern) {
			return -1
		}
		return 0
	}
	slices.SortFunc(pn, compareFn)
	return pn
}

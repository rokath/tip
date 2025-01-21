package main

import "sort"

// createB3Histogram scans data for 1-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB1Histogram(data []byte) (m map[byte]int, keys []byte) {
	// Create a histogram for 1-byte sequences.
	// The keys are the single bytes and the values are their occurance count.
	m = make(map[byte]int)
	for _, x := range data {
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort keys according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

// createB2Histogram scans data for 2-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB2Histogram(data []byte) (m map[[2]byte]int, keys [][2]byte) {
	// Create a histogram for 2-byte sequences.
	// The keys are the 2-bytes sequences and the values are their occurance count.
	m = make(map[[2]byte]int)
	for i := 0; i < len(data)-1; i++ {
		x := [2]byte{data[i], data[i+1]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([][2]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

// createB3Histogram scans data for 3-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB3Histogram(data []byte) (m map[[3]byte]int, keys [][3]byte) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m = make(map[[3]byte]int)
	for i := 0; i < len(data)-2; i++ {
		x := [3]byte{data[i], data[i+1], data[i+2]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([][3]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

// createB4Histogram scans data for 3-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB4Histogram(data []byte) (m map[[4]byte]int, keys [][4]byte) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m = make(map[[4]byte]int)
	for i := 0; i < len(data)-3; i++ {
		x := [4]byte{data[i], data[i+1], data[i+2], data[i+3]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([][4]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

// createB5Histogram scans data for 5-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB5Histogram(data []byte) (m map[[5]byte]int, keys [][5]byte) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m = make(map[[5]byte]int)
	for i := 0; i < len(data)-4; i++ {
		x := [5]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([][5]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

// createB6Histogram scans data for 6-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB6Histogram(data []byte) (m map[[6]byte]int, keys [][6]byte) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m = make(map[[6]byte]int)
	for i := 0; i < len(data)-5; i++ {
		x := [6]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([][6]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

// createB7Histogram scans data for 7-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB7Histogram(data []byte) (m map[[7]byte]int, keys [][7]byte) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m = make(map[[7]byte]int)
	for i := 0; i < len(data)-6; i++ {
		x := [7]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5], data[i+6]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([][7]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

// createB8Histogram scans data for 8-byte repetitions, stores them as values in m with the occurances count as values
// and returns them as a sorted list with max count first.
func createB8Histogram(data []byte) (m map[[8]byte]int, keys [][8]byte) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m = make(map[[8]byte]int)
	for i := 0; i < len(data)-7; i++ {
		x := [8]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5], data[i+6], data[i+7]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}
	// Get a list of all m keys.
	keys = make([][8]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})
	return
}

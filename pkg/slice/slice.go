package slice

import (
	"slices"
)

// Index returns at which index v was first found in s or -1.
func Index(s, v []byte) int {
	if len(v) > len(s) {
		return -1
	}
	for i := 0; i < len(s)-len(v)+1; i++ {
		if slices.Equal(s[i:i+len(v)], v) {
			return i
		}
	}
	return -1
}

// Count returns how often v was found in s.
// If v has len 0 the value 0 is returned.
// After a match, the seach space is reduced to to position after the match.
// Example: ff ff ff and ff ff returns 1 and not 2.
func Count(s, v []byte) (n int) {
	if len(v) == 0 {
		return
	}
	limit := len(s) - len(v) + 1
	for range limit {
		idx := Index(s, v)
		if idx < 0 {
			return
		}
		n++
		s = s[idx+len(v):]
	}
	return
}

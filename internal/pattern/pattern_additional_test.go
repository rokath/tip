package pattern

import (
	"path/filepath"
	"slices"
	"sync"
	"testing"

	"github.com/spf13/afero"
	"github.com/tj/assert"
)

func TestNewHistogram(t *testing.T) {
	var mu sync.Mutex

	h := NewHistogram(&mu)

	assert.NotNil(t, h)
	assert.NotNil(t, h.Hist)
	assert.Equal(t, 0, len(h.Hist))
	assert.Same(t, &mu, h.mu)
	assert.Nil(t, h.Keys)
}

func TestHistogramDeleteEmptyKeys(t *testing.T) {
	var mu sync.Mutex
	h := &Histogram{
		Hist: map[string]Pattern{
			s2h("ab"): {Bytes: []byte("ab"), Pos: []int{1}},
			s2h("cd"): {Bytes: []byte("cd"), Pos: nil},
			s2h("ef"): {Bytes: []byte("ef"), Pos: []int{}},
		},
		mu: &mu,
	}

	h.DeleteEmptyKeys()

	assert.Equal(t, map[string]Pattern{
		s2h("ab"): {Bytes: []byte("ab"), Pos: []int{1}},
	}, h.Hist)
}

func TestHistogramScanFile(t *testing.T) {
	var mu sync.Mutex
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	assert.NoError(t, afs.WriteFile("sample.bin", []byte("ababa"), 0o644))

	h := NewHistogram(&mu)
	err := h.ScanFile(afs, "sample.bin", 3)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{0, 2}, h.Hist[s2h("aba")].Pos)
	assert.ElementsMatch(t, []int{0, 2}, h.Hist[s2h("ab")].Pos)
	assert.ElementsMatch(t, []int{1, 3}, h.Hist[s2h("ba")].Pos)
}

func TestHistogramScanFileReadError(t *testing.T) {
	var mu sync.Mutex
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}
	h := NewHistogram(&mu)

	err := h.ScanFile(afs, "missing.bin", 3)

	assert.Error(t, err)
}

func TestHistogramScanAllFiles(t *testing.T) {
	var mu sync.Mutex
	root := t.TempDir()
	fileA := filepath.Join(root, "a.bin")
	fileB := filepath.Join(root, "nested", "b.bin")
	fs := afero.NewOsFs()
	afs := &afero.Afero{Fs: fs}

	assert.NoError(t, afs.MkdirAll(filepath.Dir(fileB), 0o755))
	assert.NoError(t, afs.WriteFile(fileA, []byte("abca"), 0o644))
	assert.NoError(t, afs.WriteFile(fileB, []byte("bcab"), 0o644))

	h := NewHistogram(&mu)
	err := h.ScanAllFiles(afs, root, 2)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{0, 2}, h.Hist[s2h("ab")].Pos)
	assert.ElementsMatch(t, []int{1, 0}, h.Hist[s2h("bc")].Pos)
	assert.ElementsMatch(t, []int{2, 1}, h.Hist[s2h("ca")].Pos)
}

func TestSortByDescWeight(t *testing.T) {
	list := []Pattern{
		{Bytes: []byte("aaaa"), Pos: []int{0}},
		{Bytes: []byte("bb"), Pos: []int{0, 2, 4}},
		{Bytes: []byte("ccc"), Pos: []int{1, 3, 5}},
		{Bytes: []byte("d"), Pos: []int{7, 8}},
	}

	got := SortByDescWeight(list)
	gotKeys := make([]string, len(got))
	for i, pattern := range got {
		gotKeys[i] = string(pattern.Bytes)
	}

	assert.Equal(t, []string{"ccc", "bb", "aaaa", "d"}, gotKeys)
	assert.True(t, slices.Equal(got[0].Pos, []int{1, 3, 5}))
}

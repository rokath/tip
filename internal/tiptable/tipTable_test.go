package tiptable

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/rokath/tip/internal/pattern"
	"github.com/spf13/afero"
	"github.com/tj/assert"
)

func TestSpaces(t *testing.T) {
	tt := []struct {
		l int
		s string
	}{
		{-2, ""},
		{-1, ""},
		{0, ""},
		{1, " "},
		{2, "  "},
	}
	for _, x := range tt {
		assert.Equal(t, x.s, spaces(x.l))
	}
}

func TestByteSliceAsASCII(t *testing.T) {
	got := byteSliceAsASCII([]byte{'A', 0x00, 'z'}, 5)
	assert.Equal(t, "`A˙z`  ", got)
}

func TestCreatePatternLineString(t *testing.T) {
	got := createPatternLineString([]byte{'A', 0x00}, 4)
	assert.Equal(t, "  2, 0x41, 0x00,             // `A˙`  ", got)
}

func TestTipPackageIDs(t *testing.T) {
	oldID1Count := ID1Count
	t.Cleanup(func() {
		ID1Count = oldID1Count
	})
	ID1Count = 3

	id1, id2 := tipPackageIDs(2)
	assert.Equal(t, uint8(2), id1)
	assert.Equal(t, -1, id2)

	id1, id2 = tipPackageIDs(4)
	assert.Equal(t, uint8(4), id1)
	assert.Equal(t, 1, id2)

	id1, id2 = tipPackageIDs(259)
	assert.Equal(t, uint8(5), id1)
	assert.Equal(t, 1, id2)
}

func TestPrintPattern(t *testing.T) {
	out := captureStdout(t, func() {
		PrintPattern(7, pattern.Pattern{
			Bytes: []byte{'A', 0x00, 'B'},
			Pos:   []int{1, 2},
		})
	})

	assert.Contains(t, out, "i:  7")
	assert.Contains(t, out, "weight:       6")
	assert.Contains(t, out, "cnt:     2")
	assert.Contains(t, out, "ascii:'")
}

func TestGenerateFromFile(t *testing.T) {
	oldBits := UnreplacableContainerBits
	oldID1Count := ID1Count
	oldMaxID := MaxID
	oldID1Max := ID1Max
	oldPatternSizeMax := pattern.PatternSizeMax
	oldVerbose := Verbose
	t.Cleanup(func() {
		UnreplacableContainerBits = oldBits
		ID1Count = oldID1Count
		MaxID = oldMaxID
		ID1Max = oldID1Max
		pattern.PatternSizeMax = oldPatternSizeMax
		Verbose = oldVerbose
	})

	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}
	assert.NoError(t, afs.WriteFile("sample.bin", []byte("ababa"), 0o644))

	UnreplacableContainerBits = 6
	ID1Count = 2
	pattern.PatternSizeMax = 3
	Verbose = false

	err := Generate(afs, "idTable.c", "sample.bin", 3)

	assert.NoError(t, err)
	content, err := afs.ReadFile("idTable.c")
	assert.NoError(t, err)
	text := string(content)
	assert.Contains(t, text, "unsigned unreplacableContainerBits = 6;")
	assert.Contains(t, text, "unsigned ID1Max = 191;")
	assert.Contains(t, text, "unsigned ID1Count = 2;")
	assert.Contains(t, text, "unsigned MaxID = 48197;")
	assert.Contains(t, text, "unsigned LastID = 3;")
	assert.Contains(t, text, "uint8_t maxPatternlength = 3;")
	assert.Contains(t, text, "static uint8_t const idTable[] = { // from sample.bin")
	assert.Contains(t, text, "0 // table end marker")
}

func TestGenerateFromDirectoryVerbose(t *testing.T) {
	oldBits := UnreplacableContainerBits
	oldID1Count := ID1Count
	oldMaxID := MaxID
	oldID1Max := ID1Max
	oldPatternSizeMax := pattern.PatternSizeMax
	oldVerbose := Verbose
	t.Cleanup(func() {
		UnreplacableContainerBits = oldBits
		ID1Count = oldID1Count
		MaxID = oldMaxID
		ID1Max = oldID1Max
		pattern.PatternSizeMax = oldPatternSizeMax
		Verbose = oldVerbose
	})

	root := t.TempDir()
	fs := afero.NewOsFs()
	afs := &afero.Afero{Fs: fs}
	assert.NoError(t, afs.WriteFile(filepath.Join(root, "a.bin"), []byte("abca"), 0o644))
	assert.NoError(t, afs.WriteFile(filepath.Join(root, "b.bin"), []byte("bcab"), 0o644))

	UnreplacableContainerBits = 7
	ID1Count = 1
	pattern.PatternSizeMax = 3
	Verbose = true

	outFile := filepath.Join(root, "idTable.c")
	err := Generate(afs, outFile, root, 3)

	assert.NoError(t, err)
	content, err := afs.ReadFile(outFile)
	assert.NoError(t, err)
	text := string(content)
	assert.Contains(t, text, "unsigned unreplacableContainerBits = 7;")
	assert.Contains(t, text, "unsigned ID1Max = 127;")
	assert.Contains(t, text, "// Informal, here are all by weight sorted pattern occuring at least twice:")
	assert.Contains(t, text, "IDTable points to the used idTable")
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	fn()

	assert.NoError(t, w.Close())
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	assert.NoError(t, err)
	assert.NoError(t, r.Close())
	return buf.String()
}

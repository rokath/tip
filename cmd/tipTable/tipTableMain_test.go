package main

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/tj/assert"
)

var FSys *afero.Afero // ram file system for the tests

func init() {
	// All id tests should be executed only on a memory mapped file system.
	FSys = &afero.Afero{Fs: afero.NewMemMapFs()}
}

func Test_doit(t *testing.T) {
	help = true
	var b bytes.Buffer
	doit(&b, FSys)
	act := b.String()
	exp := `Usage: tipTable -i inputFileName [-o outputFileName] [-z max pattern size] [-v]
Example: `+"`tipTableGen -i trice.bin`"+` creates tipTable.c
The TipUserManual explains details.
`
	assert.Equal(t, exp, act)
}

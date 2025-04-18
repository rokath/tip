package pattern

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
)

var (
	B      bytes.Buffer
	W      = io.Writer(&B) // An io.Writer is needed for some tests
	FSys   *afero.Afero    // ram file system for the tests
	global defaults        // global holds global vars default values
)

func init() {
	// All tests should be executed only on a memory mapped file system.
	FSys = &afero.Afero{Fs: afero.NewMemMapFs()}
}

// Setup should be called on the begin of each id test function, if global variables are used/changed.
func Setup(t *testing.T) func() {
	// setup code here
	global.RestoreVars(t)
	fmt.Println(t.Name(), "...")

	return func() {
		// tear-down code here
		fmt.Println(t.Name(), "...done.")
		B.Reset() // Clear output generated during the test.
	}
}

// TestMain - see for example https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc
//
// In Go, each package generates an individual test binary and they are tested parallel.
// All package tests are executed sequentially but use the same global variables.
// Therefore we have to reset the global variables in each test function.
func TestMain(m *testing.M) {
	global.StoreVars() // Do stuff BEFORE the id package tests!
	exitVal := m.Run() // Run the tests sequentially in alphabetical order.
	os.Exit(exitVal)   // Do stuff AFTER the id package tests!
}

type defaults struct {
	verbose bool
	patternSizeMax int
}

// StoreVars reads global variables for restauration later.
func (p *defaults) StoreVars() {
	p.verbose = Verbose
	p.patternSizeMax = PatternSizeMax
}

// RestoreVars sets all global variables into previous state.
func (p *defaults) RestoreVars(t *testing.T) {
	Verbose = p.verbose
	PatternSizeMax = p.patternSizeMax
}

package bar

import (
	"github.com/superhawk610/terminal"
)

// Output is a stand-in for a basic io.Writer that also exposes
// a ClearLine() function to clear the current line and return the
// cursor to the first index
type Output interface {
	ClearLine()
	Printf(format string, vals ...interface{})
}

type stdout struct{}

// ClearLine clears the current output line and returns the cursor
// to the first index
func (stdout) ClearLine() {
	terminal.ClearLine()
}

// Printf accepts a format string and any number of input values
func (stdout) Printf(format string, vals ...interface{}) {
	terminal.Overwritef(format, vals...)
}

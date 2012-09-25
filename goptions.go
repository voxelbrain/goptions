package goptions

import (
	"fmt"
	"io"
	"os"
)

var (
	globalFlagSet *FlagSet
)

func Parse(v interface{}) error {
	fs, err := NewFlagSet(v)
	if err != nil {
		return err
	}
	globalFlagSet = fs

	e := fs.Parse(os.Environ()[1:])
	if e != nil {
		return e
	}

	return nil
}

func PrintHelp() {
	if globalFlagSet == nil {
		panic("Must call Parse() before PrintHelp()")
	}
	globalFlagSet.PrintHelp(os.Stdout)
}

func DefaultHelpFunc(w io.Writer, fs *FlagSet) {
	fmt.Fprintf(w, "wat?")
}

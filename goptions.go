/*
package goptions implements a flexible parser for command line options.

Key targets were the support for both long and short flag versions, mutually
exclusive flags, and verbs. Flags and their corresponding variables are defined
by the tags in a (possibly anonymous) struct.

    var options struct {
    	Name string `goptions:"-n, --name"`
    	Force bool `goptions:"-f, --force"`
    	Verbosity int `goptions:"-v, --verbose, accumulate"`
    }

Short flags can be combined (e.g. `-nfv`). Long flags take their value after a
separating space. The equals notation (`--long-flag=value`) is NOT supported
right now.

Every member of the struct, which is supposed to catch a command line value
has to have a "goptions" tag. Multiple short and long flag names can be specified.
Each tag can also list any number of the following options:

    accumulate        - (Only valid for `int`) Counts how of then the flag has been
                        specified in the short version. The long version simply
                        accepts an int.
    obligatory        - Flag must be specified. Otherwise an error will be returned
                        when Parse() is called.
    description='...' - Set the description for this particular flag. Will be
                        used by the HelpFunc.
    mutexgroup='...'  - Sets the name of the MutexGroup. Only one flag of the
                        ones sharing a MutexGroup can be set. Otherwise an error
                        will be returned when Parse() is called. If one flag in a
                        MutexGroup is `obligatory` one flag of the group must be
                        specified.

goptions also has support for verbs. Each verb accepts its own set of flags which
take exactly the same tag format as global options. For an usage example of verbs
see the PrintHelp() example.
*/
package goptions

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	VERSION = "1.3.4"
)

var (
	globalFlagSet *FlagSet
)

// ParseAndFail is a convenience function to parse tos.Args[1:] and print
// the help if an error occurs. This should cover 90% of this library's
// applications.
func ParseAndFail(v interface{}) {
	err := Parse(v)
	if err != nil {
		errCode := 0
		if err != ErrHelpRequest {
			errCode = 1
			fmt.Printf("Error: %s\n", err)
		}
		PrintHelp()
		os.Exit(errCode)
	}
}

// Parse parses the command-line flags from os.Args[1:].
func Parse(v interface{}) error {
	globalFlagSet = NewFlagSet(filepath.Base(os.Args[0]), v)
	return globalFlagSet.Parse(os.Args[1:])
}

// PrintHelp renders the default help to os.Stderr.
func PrintHelp() {
	if globalFlagSet == nil {
		panic("Must call Parse() before PrintHelp()")
	}
	globalFlagSet.PrintHelp(os.Stderr)
}

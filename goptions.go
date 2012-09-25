package goptions

import (
	"os"
)

func Parse(v interface{}) error {
	fs, err := NewFlagSet(v)
	if err != nil {
		return err
	}

	e := fs.Parse(os.Environ()[1:])
	if e != nil {
		return e
	}

	return nil
}

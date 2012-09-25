package goptions

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrHelpRequest = errors.New("Request for Help")
)

// Catches remaining strings after parsing has finished
type Remainder []string

func Parse(args []string, v interface{}) error {
	structValue := reflect.ValueOf(v)
	if structValue.Kind() != reflect.Ptr {
		panic("Value type is not a pointer to a struct")
	}
	structValue = structValue.Elem()
	if structValue.Kind() != reflect.Struct {
		panic("Value type is not a pointer to a struct")
	}

	fs, err := newFlagSet(structValue)
	if err != nil {
		return err
	}

	e := fs.Parse(args)
	if e != nil {
		return e
	}

	if fs.helpFlag != nil && fs.helpFlag.WasSpecified {
		return ErrHelpRequest
	}

	// Check for unset, obligatory flags
	for _, f := range fs.flags {
		if f.Obligatory && !f.WasSpecified {
			return fmt.Errorf("%s must be specified", f.Name())
		}
	}

	// Check for multiple set flags in one mutex group
	mgs := fs.MutexGroups()
	for _, mg := range mgs {
		wasSpecifiedCount := 0
		names := make([]string, 0)
		for _, flag := range mg {
			names = append(names, flag.Name())
			if flag.WasSpecified {
				wasSpecifiedCount += 1
			}
		}
		if wasSpecifiedCount >= 2 {
			return fmt.Errorf("Only one of %s can be specified", strings.Join(names, ","))
		}
	}

	return nil
}

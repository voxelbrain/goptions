package goptions

import (
	"errors"
	"reflect"
)

// Flag represents a single flag of a FlagSet.
type Flag struct {
	Short        []string
	Long         []string
	MutexGroup   string
	Accumulate   bool
	Description  string
	Obligatory   bool
	WasSpecified bool
	value        reflect.Value
}

// Return the name of the flag preceding the right amount of dashes.
// The long name is preferred. If no name has been specified, "<unspecified>"
// will be returned.
func (f *Flag) Name() string {
	if len(f.Long) > 0 {
		return "--" + f.Long[0]
	}
	if len(f.Short) > 0 {
		return "-" + f.Short[0]
	}
	return "<unspecified>"
}

// Returns true if the flag expects a separate value on the command line.
func (f *Flag) NeedsExtraValue() bool {
	// Explicit over implicit
	if f.value.Kind() == reflect.Bool {
		return false
	}
	if f.value.Kind() == reflect.Int && f.Accumulate {
		return false
	}
	return true
}

func (f *Flag) set() {
	f.WasSpecified = true
	if f.value.Kind() == reflect.Bool {
		f.setValue(true)
	} else if f.value.Kind() == reflect.Int && f.Accumulate {
		f.setValue(f.value.Interface().(int) + 1)
	}
}

func (f *Flag) setValue(v interface{}) (err error) {
	defer func() {
		if x := recover(); x != nil {
			if str, ok := x.(string); ok {
				err = errors.New(str)
				return
			}
			err = x.(error)
		}
	}()
	f.value.Set(reflect.ValueOf(v))
	return
}

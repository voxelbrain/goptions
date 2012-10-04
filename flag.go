package goptions

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Flag represents a single flag of a FlagSet.
type Flag struct {
	Short        string
	Long         string
	MutexGroups  []string
	Description  string
	Obligatory   bool
	WasSpecified bool
	value        reflect.Value
}

var (
	typeBool   = reflect.TypeOf(bool(false))
	typeString = reflect.TypeOf(string("string"))
	typeInt    = reflect.TypeOf(int(0))
)

// Return the name of the flag preceding the right amount of dashes.
// The long name is preferred. If no name has been specified, "<unspecified>"
// will be returned.
func (f *Flag) Name() string {
	if len(f.Long) > 0 {
		return "--" + f.Long
	}
	if len(f.Short) > 0 {
		return "-" + f.Short
	}
	return "<unspecified>"
}

// NeedsExtraValue returns true if the flag expects a separate value.
func (f *Flag) NeedsExtraValue() bool {
	// Explicit over implicit
	if f.value.Type() == typeBool {
		return false
	}
	if _, ok := f.value.Interface().(Help); ok {
		return false
	}
	return true
}

// IsMulti returns true if the flag can be specified multiple times.
func (f *Flag) IsMulti() bool {
	if f.value.Kind() == reflect.Slice {
		return true
	}
	return false
}

func isShort(arg string) bool {
	return strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--")
}

func isLong(arg string) bool {
	return strings.HasPrefix(arg, "--")
}

func (f *Flag) Handles(arg string) bool {
	return (isShort(arg) && arg[1:2] == f.Short) ||
		(isLong(arg) && arg[2:] == f.Long)

}

func (f *Flag) Parse(args []string) ([]string, error) {
	param, value := args[0], ""
	if f.NeedsExtraValue() &&
		(len(args) < 2 || (isShort(param) && len(param) > 2)) {
		return args, fmt.Errorf("Flag %s needs an argument", f.Name())
	}
	if f.WasSpecified && !f.IsMulti() {
		return args, fmt.Errorf("Flag %s can only be specified once", f.Name())
	}
	if isShort(param) && len(param) > 2 {
		// Short flag cluster
		args[0] = "-" + param[2:]
	} else if f.NeedsExtraValue() {
		value = args[1]
		args = args[2:]
	} else {
		args = args[1:]
	}
	f.WasSpecified = true
	return args, f.setValue(value)
}

func (f *Flag) setValue(s string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	if m, ok := f.value.Interface().(Marshaler); ok {
		return m.MarshalGoption(s)
	} else if _, ok := f.value.Interface().(Help); ok {
		return ErrHelpRequest
	} else if f.value.Type() == typeBool {
		f.value.Set(reflect.ValueOf(true))
	} else if f.value.Type() == typeString {
		f.value.Set(reflect.ValueOf(s))
	} else if f.value.Type() == typeInt {
		intval, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		f.value.Set(reflect.ValueOf(int(intval)))
	} else {
		return fmt.Errorf("Unsupported flag type: %s", f.value.Type().Name())
	}
	return
}

// Old
/*
func (f *Flag) set() {
	f.WasSpecified = true
	if f.value.Kind() == reflect.Bool {
		f.setValue(true)
	} else if f.value.Kind() == reflect.Int && f.Accumulate {
		f.setValue(f.value.Interface().(int) + 1)
	}
}

func (f *Flag) setLong() {
	f.WasSpecifiedLong = true
	f.set()
}
func (f *Flag) setShort() {
	f.WasSpecifiedLong = false
	f.set()
}

func (f *Flag) setStringValue(val string) (err error) {
	switch f.value.Interface().(type) {
	case Marshaler:
		newval := reflect.New(f.value.Type()).Elem()
		if newval.Kind() == reflect.Ptr {
			newptrval := reflect.New(f.value.Type().Elem())
			newval.Set(newptrval)
		}
		err := newval.Interface().(Marshaler).MarshalGoption(val)
		f.value.Set(newval)
		return err
	case int:
		intval, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		return f.setValue(int(intval))
	case float64:
		intval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		return f.setValue(intval)
	default:
		return f.setValue(val)
	}
	return nil
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
*/

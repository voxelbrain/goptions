package goptions

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
	if f.value.Kind() == reflect.Bool {
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

func (f *Flag) Parse(args []string) ([]string, error) {
	if f.WasSpecified && !f.IsMulti() {
		return args, fmt.Errorf("Flag %s can only be specified once", f.Name())
	}
	if f.NeedsExtraValue() && len(args) <= 0 {
		return args, fmt.Errorf("Flag %s needs an argument", f.Name())
	}
	f.WasSpecified = true
	return f.parseArgument(args)
}

func (f *Flag) parseArgument(args []string) ([]string, error) {
	vkind := f.value.Kind()
	vtype := f.value.Type()

	// Loop
	panic("Invalid execution path")
}

func setValue(v reflect.Value, value interface{}) error {
}

// Old
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

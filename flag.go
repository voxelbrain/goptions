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
	DefaultValue reflect.Value
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
	if f.value.Type() == reflect.TypeOf(new(bool)).Elem() {
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

type valueParser func(val string) (reflect.Value, error)

var (
	parserMap = map[reflect.Type]valueParser{
		reflect.TypeOf(new(bool)).Elem():   boolValueParser,
		reflect.TypeOf(new(string)).Elem(): stringValueParser,
		reflect.TypeOf(new(int)).Elem():    intValueParser,
		reflect.TypeOf(new(Help)).Elem():   helpValueParser,
	}
)

func (f *Flag) setValue(s string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	if _, ok := f.value.Interface().(Marshaler); ok {
		newval := reflect.New(f.value.Type()).Elem()
		if newval.Kind() == reflect.Ptr {
			newptrval := reflect.New(f.value.Type().Elem())
			newval.Set(newptrval)
		}
		err := newval.Interface().(Marshaler).MarshalGoption(s)
		f.value.Set(newval)
		return err
	}
	vtype := f.value.Type()
	if f.value.Kind() == reflect.Slice {
		vtype = f.value.Type().Elem()
	}
	if parser, ok := parserMap[vtype]; ok {
		val, err := parser(s)
		if err != nil {
			return err
		}
		if f.value.Kind() == reflect.Slice {
			f.value.Set(reflect.Append(f.value, val))
		} else {
			f.value.Set(val)
		}
		return nil
	} else {
		return fmt.Errorf("Unsupported flag type: %s", f.value.Type().Name())
	}
	panic("Invalid execution path")
}

func boolValueParser(val string) (reflect.Value, error) {
	return reflect.ValueOf(true), nil
}

func stringValueParser(val string) (reflect.Value, error) {
	return reflect.ValueOf(val), nil
}

func intValueParser(val string) (reflect.Value, error) {
	intval, err := strconv.ParseInt(val, 10, 64)
	return reflect.ValueOf(int(intval)), err
}

func helpValueParser(val string) (reflect.Value, error) {
	return reflect.Value{}, ErrHelpRequest
}

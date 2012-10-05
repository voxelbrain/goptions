package goptions

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type valueParser func(f *Flag, val string) (reflect.Value, error)

var (
	parserMap = map[reflect.Type]valueParser{
		reflect.TypeOf(new(bool)).Elem():     boolValueParser,
		reflect.TypeOf(new(string)).Elem():   stringValueParser,
		reflect.TypeOf(new(int)).Elem():      intValueParser,
		reflect.TypeOf(new(Help)).Elem():     helpValueParser,
		reflect.TypeOf(new(*os.File)).Elem(): fileValueParser,
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
		val, err := parser(f, s)
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

func boolValueParser(f *Flag, val string) (reflect.Value, error) {
	return reflect.ValueOf(true), nil
}

func stringValueParser(f *Flag, val string) (reflect.Value, error) {
	return reflect.ValueOf(val), nil
}

func intValueParser(f *Flag, val string) (reflect.Value, error) {
	intval, err := strconv.ParseInt(val, 10, 64)
	return reflect.ValueOf(int(intval)), err
}

func fileValueParser(f *Flag, val string) (reflect.Value, error) {
	mode := 0
	if v, ok := f.optionMeta["file_mode"].(int); ok {
		mode = v
	}
	if val == "-" {
		if mode&os.O_RDONLY > 0 {
			return reflect.ValueOf(os.Stdin), nil
		} else if mode&os.O_WRONLY > 0 {
			return reflect.ValueOf(os.Stdout), nil
		}
	} else {
		perm := uint32(0644)
		if v, ok := f.optionMeta["file_perm"].(uint32); ok {
			perm = v
		}
		f, e := os.OpenFile(val, mode, os.FileMode(perm))
		return reflect.ValueOf(f), e
	}
	panic("Invalid execution path")
}

func helpValueParser(f *Flag, val string) (reflect.Value, error) {
	return reflect.Value{}, ErrHelpRequest
}

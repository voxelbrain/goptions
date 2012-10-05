package goptions

import (
	"fmt"
	"reflect"
	"strconv"
)

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

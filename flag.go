package goptions

import (
	"errors"
	"reflect"
)

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

func (f *Flag) Name() string {
	if len(f.Long) > 0 {
		return "--" + f.Long[0]
	}
	if len(f.Short) > 0 {
		return "-" + f.Short[0]
	}
	return "<unspecified>"
}

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

func (f *Flag) Set() {
	f.WasSpecified = true
	if f.value.Kind() == reflect.Bool {
		f.SetValue(true)
	} else if f.value.Kind() == reflect.Int && f.Accumulate {
		f.SetValue(f.value.Interface().(int) + 1)
	}
}

func (f *Flag) SetValue(v interface{}) (err error) {
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

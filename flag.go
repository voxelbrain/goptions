package goptions

import (
	"errors"
	"reflect"
)

type flag struct {
	Value        reflect.Value
	Short        []string
	Long         []string
	MutexGroup   string
	Accumulate   bool
	Description  string
	Obligatory   bool
	WasSpecified bool
}

func (f *flag) Name() string {
	if len(f.Long) > 0 {
		return "--" + f.Long[0]
	}
	if len(f.Short) > 0 {
		return "-" + f.Short[0]
	}
	return "<unspecified>"
}

func (f *flag) NeedsExtraValue() bool {
	// Explicit over implicit
	if f.Value.Kind() == reflect.Bool {
		return false
	}
	if f.Value.Kind() == reflect.Int && f.Accumulate {
		return false
	}
	return true
}

func (f *flag) Set() {
	f.WasSpecified = true
	if f.Value.Kind() == reflect.Bool {
		f.SetValue(true)
	} else if f.Value.Kind() == reflect.Int && f.Accumulate {
		f.SetValue(f.Value.Interface().(int) + 1)
	}
}

func (f *flag) SetValue(v interface{}) (err error) {
	defer func() {
		if x := recover(); x != nil {
			if str, ok := x.(string); ok {
				err = errors.New(str)
				return
			}
			err = x.(error)
		}
	}()
	f.Value.Set(reflect.ValueOf(v))
	return
}

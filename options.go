package goptions

import (
	"fmt"
	"reflect"
)

type optionFunc func(f *Flag, option, value string) error
type optionMap map[string]optionFunc

var (
	typeOptionMap = map[reflect.Type]optionMap{
		// Global options
		nil: optionMap{
			"description": description,
			"obligatory":  obligatory,
			"mutexgroup":  mutexgroup,
		},
		typeInt: optionMap{
			"accumulate": accumulate,
		},
	}
)

func description(f *Flag, option, value string) error {
	f.Description = value
	return nil
}

func obligatory(f *Flag, option, value string) error {
	f.Obligatory = true
	return nil
}

func mutexgroup(f *Flag, option, value string) error {
	if len(value) <= 0 {
		return fmt.Errorf("Mutexgroup option needs a value")
	}
	f.MutexGroups = append(f.MutexGroups, value)
	return nil
}

func accumulate(f *Flag, option, value string) error {
	return nil
}

func optionMapForType(t reflect.Type) optionMap {
	g := typeOptionMap[nil]
	m, _ := typeOptionMap[t]
	r := make(optionMap)
	for k, v := range g {
		r[k] = v
	}
	for k, v := range m {
		r[k] = v
	}
	return r
}

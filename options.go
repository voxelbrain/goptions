package goptions

import (
	"fmt"
	"reflect"
	"strings"
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
	}
)

func description(f *Flag, option, value string) error {
	f.Description = strings.Replace(value, `\`, ``, -1)
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
	for _, group := range strings.Split(value, ",") {
		f.MutexGroups = append(f.MutexGroups, group)
	}
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

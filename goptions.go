package goptions

import (
	"fmt"
	"reflect"
)

type Verbs map[string]interface{}

func Parse(args []string, v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		panic("Value type is not a pointer to a struct")
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		panic("Value type is not a pointer to a struct")
	}

	flags := parseFields(value)
	_ = flags
	return nil
}

type Flag struct {
	Value       reflect.Value
	Short       []string
	Long        []string
	MutexGroup  string
	Accumulate  bool
	Description string
	NonZero     bool
}

func parseFields(value reflect.Value) []*Flag {
	r := make([]*Flag, 0)
	for i := 0; i < value.Type().NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get("goptions")
		value := value.Field(i)
		flag, e := parseTag(tag)
		if e != nil {
			panic(fmt.Sprintf("Invalid tagline: %s", e))
		}
		flag.Value = value
		r = append(r, flag)

	}
	return r
}

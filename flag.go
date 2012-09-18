package goptions

import (
	"fmt"
	"reflect"
	"strings"
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

func (f *flag) NeedsValue() bool {
	// Explicit over implicit
	if f.Value.Kind() == reflect.Bool {
		return false
	}
	if f.Value.Kind() == reflect.Int && f.Accumulate {
		return false
	}
	return true
}

type flagSet []*flag

func newFlagSet(structValue reflect.Value) (flagSet, error) {
	r := make(flagSet, 0)
	for i := 0; i < structValue.Type().NumField(); i++ {
		tag := structValue.Type().Field(i).Tag.Get("goptions")
		fieldValue := structValue.Field(i)
		flag, err := parseTag(tag)
		if err != nil {
			return nil, fmt.Errorf("Invalid tagline: %s", err)
		}
		flag.Value = fieldValue
		r = append(r, flag)

	}
	return r, nil
}

func (fs flagSet) ShortFlagMap() map[string]*flag {
	r := make(map[string]*flag)
	for _, flag := range fs {
		for _, short := range flag.Short {
			r[short] = flag
		}
	}
	return r
}

func (fs flagSet) LongFlagMap() map[string]*flag {
	r := make(map[string]*flag)
	for _, flag := range fs {
		for _, long := range flag.Long {
			r[long] = flag
		}
	}
	return r
}

func (fs flagSet) MutexGroups() map[string][]*flag {
	r := make(map[string][]*flag)
	for _, f := range fs {
		mg := f.MutexGroup
		if _, ok := r[mg]; !ok {
			r[mg] = make([]*flag, 0)
		}
		r[mg] = append(r[mg], f)
	}
	return r
}

func (fs flagSet) Parse(args []string) error {
	shortMap, longMap := fs.ShortFlagMap(), fs.LongFlagMap()

	for i := 0; i < len(args); i++ {
		arg := args[i]
		var f *flag
		if strings.HasPrefix(arg, "--") {
			longname := arg[2:]
			if _, ok := longMap[longname]; !ok {
				return fmt.Errorf("Unknown flag %s", arg)
			}
			f = longMap[longname]
		} else if strings.HasPrefix(arg, "-") {
			shortname := arg[1:]
			if _, ok := shortMap[shortname]; !ok {
				return fmt.Errorf("Unknown flag %s", arg)
			}
			f = shortMap[shortname]
		}
		switch f.Value.Kind() {
		case reflect.String:
			f.Value.SetString(args[i+1])
			i++
		default:
			return fmt.Errorf("Unsupported type %s", f.Value.Kind().String())
		}
		f.WasSpecified = true
	}
	return nil
}

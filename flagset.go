package goptions

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
)

// HelpFunc is the signature of a function responsible for printing the help.
type HelpFunc func(w io.Writer, fs *FlagSet)

// A FlagSet represents one set of flags which belong to one particular program.
// A FlagSet is also used to represent a subset of flags belonging to one verb.
type FlagSet struct {
	// This HelpFunc will be called when PrintHelp() is called.
	HelpFunc
	// Name of the program. Might be used by HelpFunc.
	Name          string
	helpFlag      *Flag
	remainderFlag *Flag
	shortMap      map[string]*Flag
	longMap       map[string]*Flag
	// Global option flags
	Flags []*Flag
	// Verbs and corresponding FlagSets
	Verbs map[string]*FlagSet
}

// NewFlagSet returns a new FlagSet containing all the flags which result from
// parsing the tags of the struct. Said struct as to be passed to the function
// as a pointer.
// If a tag line is erroneous, NewFlagSet() panics as this is considered a
// compile time error rather than a runtme error.
func NewFlagSet(name string, v interface{}) *FlagSet {
	structValue := reflect.ValueOf(v)
	if structValue.Kind() != reflect.Ptr {
		panic("Value type is not a pointer to a struct")
	}
	structValue = structValue.Elem()
	if structValue.Kind() != reflect.Struct {
		panic("Value type is not a pointer to a struct")
	}
	return newFlagset(name, structValue, nil)
}

// Internal version which skips type checking.
// Can't obtain a pointer to a struct field using reflect.
func newFlagset(name string, structValue reflect.Value, remainder *Flag) *FlagSet {
	var once sync.Once
	r := &FlagSet{
		Name:          name,
		Flags:         make([]*Flag, 0),
		HelpFunc:      DefaultHelpFunc,
		remainderFlag: remainder,
	}

	var i int
	// Parse Option fields
	for i = 0; i < structValue.Type().NumField(); i++ {
		fieldValue := structValue.Field(i)
		if fieldValue.Type().Name() == "Verbs" {
			break
		}
		tag := structValue.Type().Field(i).Tag.Get("goptions")
		flag, err := parseTag(tag)
		if err != nil {
			panic(fmt.Sprintf("Invalid tagline: %s", err))
		}
		flag.value = fieldValue

		switch flag.value.Interface().(type) {
		case Help:
			r.helpFlag = flag
		case Remainder:
			if r.remainderFlag == nil {
				r.remainderFlag = flag
			}
		}

		if len(tag) != 0 {
			r.Flags = append(r.Flags, flag)
		}
	}
	r.shortMap, r.longMap = r.shortFlagMap(), r.longFlagMap()

	// Parse verb fields
	for i++; i < structValue.Type().NumField(); i++ {
		once.Do(func() {
			r.Verbs = make(map[string]*FlagSet)
		})
		fieldValue := structValue.Field(i)
		tag := structValue.Type().Field(i).Tag.Get("goptions")
		r.Verbs[tag] = newFlagset(tag, fieldValue, r.remainderFlag)
	}

	return r
}

func (fs *FlagSet) shortFlagMap() map[string]*Flag {
	r := make(map[string]*Flag)
	for _, flag := range fs.Flags {
		for _, short := range flag.Short {
			r[short] = flag
		}
	}
	return r
}

func (fs *FlagSet) longFlagMap() map[string]*Flag {
	r := make(map[string]*Flag)
	for _, flag := range fs.Flags {
		for _, long := range flag.Long {
			r[long] = flag
		}
	}
	return r
}

// MutexGroups returns a map of Flag lists which contain mutually
// exclusive flags.
func (fs *FlagSet) MutexGroups() map[string][]*Flag {
	r := make(map[string][]*Flag)
	for _, f := range fs.Flags {
		mg := f.MutexGroup
		if len(mg) == 0 {
			continue
		}
		if _, ok := r[mg]; !ok {
			r[mg] = make([]*Flag, 0)
		}
		r[mg] = append(r[mg], f)
	}
	return r
}

var (
	ErrHelpRequest = errors.New("Request for Help")
)

// Parse takes the command line arguments and sets the corresponding values
// in the FlagSet's struct.
func (fs *FlagSet) Parse(args []string) error {
	for len(args) > 0 {
		restArgs, err := fs.parseNextItem(args)
		if err != nil {
			return err
		}
		args = restArgs
	}

	if fs.helpFlag != nil && fs.helpFlag.WasSpecified {
		return ErrHelpRequest
	}

	// Check for unset, obligatory Flags
	for _, f := range fs.Flags {
		if f.Obligatory && !f.WasSpecified {
			return fmt.Errorf("%s must be specified", f.Name())
		}
	}

	// Check for multiple set Flags in one mutex group
	mgs := fs.MutexGroups()
	for _, mg := range mgs {
		wasSpecifiedCount := 0
		names := make([]string, 0)
		for _, flag := range mg {
			names = append(names, flag.Name())
			if flag.WasSpecified {
				wasSpecifiedCount += 1
			}
		}
		if wasSpecifiedCount >= 2 {
			return fmt.Errorf("Only one of %s can be specified", strings.Join(names, ", "))
		}
	}
	return nil
}

func (fs *FlagSet) parseNextItem(args []string) ([]string, error) {
	if strings.HasPrefix(args[0], "--") {
		return fs.parseLongFlag(args)
	} else if strings.HasPrefix(args[0], "-") {
		return fs.parseShortFlagCluster(args)
	}

	if fs.Verbs != nil {
		verb, ok := fs.Verbs[args[0]]
		if !ok {
			return args, fmt.Errorf("Unknown verb: %s", args[0])
		}
		err := verb.Parse(args[1:])
		if err != nil {
			return args, err
		}
		return []string{}, nil
	}
	if fs.remainderFlag != nil {
		remainder := reflect.MakeSlice(fs.remainderFlag.value.Type(), len(args), len(args))
		reflect.Copy(remainder, reflect.ValueOf(args))
		fs.remainderFlag.value.Set(remainder)
		return []string{}, nil
	}
	return args, fmt.Errorf("Invalid trailing arguments: %v", args)
}

func (fs *FlagSet) parseLongFlag(args []string) ([]string, error) {
	longflagname := args[0][2:]
	f, ok := fs.longMap[longflagname]
	if !ok {
		return args, fmt.Errorf("Unknown flag --%s", longflagname)
	}
	args = args[1:]
	f.set()
	if f.NeedsExtraValue() {
		err := f.setStringValue(args[0])
		if err != nil {
			return args, err
		}
		args = args[1:]
	}
	return args, nil
}

func (fs *FlagSet) parseShortFlagCluster(args []string) ([]string, error) {
	shortflagnames := args[0][1:]
	args = args[1:]
	for idx, shortflagname := range shortflagnames {
		flag, ok := fs.shortMap[string(shortflagname)]
		if !ok {
			return args, fmt.Errorf("Unknown flag -%s", string(shortflagname))
		}
		flag.set()
		// If value-flag is given but is not the last in a short flag cluster,
		// it's an error.
		if flag.NeedsExtraValue() && idx != len(shortflagnames)-1 {
			return args, fmt.Errorf("Flag %s needs a value", flag.Name())
		} else if flag.NeedsExtraValue() {
			err := flag.setStringValue(args[0])
			if err != nil {
				return args, err
			}
			args = args[1:]
		}
	}
	return args, nil
}

// Prints the FlagSet's help to the given writer.
func (fs *FlagSet) PrintHelp(w io.Writer) {
	fs.HelpFunc(w, fs)
}

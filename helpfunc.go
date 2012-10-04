package goptions

import (
	"io"
	"reflect"
	"sync"
	"text/tabwriter"
	"text/template"
	// "log"
)

// HelpFunc is the signature of a function responsible for printing the help.
type HelpFunc func(w io.Writer, fs *FlagSet)

var defaultFuncMap = map[string]interface{}{
	"notZero": func(arg interface{}) bool {
		switch x := arg.(type) {
		case reflect.Value:
			return !reflect.DeepEqual(reflect.Zero(x.Type()).Interface(), x.Interface())
		default:
			return !reflect.DeepEqual(reflect.Zero(reflect.TypeOf(arg)), reflect.ValueOf(arg))
		}
		panic("Invalid execution path")
	},
	"dereflect": func(arg reflect.Value) interface{} {
		return arg.Interface()
	},
}

// Generates a new HelpFunc taking a `text/template.Template`-formatted
// string as an argument. The resulting template will be executed with the FlagSet
// as its data.
func NewTemplatedHelpFunc(tpl string) HelpFunc {
	var once sync.Once
	var t *template.Template
	return func(w io.Writer, fs *FlagSet) {
		once.Do(func() {
			t = template.Must(template.New("helpTemplate").Funcs(defaultFuncMap).Parse(tpl))
		})
		err := t.Execute(w, fs)
		if err != nil {
			panic(err)
		}
	}
}

const (
	_DEFAULT_HELP = `Usage: {{.Name}} [global options] {{with .Verbs}}<verb> [verb options]{{end}}

Global options:{{range .Flags}}
	{{with .Short}}-{{.}},{{end}}	{{with .Long}}--{{.}}{{end}}	{{.Description}}{{if notZero .DefaultValue}} (default: {{dereflect .DefaultValue}}){{end}}{{if .Obligatory}} (*){{end}}{{end}}

{{with .Verbs}}Verbs:{{range .}}
	{{.Name}}:{{range .Flags}}
		{{with .Short}}-{{.}},{{end}}	{{with .Long}}--{{.}}{{end}}	{{.Description}}{{if notZero .DefaultValue}} (default: {{dereflect .DefaultValue}}){{end}}{{if .Obligatory}} (*){{end}}{{end}}{{end}}{{end}}

`
)

// DefaultHelpFunc is a HelpFunc which renders the default help template and pipes
// the output through a text/tabwriter.Writer before flushing it to the output.
func DefaultHelpFunc(w io.Writer, fs *FlagSet) {
	tw := &tabwriter.Writer{}
	tw.Init(w, 4, 4, 1, ' ', 0)
	NewTemplatedHelpFunc(_DEFAULT_HELP)(tw, fs)
	tw.Flush()
}

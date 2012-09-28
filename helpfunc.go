package goptions

import (
	"io"
	"sync"
	"text/tabwriter"
	"text/template"
)

// HelpFunc is the signature of a function responsible for printing the help.
type HelpFunc func(w io.Writer, fs *FlagSet)

// Generates a new HelpFunc taking a `text/template.Template`-formatted
// string as an argument. The resulting template will be executed with the FlagSet
// as its data.
func NewTemplatedHelpFunc(tpl string) HelpFunc {
	var once sync.Once
	var t *template.Template
	return func(w io.Writer, fs *FlagSet) {
		once.Do(func() {
			t = template.Must(template.New("helpTemplate").Parse(tpl))
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
	{{if len .Short}}-{{index .Short 0}},{{end}}	{{if len .Long}}--{{index .Long 0}}{{end}}	{{.Description}}{{if .Obligatory}} (*){{end}}{{end}}

{{if .Verbs}}Verbs:{{range .Verbs}}
	{{.Name}}:{{range .Flags}}
		{{if len .Short}}-{{index .Short 0}},{{end}}	{{if len .Long}}--{{index .Long 0}}{{end}}	{{.Description}}{{if .Obligatory}} (*){{end}}{{end}}{{end}}{{end}}

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

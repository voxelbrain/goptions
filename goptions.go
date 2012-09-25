package goptions

import (
	"io"
	"os"
	"sync"
	"text/tabwriter"
	"text/template"
)

var (
	globalFlagSet *FlagSet
)

func Parse(v interface{}) error {
	fs, err := NewFlagSet(os.Args[0], v)
	if err != nil {
		return err
	}
	globalFlagSet = fs

	e := fs.Parse(os.Args[1:])
	if e != nil {
		return e
	}

	return nil
}

func PrintHelp() {
	if globalFlagSet == nil {
		panic("Must call Parse() before PrintHelp()")
	}
	globalFlagSet.PrintHelp(os.Stdout)
}

// Generates a new HelpFunc taking a `text/template.Template`-formatted
// string as an argument.
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

func Must(fs *FlagSet, err error) *FlagSet {
	if err != nil {
		panic(err)
	}
	return fs
}

const (
	DEFAULT_HELP = `
Usage: {{.Name}} [global options] {{with .Verbs}}<verb> [verb options]{{end}}

Global options:{{range .Flags}}
	{{if len .Short}}-{{index .Short 0}},{{end}}	{{if len .Long}}--{{index .Long 0}}{{end}}	{{.Description}}{{if .Obligatory}} (*){{end}}{{end}}

{{if .Verbs}}Verbs:{{range .Verbs}}
	{{.Name}}:{{range .Flags}}
		{{if len .Short}}-{{index .Short 0}},{{end}}	{{if len .Long}}--{{index .Long 0}}{{end}}	{{.Description}}{{if .Obligatory}} (*){{end}}{{end}}{{end}}{{end}}
`
)

var DefaultHelpFunc = func() HelpFunc {
	tplhf := NewTemplatedHelpFunc(DEFAULT_HELP)
	return func(w io.Writer, fs *FlagSet) {
		tw := &tabwriter.Writer{}
		tw.Init(w, 4, 4, 1, ' ', 0)
		tplhf(tw, fs)
		tw.Flush()
	}
}()

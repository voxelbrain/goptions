`goptions` implements a flexible parser for command line options.

Key targets were the support for both long and short flag versions, mutually
exclusive flags, and verbs. Flags and their corresponding variables are defined
by the tags in a (possibly anonymous) struct.

# Example

```Go
package main

import (
	"github.com/voxelbrain/goptions"
	"os"
)

func main() {
	options := struct {
		Server        *net.TCPAddr `goptions:"-s, --server, obligatory, description='Server to connect to'"`
		Password      string       `goptions:"-p, --password, description='Don\\'t prompt for password'"`
		Timeout       int          `goptions:"-t, --timeout, description='Connection timeout in seconds'"`
		goptions.Help              `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
		Execute struct {
			Command string   `goptions:"--command, mutexgroup='input', description='Command to exectute', obligatory"`
			Script  *os.File `goptions:"--script, mutexgroup='input', description='Script to execture', create, wronly, append"`
		} `goptions:"execute"`
		Delete struct {
			Path  string `goptions:"-n, --name, obligatory, description='Name of the entity to be deleted'"`
			Force bool   `goptions:"-f, --force, description='Force removal'"`
		} `goptions:"delete"`
	}{ // Default values goes here
		Timeout: 10,
	}
	goptions.ParseAndFail(&options)
}
```

```
$ go run examples/readme_example.go --help
Usage: a.out [global options] <verb> [verb options]

Global options:
    -s, --server   Server to connect to (*)
    -p, --password Don't prompt for password
    -t, --timeout  Connection timeout in seconds (default: 10)
    -h, --help     Show this help

Verbs:
    delete:
        -n, --name  Name of the entity to be deleted (*)
        -f, --force Force removal
    execute:
            --command Command to exectute (*)
            --script  Script to execture
```

---
Version 2.4.0

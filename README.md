`goptions` implements a flexible parser for command line options.

Key targets were the support for both long and short flag versions, mutually
exclusive flags, and verbs. Flags and their corresponding variables are defined
by the tags in a (possibly anonymous) struct.

![](https://circleci.com/gh/voxelbrain/goptions.png?circle-token=27cd98362d475cfa8c586565b659b2204733f25c)

# Example

```Go
package main

import (
	"github.com/voxelbrain/goptions"
	"os"
	"time"
)

func main() {
	options := struct {
		Servers  []string      `goptions:"-s, --server, obligatory, description='Servers to connect to'"`
		Password string        `goptions:"-p, --password, description='Don\\'t prompt for password'"`
		Timeout  time.Duration `goptions:"-t, --timeout, description='Connection timeout in seconds'"`
		Help     goptions.Help `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
		Execute struct {
			Command string   `goptions:"--command, mutexgroup='input', description='Command to exectute', obligatory"`
			Script  *os.File `goptions:"--script, mutexgroup='input', description='Script to exectute', rdonly"`
		} `goptions:"execute"`
		Delete struct {
			Path  string `goptions:"-n, --name, obligatory, description='Name of the entity to be deleted'"`
			Force bool   `goptions:"-f, --force, description='Force removal'"`
		} `goptions:"delete"`
	}{ // Default values goes here
		Timeout: 10 * time.Second,
	}
	goptions.ParseAndFail(&options)
}
```

```
$ go run examples/readme_example.go --help
Usage: a.out [global options] <verb> [verb options]

Global options:
        -s, --server   Servers to connect to (*)
        -p, --password Don't prompt for password
        -t, --timeout  Connection timeout in seconds (default: 10s)
        -h, --help     Show this help

Verbs:
    delete:
        -n, --name     Name of the entity to be deleted (*)
        -f, --force    Force removal
    execute:
            --command  Command to exectute (*)
            --script   Script to exectute
```

---
Version 2.5.8

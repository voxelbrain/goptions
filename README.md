`goptions` implements a flexible parser for command line options.

Key targets were the support for both long and short flag version, mutually
exclusive flags and verbs. Flags are defined by a (possibly anonymous) struct
and the corresponding tags.

# Example

```Go
func ExampleFlagSet_PrintHelp() {
	var options struct {
		Server    string `goptions:"-s, --server, obligatory, description='Server to connect to'"`
		Password  string `goptions:"-p, --password, description='Don\\'t prompt for password'"`
		Verbosity int    `goptions:"-v, --verbose, accumulate, description='Set output threshold level'"`
		Help      `goptions:"-h, --help, description='Show this help'"`

		Verbs
		Create struct {
			Name      string `goptions:"-n, --name, obligatory, description='Name of the entity to be created'"`
			Directory bool   `goptions:"--directory, mutexgroup='type', description='Create a directory'"`
			File      bool   `goptions:"--file, mutexgroup='type', description='Create a file'"`
		} `goptions:"create"`
		Delete struct {
			Name      string `goptions:"-n, --name, obligatory, description='Name of the entity to be deleted'"`
			Directory bool   `goptions:"--directory, mutexgroup='type', description='Delete a directory'"`
			File      bool   `goptions:"--file, mutexgroup='type', description='Delete a file'"`
		} `goptions:"delete"`
	}
	args := []string{"--help"}
	fs := Must(NewFlagSet("goptions", &options))
	err := fs.Parse(args)
	if err == ErrHelpRequest {
		fs.PrintHelp(os.Stderr)
		return
	} else if err != nil {
		panic(err)
	}

	// Output:
	// Usage: goptions [global options] <verb> [verb options]
	//
	// Global options:
	//     -s, --server   Server to connect to (*)
	//     -p, --password Don't prompt for password
	//     -v, --verbose  Set output threshold level
	//     -h, --help     Show this help
	//
	// Verbs:
	//     create:
	//         -n, --name      Name of the entity to be created (*)
	//             --directory Create a directory
	//             --file      Create a file
	//     delete:
	//         -n, --name      Name of the entity to be deleted (*)
	//             --directory Delete a directory
	//             --file      Delete a file
}
```

---
Version 1.0

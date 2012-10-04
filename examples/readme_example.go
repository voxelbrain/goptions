package main

import (
	"github.com/voxelbrain/goptions"
)

func main() {
	options := struct {
		Server    string `goptions:"-s, --server, obligatory, description='Server to connect to'"`
		Password  string `goptions:"-p, --password, description='Don\\'t prompt for password'"`
		Verbosity int    `goptions:"-v, --verbose, description='Be verbose'"`
		Timeout   int    `goptions:"-t, --timeout, description='Connection timeout in seconds'"`
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
	}{ // Default values goes here
		Timeout: 10,
	}
	goptions.ParseAndFail(&options)
}

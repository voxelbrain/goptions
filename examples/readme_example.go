package main

import (
	"github.com/voxelbrain/goptions"
	"os"
)

func main() {
	options := struct {
		Server        string `goptions:"-s, --server, obligatory, description='Server to connect to'"`
		Password      string `goptions:"-p, --password, description='Don\\'t prompt for password'"`
		Timeout       int    `goptions:"-t, --timeout, description='Connection timeout in seconds'"`
		goptions.Help `goptions:"-h, --help, description='Show this help'"`

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

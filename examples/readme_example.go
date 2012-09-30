package main

import (
	"fmt"
	"github.com/voxelbrain/goptions"
)

func main() {
	options := struct {
		Server        string `goptions:"-s, --server, description='Server to connect to'"`
		Password      string `goptions:"-p, --password, description='Don\\'t prompt for password'"`
		Verbosity     int    `goptions:"-v, --verbose, accumulate, description='Set output threshold level'"`
		goptions.Help `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
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
	}{ // Default values go here
		Server: "localhost",
	}
	err := goptions.Parse(&options)
	if err != nil {
		if err != goptions.ErrHelpRequest {
			fmt.Printf("Error: %s\n", err)
		}
		goptions.PrintHelp()
		return
	}
}

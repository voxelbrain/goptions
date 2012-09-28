package goptions

import (
	"fmt"
	"strings"
	"testing"
)

type Name struct {
	FirstName string
	LastName  string
}

func (n *Name) MarshalGoption(val string) error {
	f := strings.SplitN(val, " ", 2)
	if len(f) != 2 {
		return fmt.Errorf("Incomplete name")
	}
	n.FirstName = f[0]
	n.LastName = f[1]
	return nil
}

func TestMarshaler(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Name *Name `goptions:"--name"`
	}
	args = []string{"--name", "Alexander Surma"}
	fs = NewFlagSet("goptions", &options)
	err = fs.Parse(args)
	if err != nil ||
		options.Name.FirstName != "Alexander" ||
		options.Name.LastName != "Surma" {
		t.Fatalf("Unexpected value: %#v", options)
	}
}

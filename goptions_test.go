package goptions

import (
	"testing"
)

func TestParse_StringValue(t *testing.T) {
	var args []string
	var err error
	var options struct {
		Name string `goptions:"--name, -n"`
	}
	expected := "SomeName"

	args = []string{"--name", "SomeName"}
	err = Parse(args, &options)
	if err != nil {
		t.Fatalf("flag parsing failed: %s", err)
	}
	if options.Name != expected {
		t.Fatalf("Expected %s for options.Name, got %s", expected, options.Name)
	}

	options.Name = ""

	args = []string{"-n", "SomeName"}
	err = Parse(args, &options)
	if err != nil {
		t.Fatalf("flag parsing failed: %s", err)
	}
	if options.Name != expected {
		t.Fatalf("Expected %s for options.Name, got %s", expected, options.Name)
	}
}

func TestParse_ObligatoryStringValue(t *testing.T) {
	var args []string
	var err error
	var options struct {
		Name string `goptions:"-n, obligatory"`
	}
	args = []string{}
	err = Parse(args, &options)
	if err == nil {
		t.Fatalf("parsing should have failed.")
	}

	args = []string{"-n", "SomeName"}
	err = Parse(args, &options)
	if err != nil {
		t.Fatalf("parsing failed: %s", err)
	}

	expected := "SomeName"
	if options.Name != expected {
		t.Fatalf("Expected %s for options.Name, got %s", expected, options.Name)
	}
}

func TestParse_UnknownFlag(t *testing.T) {
	var args []string
	var err error
	var options struct {
		Name string `goptions:"--name, -n"`
	}
	args = []string{"-k", "4"}
	err = Parse(args, &options)
	if err == nil {
		t.Fatalf("Parsing should have failed.")
	}
}

func TestParse_FlagCluster(t *testing.T) {
	var args []string
	var err error
	var options struct {
		Fast    bool `goptions:"-f"`
		Silent  bool `goptions:"-q"`
		Serious bool `goptions:"-s"`
		Crazy   bool `goptions:"-c"`
		Verbose int  `goptions:"-v, accumulate"`
	}
	args = []string{"-fqcvvv"}
	err = Parse(args, &options)
	if err != nil {
		t.Fatalf("parsing failed: %s", err)
	}

	if !(options.Fast &&
		options.Silent &&
		!options.Serious &&
		options.Crazy &&
		options.Verbose == 3) {
		t.Fatalf("Unexpected value: %v", options)
	}
}

func TestParseTag_minimal(t *testing.T) {
	var tag string
	tag = `--name, -n, description='Some name'`
	f, e := parseTag(tag)
	if e != nil {
		t.Fatalf("Tag parsing failed: %s", e)
	}
	expected := &flag{
		Long:        []string{"name"},
		Short:       []string{"n"},
		Description: "Some name",
	}
	if !flagEqual(f, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, f)
	}
}

func TestParseTag_more(t *testing.T) {
	var tag string
	tag = `--name, -n, description='Some name', mutexgroup='selector', obligatory`
	f, e := parseTag(tag)
	if e != nil {
		t.Fatalf("Tag parsing failed: %s", e)
	}
	expected := &flag{
		Long:        []string{"name"},
		Short:       []string{"n"},
		Accumulate:  false,
		Description: "Some name",
		MutexGroup:  "selector",
		Obligatory:  true,
	}
	if !flagEqual(f, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, f)
	}
}

func flagEqual(f1, f2 *flag) bool {
	if !stringArrayEqual(f1.Long, f2.Long) {
		return false
	}
	if !stringArrayEqual(f1.Short, f2.Short) {
		return false
	}
	if f1.MutexGroup != f2.MutexGroup {
		return false
	}
	if f1.Accumulate != f2.Accumulate {
		return false
	}
	if f1.Description != f2.Description {
		return false
	}
	if f1.Obligatory != f2.Obligatory {
		return false
	}
	return true
}

func stringArrayEqual(a1, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

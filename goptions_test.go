package goptions

import (
	"testing"
)

func TestSimpleStruct(t *testing.T) {
	var args []string
	var err error
	var options struct {
		Force     bool   `goptions:"--force, -f, description='Force action'"`
		Verbosity int    `goptions:"--verbose, -v, description='Level of verbosity, accumulate'"`
		Name      string `goptions:"--name, -n, description='Some name', non-zero"`
	}
	args = []string{}
	err = Parse(args, &options)
	_ = err
}

func TestParseTag1(t *testing.T) {
	var tag string
	tag = `--name, -n, description='Some name', mutexgroup='selector', non-zero`
	f, e := parseTag(tag)
	if e != nil {
		t.Fatalf("Tag parsing failed: %s", e)
	}
	expected := &Flag{
		Long:        []string{"name"},
		Short:       []string{"n"},
		Accumulate:  false,
		Description: "Some name",
		MutexGroup:  "selector",
		NonZero:     true,
	}
	if !flagEqual(f, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, f)
	}
}

func TestParseTag2(t *testing.T) {
	var tag string
	tag = `--verbose, -v, description='Increase verbosity', accumulate`
	f, e := parseTag(tag)
	if e != nil {
		t.Fatalf("Tag parsing failed: %s", e)
	}
	expected := &Flag{
		Long:        []string{"verbose"},
		Short:       []string{"v"},
		Accumulate:  true,
		Description: "Increase verbosity",
	}
	if !flagEqual(f, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, f)
	}
}

func flagEqual(f1, f2 *Flag) bool {
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

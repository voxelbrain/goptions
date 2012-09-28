package goptions

import (
	"reflect"
	"testing"
)

func TestParseTag_Minimal(t *testing.T) {
	var tag string
	tag = `--name, -n, description='Some name'`
	f, e := parseTag(tag)
	if e != nil {
		t.Fatalf("Tag parsing failed: %s", e)
	}
	expected := &Flag{
		Long:        "name",
		Short:       "n",
		Description: "Some name",
	}
	if !reflect.DeepEqual(f, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, f)
	}
}

func TestParseTag_More(t *testing.T) {
	var tag string
	tag = `--name, -n, description='Some name', mutexgroup='selector', obligatory`
	f, e := parseTag(tag)
	if e != nil {
		t.Fatalf("Tag parsing failed: %s", e)
	}
	expected := &Flag{
		Long:        "name",
		Short:       "n",
		Accumulate:  false,
		Description: "Some name",
		MutexGroup:  "selector",
		Obligatory:  true,
	}
	if !reflect.DeepEqual(f, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, f)
	}
}

func TestParseTag_MultipleFlags(t *testing.T) {
	var tag string
	var e error
	tag = `--name1, --name2`
	_, e = parseTag(tag)
	if e == nil {
		t.Fatalf("Parsing should have failed")
	}

	tag = `-n, -v`
	_, e = parseTag(tag)
	if e == nil {
		t.Fatalf("Parsing should have failed")
	}
}

package goptions

import (
	"reflect"
	"testing"
)

func TestParseTag_minimal(t *testing.T) {
	var tag string
	tag = `--name, -n, description='Some name'`
	f, e := parseTag(tag)
	if e != nil {
		t.Fatalf("Tag parsing failed: %s", e)
	}
	expected := &Flag{
		Long:        []string{"name"},
		Short:       []string{"n"},
		Description: "Some name",
	}
	if !reflect.DeepEqual(f, expected) {
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
	expected := &Flag{
		Long:        []string{"name"},
		Short:       []string{"n"},
		Accumulate:  false,
		Description: "Some name",
		MutexGroup:  "selector",
		Obligatory:  true,
	}
	if !reflect.DeepEqual(f, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, f)
	}
}

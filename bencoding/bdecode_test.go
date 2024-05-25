package bencoding

import (
	"reflect"
	"testing"
)

//TODO: implement more tests (nested lists, dicts, etc.)

func TestInteger(t *testing.T) {
	source := "i34e"
	expected := 34
	value, _ := Decode(source)
	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}
}

func TestString(t *testing.T) {
	source := "11:Hello World"
	expected := "Hello World"
	value, _ := Decode(source)

	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}
}

func TestList(t *testing.T) {
	source := "l4:spam4:eggse"
	expected := []any{"spam", "eggs"}
	value, err := Decode(source)
	if !reflect.DeepEqual(expected, value) {
		t.Errorf("Expected %v got %v, Err: %v", expected, value, err)
	}
}

func TestDict(t *testing.T) {
	source := "d3:cow3:moo4:spam4:eggse"
	expected := map[string]any{"cow": "moo", "spam": "eggs"}
	value, err := Decode(source)

	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected %v got %v, Err: %v", expected, value, err)
	}
}

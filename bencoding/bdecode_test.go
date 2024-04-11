package bencoding

import (
	"reflect"
	"testing"
)

//TODO: implement more tests (nested lists, dicts, etc.)

func TestInteger(t *testing.T) {
	source := "i34e"
	expected := BValue{BTypes.Integer, 34}
	value, _ := Decode(source)
	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}
}

func TestString(t *testing.T) {
	source := "11:Hello World"
	expected := BValue{BTypes.String, "Hello World"}
	value, _ := Decode(source)

	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}
}

func TestList(t *testing.T) {
	source := "l4:spam4:eggse"
	expected := BValue{BTypes.List, []BValue{{BTypes.String, "spam"}, {BTypes.String, "eggs"}}}
	value, err := Decode(source)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected %v got %v, Err: %v", expected, value, err)
	}
}

func TestDict(t *testing.T) {
	source := "d3:cow3:moo4:spam4:eggse"
	expected := BValue{BTypes.Dict, map[string]BValue{"cow": {BTypes.String, "moo"}, "spam": {BTypes.String, "eggs"}}}
	value, err := Decode(source)

	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected %v got %v, Err: %v", expected, value, err)
	}
}

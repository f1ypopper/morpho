package bencoding_test

import (
	"morpho/bencoding"
	"reflect"
	"testing"
)

//TODO: implement more tests (nested lists, dicts, etc.)

func TestInteger(t *testing.T) {
	source := "i34e"
	expected := bencoding.BValue{bencoding.BTypes.Integer, 34}
	value, _ := bencoding.Decode(source)
	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}
}

func TestString(t *testing.T) {
	source := "11:Hello World"
	expected := bencoding.BValue{bencoding.BTypes.String, "Hello World"}
	value, _ := bencoding.Decode(source)

	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}
}

func TestList(t *testing.T) {
	source := "l4:spam4:eggse"
	expected := bencoding.BValue{bencoding.BTypes.List, []bencoding.BValue{{bencoding.BTypes.String, "spam"}, {bencoding.BTypes.String, "eggs"}}}
	value, err := bencoding.Decode(source)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected %v got %v, Err: %v", expected, value, err)
	}
}

func TestDict(t *testing.T) {
	source := "d3:cow3:moo4:spam4:eggse"
	expected := bencoding.BValue{bencoding.BTypes.Dict, map[string]bencoding.BValue{"cow": {bencoding.BTypes.String, "moo"}, "spam": {bencoding.BTypes.String, "eggs"}}}
	value, err := bencoding.Decode(source)

	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected %v got %v, Err: %v", expected, value, err)
	}
}

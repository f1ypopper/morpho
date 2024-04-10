package bencoding_test

import (
	"morpho/bencoding"
	"testing"
)

//TODO: implement all tests

func TestInteger(t *testing.T) {
	source := "i34e"
	expected := bencoding.BValue{bencoding.BTypes.Integer, 34}
	value, _ := bencoding.Decode(source)

	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}
}

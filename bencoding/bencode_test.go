package bencoding

import "testing"

func TestEncodeString(t *testing.T) {
	source := "shame"
	expected := "5:shame"
	value := Encode(source)
	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}

}
func TestEncodeInteger(t *testing.T) {
	source := 34
	expected := "i34e"
	value := Encode(source)
	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}

}
func TestEncodeDict(t *testing.T) {
	source := map[string]interface{}{
		"cow": "moo",
		"key": 23,
	}
	expected := "d3:cow3:moo3:keyi23ee"
	value := Encode(source)
	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}

}
func TestEncodeList(t *testing.T) {
	source := []interface{}{"hello", "cow", 23}
	expected := "l5:hello3:cowi23ee"
	value := Encode(source)
	if value != expected {
		t.Errorf("Expected %v got %v", expected, value)
	}

}

package bencoding

import (
	// "fmt"
	"strconv"
)

func (e *BEncoder) encode() string {
	switch e.source.(type) {
	case string:
		return e.encodeString()
	case int:
		return e.encodeInteger()

	case map[string]interface{}:
		return e.encodeDict()
	case []interface{}:
		return e.encodeList()
	}

	return ""

}

func (e *BEncoder) encodeString() string {
	var str string
	str = e.source.(string)
	out := strconv.Itoa(len(str)) + ":" + str
	return out

}
func (e *BEncoder) encodeInteger() string {
	str := e.source.(int)
	out := "i" + strconv.Itoa(str) + "e"
	return out

}
func (e *BEncoder) encodeDict() string {
	dict := e.source.(map[string]interface{})
	var out string
	var keyEncoder BEncoder
	var valEncoder BEncoder
	for k, v := range dict {
		keyEncoder.source = k
		valEncoder.source = v

		out += keyEncoder.encode()
		out += valEncoder.encode()

	}
	return "d" + out + "e"

}
func (e *BEncoder) encodeList() string {
	var out string
	list := e.source.([]interface{})
	var listEncoder BEncoder
	for _, v := range list {
		listEncoder.source = v
		out += listEncoder.encode()
	}
	return "l" + out + "e"

}

func Encode(source any) string {
	encoder := BEncoder{source: source}
	return encoder.encode()

}

package bencoding

type BType uint8

var BTypes = struct {
	String  BType
	Integer BType
	List    BType
	Dict    BType
}{
	String:  0,
	Integer: 1,
	List:    2,
	Dict:    3,
}

type BValue struct {
	Btype BType
	Value any
}

func Encode(bvalue BValue) string {
	panic("todo")
}

func Decode(source string) (BValue, error) {
	panic("todo")
}

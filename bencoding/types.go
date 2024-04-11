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

type BDecoder struct {
	current uint
	source  string
}
type BEncoder struct {
	source any
}

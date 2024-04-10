package bencoding

type BType uint8

const (
	String BType = iota
	Integer
	List
	Dict
)

type BValue struct {
	btype BType
	value any
}

func Encode(value BValue) string {
	panic("todo")
}

func Decode(source string) (BValue, error) {
	panic("todo")
}

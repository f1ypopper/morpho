package bencoding

type BDecoder struct {
	current uint
	source  string
}
type BEncoder struct {
	source any
}

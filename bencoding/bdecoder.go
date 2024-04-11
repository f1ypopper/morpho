package bencoding

import (
	"errors"
	"fmt"
	"strconv"
)

func (decoder *BDecoder) decode() (BValue, error) {
	c := decoder.advance()
	switch c {
	case 'i':
		return decoder.decodeInt()
	case 'l':
		return decoder.decodeList()
	case 'd':
		return decoder.decodeDict()
	default:
		if '0' <= c && c <= '9' {
			return decoder.decodeString()
		} else {
			return BValue{}, fmt.Errorf("unknown character %c", c)
		}
	}
}
func (decoder *BDecoder) advance() byte {
	decoder.current += 1
	return decoder.source[decoder.current-1]
}
func (decoder *BDecoder) peek() byte {
	if decoder.atEnd() {
		return 0x0
	}
	return decoder.source[decoder.current]
}
func (decoder *BDecoder) atEnd() bool {
	return decoder.current >= uint(len(decoder.source))
}
func (decoder *BDecoder) decodeInt() (BValue, error) {
	start := decoder.current
	for decoder.peek() != 'e' && !decoder.atEnd() {
		decoder.advance()
	}

	if decoder.atEnd() {
		return BValue{}, errors.New("unterminated integer")
	}
	integer, err := strconv.Atoi(decoder.source[start:decoder.current])

	if err != nil {
		return BValue{}, err
	}

	//consume the 'e'
	decoder.advance()
	return BValue{BTypes.Integer, integer}, nil
}

func (decoder *BDecoder) decodeString() (BValue, error) {
	start := decoder.current - 1
	for decoder.peek() != ':' {
		decoder.advance()
	}
	slen, err := strconv.ParseUint(decoder.source[start:decoder.current], 10, 64)
	if err != nil {
		return BValue{}, err
	}
	//consume the ':'
	decoder.advance()
	str := decoder.source[decoder.current : decoder.current+uint(slen)]

	decoder.current += uint(slen)

	return BValue{BTypes.String, str}, nil
}

func (decoder *BDecoder) decodeList() (BValue, error) {
	var list = make([]BValue, 0, 10)
	for decoder.peek() != 'e' {
		elem, err := decoder.decode()
		if err != nil {
			return BValue{}, err
		}
		list = append(list, elem)
	}
	//consume the 'e'
	decoder.advance()
	return BValue{BTypes.List, list}, nil
}

func (decoder *BDecoder) decodeDict() (BValue, error) {
	var dict = make(map[string]BValue)
	for decoder.peek() != 'e' {
		key, err := decoder.decode()
		if err != nil {
			return BValue{}, err
		}
		if key.Btype != BTypes.String {
			return BValue{}, errors.New("expected key to be a string.")
		}
		value, err := decoder.decode()
		if err != nil {
			return BValue{}, err
		}
		dict[key.Value.(string)] = value
	}
	//consume the 'e'
	decoder.advance()
	return BValue{BTypes.Dict, dict}, nil
}

func Decode(source string) (BValue, error) {
	decoder := BDecoder{0, source}
	return decoder.decode()
}

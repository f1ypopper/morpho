package bencoding

import (
	"errors"
	"fmt"
	"strconv"
)

func (decoder *BDecoder) decode() (any, error) {
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
			return nil, fmt.Errorf("unknown character %c", c)
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
func (decoder *BDecoder) decodeInt() (int, error) {
	start := decoder.current
	for decoder.peek() != 'e' && !decoder.atEnd() {
		decoder.advance()
	}

	if decoder.atEnd() {
		return 0, errors.New("unterminated integer")
	}
	integer, err := strconv.Atoi(decoder.source[start:decoder.current])

	if err != nil {
		return 0, err
	}

	//consume the 'e'
	decoder.advance()
	return integer, nil
}

func (decoder *BDecoder) decodeString() (string, error) {
	start := decoder.current - 1
	for decoder.peek() != ':' {
		decoder.advance()
	}
	slen, err := strconv.ParseUint(decoder.source[start:decoder.current], 10, 64)
	if err != nil {
		return "", err
	}
	//consume the ':'
	decoder.advance()
	str := decoder.source[decoder.current : decoder.current+uint(slen)]

	decoder.current += uint(slen)

	return str, nil
}

func (decoder *BDecoder) decodeList() ([]any, error) {
	var list = []any{}
	for decoder.peek() != 'e' {
		elem, err := decoder.decode()
		if err != nil {
			return nil, err
		}
		list = append(list, elem)
	}
	//consume the 'e'
	decoder.advance()
	return list, nil
}

func (decoder *BDecoder) decodeDict() (map[string]any, error) {
	var dict = make(map[string]any)
	for decoder.peek() != 'e' {
		key, err := decoder.decode()
		if err != nil {
			return nil, err
		}
		if _, ok := key.(string); !ok {
			return nil, errors.New("expected key to be a string.")
		}
		value, err := decoder.decode()
		if err != nil {
			return nil, err
		}
		dict[key.(string)] = value
	}
	//consume the 'e'
	decoder.advance()
	return dict, nil
}

func Decode(source string) (any, error) {
	decoder := BDecoder{0, source}
	return decoder.decode()
}

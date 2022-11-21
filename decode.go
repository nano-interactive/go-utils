package utils

import "errors"

func UnHex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func IsHex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

var ErrEscape = errors.New("invalid URL escape")

func DecodeQueryUnsafeString(src []byte) (string, error) {
	n, err := DecodeQuery(src, src)
	if err != nil {
		return "", err
	}

	return UnsafeString(src[:n]), nil
}

func DecodeQueryUnsafe(src []byte) ([]byte, error) {
	n, err := DecodeQuery(src, src)
	if err != nil {
		return nil, err
	}

	return src[:n], nil
}

func DecodeQuery(dst []byte, src []byte) (int, error) {
	writer := 0

	ptr := src
	writer = 0
	for reader := 0; reader < len(ptr); reader++ {
		switch ptr[reader] {
		case '%':
			// DO WHILE
			if reader+2 >= len(ptr) || !IsHex(ptr[reader+1]) || !IsHex(ptr[reader+2]) {
				return 0, ErrEscape
			}

			data := UnHex(ptr[reader+1])<<4 | UnHex(ptr[reader+2])
			reader += 2

			for data == '%' {
				if reader+2 >= len(ptr) || !IsHex(ptr[reader+1]) || !IsHex(ptr[reader+2]) {
					return 0, ErrEscape
				}

				data = UnHex(ptr[reader+1])<<4 | UnHex(ptr[reader+2])
				reader += 2
			}

			dst[writer] = data
		case '+':
			dst[writer] = ' '
		default:
			dst[writer] = ptr[reader]
		}

		writer++
	}

	return writer, nil
}

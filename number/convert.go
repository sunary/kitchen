package number

import (
	"strconv"
	"unicode/utf8"
)

const (
	GroupLen = 3
)

// ToThousandFormat To show your numbers in thousands
// ToThousandFormat(1234567, ',')
// => "1,234,567"
func ToThousandFormat(n int64, sep rune) string {
	s := strconv.FormatInt(n, 10)

	startOffset := 0
	if n < 0 {
		startOffset = 1
	}

	groups := (len(s) - startOffset - 1) / GroupLen

	if groups == 0 {
		return s
	}

	sepLen := utf8.RuneLen(sep)
	sepBytes := make([]byte, sepLen)
	_ = utf8.EncodeRune(sepBytes, sep)

	buf := make([]byte, groups*(GroupLen+sepLen)+len(s)-(groups*GroupLen))

	startOffset += GroupLen
	p := len(s)
	q := len(buf)
	for p > startOffset {
		p -= GroupLen
		q -= GroupLen
		copy(buf[q:q+GroupLen], s[p:])
		q -= sepLen
		copy(buf[q:], sepBytes)
	}

	if q > 0 {
		copy(buf[:q], s)
	}

	return string(buf)
}

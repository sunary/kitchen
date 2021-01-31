package number

import (
	"strconv"
	"unicode/utf8"
)

const (
	groupLen = 3
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

	groups := (len(s) - startOffset - 1) / groupLen

	if groups == 0 {
		return s
	}

	sepLen := utf8.RuneLen(sep)
	sepBytes := make([]byte, sepLen)
	_ = utf8.EncodeRune(sepBytes, sep)

	buf := make([]byte, groups*(groupLen+sepLen)+len(s)-(groups*groupLen))

	startOffset += groupLen
	p := len(s)
	q := len(buf)
	for p > startOffset {
		p -= groupLen
		q -= groupLen
		copy(buf[q:q+groupLen], s[p:])
		q -= sepLen
		copy(buf[q:], sepBytes)
	}

	if q > 0 {
		copy(buf[:q], s)
	}

	return string(buf)
}

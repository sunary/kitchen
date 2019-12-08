package str

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"sort"
	"strings"
	"unicode"
)

// ToTitle
// ToTitle("to_title")
// => "ToTitle"
func ToTitle(input string) string {
	var output strings.Builder
	ss := strings.Split(input, "_")
	for _, s := range ss {
		if s == "" {
			continue
		}

		output.WriteString(strings.ToUpper(string(s[0])) + s[1:])
	}

	return output.String()
}

// ToTitleNorm ...
func ToTitleNorm(input string) string {
	var output strings.Builder
	var upperCount int
	for i, c := range input {
		switch {
		case isUppercase(c):
			if upperCount == 0 || nextIsLower(input, i) {
				output.WriteByte(byte(c))
			} else {
				output.WriteByte(byte(c - 'A' + 'a'))
			}
			upperCount++

		case isLowercase(c):
			output.WriteByte(byte(c))
			upperCount = 0

		case isDigit(c):
			if i == 0 {
				panic("go-common/str: Identifier must start with a character: `" + input + "`")
			}
			output.WriteByte(byte(c))
			upperCount = 0
		}
	}

	return output.String()
}

// ToSnake
// ToSnake("ToSnake")
// => "to_snake"
func ToSnake(input string) string {
	var output strings.Builder
	var upperCount int
	for i, c := range input {
		switch {
		case isUppercase(c):
			if i > 0 && (upperCount == 0 || nextIsLower(input, i)) {
				output.WriteByte('_')
			}
			output.WriteByte(byte(c - 'A' + 'a'))
			upperCount++

		case isLowercase(c):
			output.WriteByte(byte(c))
			upperCount = 0

		case isDigit(c):
			if i == 0 {
				panic("go-common/str: Identifier must start with a character: `" + input + "`")
			}
			output.WriteByte(byte(c))

		default:
			panic("go-common/str: Invalid identifier: `" + input + "`")
		}
	}

	return output.String()
}

// nextIsLower The next character is lower case, but not the last 's'.
// nextIsLower("HTMLFile", 1) expected: "html_file"
// => true
// nextIsLower("URLs", -1) expected: "urls"
// => false
func nextIsLower(input string, i int) bool {
	i++
	if i >= len(input) {
		return false
	}

	c := input[i]
	if c == 's' && i == len(input)-1 {
		return false
	}

	return isLowercase(rune(c))
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLowercase(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isUppercase(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// Normalize ...
func Normalize(input string) string {
	input = strings.TrimSpace(input)

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	strTransform, _, _ := transform.String(t, input)

	sortedSpecialRunes := []rune{'Đ', 'đ', 'Ł'}
	replacedByRunes := []rune{'D', 'd', 'L'}
	var result strings.Builder

	for _, r := range strTransform {
		pos := sort.Search(len(sortedSpecialRunes), func(i int) bool { return sortedSpecialRunes[i] >= r })
		if pos != -1 && r == sortedSpecialRunes[pos] {
			result.WriteRune(replacedByRunes[pos])
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// HashPassword ...
func HashPassword(password, salt []byte) string {
	mac := hmac.New(sha256.New, salt)
	mac.Write([]byte(password))
	return hex.EncodeToString(mac.Sum(nil))
}

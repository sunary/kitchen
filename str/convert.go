package str

import (
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Coalesce Returns the first non-empty string value
func Coalesce(input ...string) string {
	for _, s := range input {
		if s != "" {
			return s
		}
	}

	return ""
}

// ToTitleCase Convert snake_case to TitleCase
// ToTitleCase("to_title_case")
// => "ToTitleCase"
func ToTitleCase(input string) string {
	var sb strings.Builder
	ss := strings.Split(input, "_")
	for _, s := range ss {
		if s == "" {
			continue
		}

		sb.WriteString(strings.ToUpper(string(s[0])) + s[1:])
	}

	return sb.String()
}

// ToCamelCase Convert TitleCase to camelCase
// ToCamelCase("ToCamelCase")
// => "toCamelCase"
func ToCamelCase(input string) string {
	if input == "" {
		return ""
	}

	a := []rune(input)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

// NormTitleCase Normalize TitleCase
// NormTitleCase("GRPCError")
// => "GrpcError"
func NormTitleCase(input string) string {
	var sb strings.Builder
	var upperCount int
	for i, c := range input {
		switch {
		case isUppercase(c):
			if upperCount == 0 || nextIsLower(input, i) {
				sb.WriteByte(byte(c))
			} else {
				sb.WriteByte(byte(c - 'A' + 'a'))
			}
			upperCount++

		case isLowercase(c):
			sb.WriteByte(byte(c))
			upperCount = 0

		case isDigit(c):
			if i == 0 {
				panic("go-common/str: Identifier must start with a character: `" + input + "`")
			}
			sb.WriteByte(byte(c))
			upperCount = 0
		}
	}

	return sb.String()
}

// ToSnakeCase Convert TitleCase to snake_case
// ToSnakeCase("ToSnakeCase")
// => "to_snake_case"
func ToSnakeCase(input string) string {
	var sb strings.Builder
	var upperCount int
	for i, c := range input {
		switch {
		case isUppercase(c):
			if i > 0 && (upperCount == 0 || nextIsLower(input, i)) {
				sb.WriteByte('_')
			}
			sb.WriteByte(byte(c - 'A' + 'a'))
			upperCount++

		case isLowercase(c):
			sb.WriteByte(byte(c))
			upperCount = 0

		case isDigit(c):
			if i == 0 {
				panic("go-common/str: Identifier must start with a character: `" + input + "`")
			}
			sb.WriteByte(byte(c))

		default:
			panic("go-common/str: Invalid identifier: `" + input + "`")
		}
	}

	return sb.String()
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

// Normalize Remove diacritics
func Normalize(input string) string {
	input = strings.TrimSpace(input)

	rMapping := func(r rune) rune {
		sortedSpecialRunes := []rune{'Đ', 'đ', 'Ł'}
		replacedByRunes := []rune{'D', 'd', 'L'}

		pos := sort.Search(len(sortedSpecialRunes), func(i int) bool { return sortedSpecialRunes[i] >= r })
		if pos != -1 && r == sortedSpecialRunes[pos] {
			return replacedByRunes[pos]
		}
		return r
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC, runes.Map(rMapping))
	strTransform, _, _ := transform.String(t, input)
	return strTransform
}

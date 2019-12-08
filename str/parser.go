package str

import (
	"unicode"
)

// IsMediumPassword ...
func IsMediumPassword(p string) bool {
	if len(p) < 8 {
		return false
	}

	has := map[string]bool{
		"upper":   false,
		"lower":   false,
		"numeric": false,
	}

	for name, classes := range map[string][]*unicode.RangeTable{
		"upper":   {unicode.Upper, unicode.Title},
		"lower":   {unicode.Lower},
		"numeric": {unicode.Number, unicode.Digit},
	} {
		for _, r := range p {
			if !has[name] && unicode.IsOneOf(classes, r) {
				has[name] = true
			}
		}
	}

	return has["upper"] && has["lower"] && has["numeric"]
}

// IsStrongPassword ...
func IsStrongPassword(p string) bool {
	if len(p) < 8 {
		return false
	}

	has := map[string]bool{
		"upper":   false,
		"lower":   false,
		"numeric": false,
		"special": false,
	}

	for name, classes := range map[string][]*unicode.RangeTable{
		"upper":   {unicode.Upper, unicode.Title},
		"lower":   {unicode.Lower},
		"numeric": {unicode.Number, unicode.Digit},
		"special": {unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range p {
			if !has[name] && unicode.IsOneOf(classes, r) {
				has[name] = true
			}
		}
	}

	return has["upper"] && has["lower"] && has["numeric"] && has["special"]
}

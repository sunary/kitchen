package str

import (
	"unicode"
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

// IsMediumPassword Must contain lower, upper and numeric
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

// IsStrongPassword Must contain lower, upper, numeric and special characters
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

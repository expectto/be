// Package be_string provides Be matchers for string-related assertions.
package be_string

import (
	"fmt"
	"github.com/IGLOU-EU/go-wildcard" // used specifically for MatchWildcard matcher
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// todo: do a transform Stringish->String that is configurable (string/fuzzy strings)
//
//	and use it in any following string matcher
//
// Deprecated
var expectAvailableStringFormat = func(actual any) error {
	if !cast.IsString(actual, cast.AllowCustomTypes(), cast.AllowPointers()) {
		return fmt.Errorf("string expected, got %T", actual)
	}

	return nil
}

// NonEmptyString succeeds if actual is not an empty string.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func NonEmptyString() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		psi_matchers.NewNotMatcher(gomega.BeEmpty()),
	), "be non-empty string")
}

// EmptyString succeeds if actual is an empty string.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func EmptyString() types.BeMatcher {
	return Psi(psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		gomega.BeEmpty(),
	), "be an empty string")
}

// Alpha succeeds if actual is a string containing only alphabetical characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Alpha() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered Alpha
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it contains only letters
		for _, char := range str {
			if !unicode.IsLetter(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetical characters")
}

// AlphaWhitespace succeeds if actual is a string containing only alphabetical and whitespace characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaWhitespace() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaWithSpace
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it contains only letters
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetical and whitespace characters")
}

// AlphaWithPunctuation succeeds if actual is a string containing only alphabetical characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaWithPunctuation() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaWithPunctuation
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it contains only letters
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsPunct(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetical and punctuational characters")
}

// AlphaWhitespaceWithPunctuation succeeds if actual is a string containing only alphabetical characters with whitespace and punctuation..
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaWhitespaceWithPunctuation() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaWithPunctuation
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it contains only letters
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsPunct(char) && !unicode.IsSpace(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetical characters and punctuation")
}

// Whitespace succeeds if actual is a string containing only whitespace characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Whitespace() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaWithPunctuation
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it contains only letters
		for _, char := range str {
			if !unicode.IsSpace(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only whitespace characters")
}

// Numeric succeeds if actual is a string representing a valid numeric integer.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Numeric() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered Numeric
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it's a numeric string
		_, err := strconv.ParseInt(str, 10, 64)
		return err == nil, nil
	}, "contain only numeric characters")
}

// NumericWhitespace succeeds if actual is a string containing only number characters and whitespace.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func NumericWhitespace() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered Numeric
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it contains only letters
		for _, char := range str {
			if !unicode.IsSpace(char) && !unicode.IsNumber(char) {
				return false, nil
			}
		}
		return true, nil
	}, "contain only numeric characters and whitespace")
}

// AlphaNumeric succeeds if actual is a string containing only alphanumeric characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
// As Numeric() matcher is considered to match on integers, AlphaNumeric() doesn't match on dots
// So, consider AlphaNumericWithDots() then
func AlphaNumeric() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaNumeric
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it's an alphanumeric string
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetic and numeric characters")
}

// AlphaNumericWhitespace succeeds if actual is a string containing only alphanumeric characters and whitespace.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaNumericWhitespace() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaNumericWhitespace
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it's an alphanumeric string
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && !unicode.IsSpace(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetic, numeric characters and punctuation")
}

// AlphaNumericWithPunctuation succeeds if actual is a string containing only alphanumeric characters and punctuation.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaNumericWithPunctuation() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaNumericWithPunctuation
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it's an alphanumeric string
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && !unicode.IsPunct(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetic, numeric characters and punctuation")
}

// AlphaNumericWhitespaceWithPunctuation succeeds if actual is a string containing alphanumeric, whitespace and punctuation.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaNumericWhitespaceWithPunctuation() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaNumericWithPunctuation
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it's an alphanumeric string
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && !unicode.IsPunct(char) && !unicode.IsSpace(char) {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphanumeric characters with whitespace and punctuation")
}

// AlphaNumericWithDots succeeds if actual is a string containing only alphanumeric characters and dots.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaNumericWithDots() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered AlphaNumericWithDots
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// Check if it's an alphanumeric string
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '.' {
				return false, nil
			}
		}

		return true, nil
	}, "contain only alphabetic, numeric characters and dot")
}

// Float succeeds if actual is a string representing a valid floating-point number.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Float() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// Check if it's a numeric string
		_, err := strconv.ParseFloat(cast.AsString(actual), 64)
		return err == nil, nil
	}, "be a string representation of a float value")
}

// Titled succeeds if actual is a string with the first letter of each word capitalized.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Titled(languageArg ...language.Tag) types.BeMatcher {
	lang := language.English
	if len(languageArg) > 0 {
		lang = languageArg[0]
	}

	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		str := cast.AsString(actual)

		return cases.Title(lang).String(str) == str, nil
	}, "be titled")
}

// LowerCaseOnly succeeds if actual is a string containing only lowercase characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func LowerCaseOnly() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		str := cast.AsString(actual)

		return strings.ToLower(str) == str, nil
	}, "be lower-case")
}

// UpperCaseOnly succeeds if actual is a string containing only uppercase characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func UpperCaseOnly() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		str := cast.AsString(actual)

		return strings.ToUpper(str) == str, nil
	}, "be upper-case")
}

// ContainingSubstring succeeds if actual is a string containing only characters from a given set
func ContainingSubstring(substr string) types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered ContainsOf
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		return strings.Contains(str, substr), nil

	}, fmt.Sprintf("contain `%s` substring", substr))
}

// ContainingOnlyCharacters succeeds if actual is a string containing only characters from a given set
func ContainingOnlyCharacters(characters string) types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// empty string is not considered ContainsOf
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// string -> map[rune]struct{}{} as a lookup-table
		allowedSet := make(map[rune]struct{})
		for _, char := range characters {
			allowedSet[char] = struct{}{}
		}

		// Check if it's an alphanumeric string
		for _, char := range str {
			if _, ok := allowedSet[char]; !ok {
				return false, nil
			}
		}

		return true, nil
	}, fmt.Sprintf("contain only `%s` characters", characters))
}

// ContainingCharacters succeeds if actual is a string containing all characters from a given set
func ContainingCharacters(characters string) types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}
		if len(characters) == 0 {
			return true, nil
		}

		// empty string is not considered ContainingCharacters
		str := cast.AsString(actual)
		if str == "" {
			return false, nil
		}

		// string -> map[rune]struct{}{} as a lookup-table
		requiredSet := make(map[rune]struct{})
		for _, char := range characters {
			requiredSet[char] = struct{}{}
		}

		actualSet := make(map[rune]struct{})
		for _, char := range str {
			actualSet[char] = struct{}{}
		}

		// Check if it's an alphanumeric string
		for _, requiredChar := range characters {
			if _, ok := actualSet[requiredChar]; !ok {
				return false, nil
			}
		}

		return true, nil
	}, fmt.Sprintf("contain all of `%s` characters", characters))
}

// MatchWildcard succeeds if actual matches given wildcard pattern.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func MatchWildcard(pattern string) types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		return wildcard.Match(pattern, cast.AsString(actual)), nil
	}, fmt.Sprintf("match wildcard %s", pattern))
}

// ValidEmail succeeds if actual is a valid email.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func ValidEmail() types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		_, err := mail.ParseAddress(cast.AsString(actual))
		return err == nil, nil

	}, fmt.Sprintf("be a valid email"))
}

// MatchTemplate succeeds if actual matches given template pattern.
// Provided template must have `{{Field}}` placeholders.
// Each distinct placeholder from template requires a var to be passed in list of `vars`.
// Var can be a raw value or a matcher
//
// E.g.
//
//	Expect(someString).To(be_string.MatchTemplate("Hello {{Name}}. Your number is {{Number}}", be_string.Var("Name", "John"), be_string.Var("Number", 3)))
//	Expect(someString).To(be_string.MatchTemplate("Hello {{Name}}. Good bye, {{Name}}.", be_string.Var("Name", be_string.Titled()))
func MatchTemplate(template string, vars ...*V) types.BeMatcher {
	return Psi(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// Idea here is to switch from templating to regexp (Ugly, but ok for first attempt)
		// {{Name}} => (?P<Name>.+)
		variableRegex := regexp.MustCompile(`{{\s*([^}\s]+)\s*}}`)
		regexStr := variableRegex.ReplaceAllString(template, "(?P<$1>.+)")

		regex, err := regexp.Compile(regexStr)
		if err != nil {
			return false, fmt.Errorf("bad template: %w", err)
		}

		match := regex.FindStringSubmatch(cast.AsString(actual))
		if len(match) != len(regex.SubexpNames()) {
			// todo: provide better error handling
			return false, nil
		}

		results := make(map[string]string)
		for i, name := range regex.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}
			name = strings.ToLower(name)

			if savedResult, ok := results[name]; ok {
				if savedResult != match[i] {
					return false, fmt.Errorf("var %s has multiple values: %s != %s", name, savedResult, match[i])
				}
			}

			results[name] = match[i]
		}

		// if no vars are given: we simply verified that whole string matches template
		// without matching specifically templates variables
		if len(vars) == 0 {
			return true, nil
		}

		for _, v := range vars {
			name := strings.ToLower(v.Name)
			result, ok := results[name]
			if !ok {
				return false, fmt.Errorf("var %s given but not met in actual value", name)
			}

			if matched, err := v.Matcher.Match(result); err != nil {
				return false, fmt.Errorf("var %s failed: %w", name, err)
			} else if !matched {
				// todo: transmit failure to the error message
				return false, nil
			}
		}

		return true, nil
	}, "match given template")
	// todo: it hides the actual reason of failing from template variables
	//       should be exposed
}

type V struct {
	Name    string
	Matcher types.BeMatcher
}

// Var creates a var used for replacing placeholders for templates in `MatchTemplate`
func Var(name string, matching any) *V {
	return &V{Name: name, Matcher: Psi(matching)}
}

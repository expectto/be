// Package be_strings provides Be matchers for string-related assertions
package be_strings

import (
	"fmt"
	"github.com/IGLOU-EU/go-wildcard" // used specifically for MatchWildcard matcher
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
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
	return psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		psi_matchers.NewNotMatcher(gomega.BeEmpty()),
	)
}

// EmptyString succeeds if actual is an empty string.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func EmptyString() types.BeMatcher {
	return psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		Psi(gomega.BeEmpty()),
	)
}

// Alpha succeeds if actual is a string containing only alphabetical characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Alpha() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// Check if it contains only letters
		for _, char := range cast.AsString(actual) {
			if !unicode.IsLetter(char) {
				return false, nil
			}
		}

		return true, nil
	}))
}

// Numeric succeeds if actual is a string representing a valid numeric integer.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Numeric() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// Check if it's a numeric string
		_, err := strconv.ParseInt(cast.AsString(actual), 10, 64)
		return err == nil, nil
	}))
}

// AlphaNumeric succeeds if actual is a string containing only alphanumeric characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
// As Numeric() matcher is considered to match on integers, AlphaNumeric() doesn't match on dots
// So, consider AlphaNumericWithDots() then
func AlphaNumeric() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// Check if it's an alphanumeric string
		for _, char := range cast.AsString(actual) {
			if !(('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9')) {
				return false, nil
			}
		}

		return true, nil
	}))
}

// AlphaNumericWithDots succeeds if actual is a string containing only alphanumeric characters and dots.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func AlphaNumericWithDots() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// Check if it's an alphanumeric string
		for _, char := range cast.AsString(actual) {
			if !(('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || char != '.') {
				return false, nil
			}
		}

		return true, nil
	}))
}

// Float succeeds if actual is a string representing a valid floating-point number.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Float() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// Check if it's a numeric string
		_, err := strconv.ParseFloat(cast.AsString(actual), 64)
		return err == nil, nil
	}))
}

// Titled succeeds if actual is a string with the first letter of each word capitalized.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Titled() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		str := cast.AsString(actual)

		// todo: switch to cases
		return strings.Title(str) == str, nil
	}))
}

// LowerCaseOnly succeeds if actual is a string containing only lowercase characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func LowerCaseOnly() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		str := cast.AsString(actual)

		return strings.ToLower(str) == str, nil
	}))
}

// MatchWildcard succeeds if actual matches given wildcard pattern.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func MatchWildcard(pattern string) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		return wildcard.Match(pattern, cast.AsString(actual)), nil
	}))
}

// ValidEmail succeeds if actual is a valid email.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func ValidEmail() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		_, err := mail.ParseAddress(cast.AsString(actual))
		// todo: do not lose the err
		return err == nil, nil
	}))
}

// MatchTemplate succeeds if actual matches given template pattern.
// Provided template must have `{{Field}}` placeholders.
// Each distinct placeholder from template requires a var to be passed in list of `vars`.
// Var can be a raw value or a matcher
//
// E.g.
//
//	Expect(someString).To(be_strings.MatchTemplate("Hello {{Name}}. Your number is {{Number}}", be_strings.Var("Name", "John"), be_strings.Var("Number", 3)))
//	Expect(someString).To(be_strings.MatchTemplate("Hello {{Name}}. Good bye, {{Name}}.", be_strings.Var("Name", be_strings.Titled()))
func MatchTemplate(template string, vars ...*V) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
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
	}))
}

type V struct {
	Name    string
	Matcher types.BeMatcher
}

// Var creates a var used for replacing placeholders for templates in `MatchTemplate`
func Var(name string, matching any) *V {
	return &V{Name: name, Matcher: Psi(matching)}
}

// Package be_strings provides Be matchers for string-related assertions.
// It includes some experimental matchers like Alpha/UpperCaseOnly/LowerCaseOnly/etc
//
// Note: The package relies on the IGLOU-EU/go-wildcard libraries for wildcard mathcing.
package be_strings

import (
	"fmt"
	"github.com/IGLOU-EU/go-wildcard"
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"net/mail"
	"regexp"
	"strings"
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

// NonEmptyString succeeds if actual is an empty string.
// Actual must be of string kind (can be adjusted via SetStringFormat method)
func NonEmptyString() types.BeMatcher {
	return psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		psi_matchers.NewNotMatcher(gomega.BeEmpty()),
	)
}

func EmptyString() types.BeMatcher {
	return psi_matchers.NewAllMatcher(
		be_reflected.AsString(),
		Psi(gomega.BeEmpty()),
	)
}

// Alpha succeeds if actual contains only letters
// Actual must be string (pointers and custom types are OK)
func Alpha() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		for _, char := range cast.AsString(actual) {
			if char < 'a' || char > 'Z' {
				return false, nil
			}
		}

		return true, nil
	}))
}

// Num succeeds if actual contains only numeric values (with dots)
func Num() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// todo:
		// check if it's a numeric stirng

		return true, nil
	}))
}

// AlphaNum succeeds if actual contains only numeric values (with dots)
func AlphaNum() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		// todo:
		// check if it's a numeric stirng + alpha

		return true, nil
	}))
}

func Titled() types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		str := cast.AsString(actual)

		return strings.ToTitle(str) == str, nil
	}))
}

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
// Actual must be string (pointers and custom types are OK)
func MatchWildcard(pattern string) types.BeMatcher {
	return Psi(gcustom.MakeMatcher(func(actual interface{}) (bool, error) {
		if err := expectAvailableStringFormat(actual); err != nil {
			return false, err
		}

		return wildcard.Match(pattern, cast.AsString(actual)), nil
	}))
}

// ValidEmail succeeds if actual is a valid email
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
// Provided template must have `{{Something}}` placeholders.
// Each distinct placeholder from template requires an arg to be passed in list of `args`.
// Arg can be a raw value or a matcher
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
		// special handling for core variables like `...`, `...?`
		// we have to rename them for valid regex
		for _, cv := range coreVariables {
			regexStr = strings.ReplaceAll(regexStr, "?P<"+cv.Name+">.", "?P<"+cv.Placeholder+">.")
		}

		regex, err := regexp.Compile(regexStr)
		if err != nil {
			return false, fmt.Errorf("bad template: %w", err)
		}

		match := regex.FindStringSubmatch(cast.AsString(actual))
		if len(match) != len(regex.SubexpNames()) {
			return false, fmt.Errorf("invalid template")
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
		for _, v := range coreVariables {
			name := v.Placeholder
			result, ok := results[name]
			if !ok {
				continue
			}

			if matched, err := v.Matcher.Match(result); err != nil {
				return false, fmt.Errorf("core var %s failed: %w", name, err)
			} else if !matched {
				// todo: transmit failure to the error message
				return false, nil
			}
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

type coreVariable struct {
	*V
	Placeholder string
}

func Var(name string, matching any) *V {
	return &V{Name: name, Matcher: Psi(matching)}
}

var coreVariables = []*coreVariable{
	{V: &V{Name: "...", Matcher: NonEmptyString()}, Placeholder: "__ANYTHING__"},
	{V: &V{Name: "..?", Matcher: NonEmptyString()}, Placeholder: "__ANYTHING_OPTIONAL__"},
}

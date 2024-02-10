// Package be_string provides Be matchers for string-related assertions.
package be_string

import (
	"fmt"
	"github.com/IGLOU-EU/go-wildcard" // used specifically for MatchWildcard matcher
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	. "github.com/expectto/be/options"
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

// psiString is a convenient wrapper for creating every be_string matcher
// It simply wraps `Psi()` call to have a pre-matching for `checking type string`
// and returning a separate clear message if `actual` is not of string type.
var psiString = func(args ...any) types.BeMatcher {
	return Psi(
		Psi(func(actual any) (bool, error) {
			// todo: IsString options should be configurable
			return cast.IsString(actual, cast.AllowCustomTypes()), nil
		}, "be type of string"),
		Psi(args...),
	)
}

// validateStringOption checks if string options satisfies given rune
func validateStringOption(opt StringOption, r rune) bool {
	switch opt {
	case Alpha:
		return unicode.IsLetter(r)
	case Numeric:
		return unicode.IsNumber(r)
	case Whitespace:
		return unicode.IsSpace(r)
	case Punctuation:
		return unicode.IsPunct(r)
	case Dots:
		return r == '.'
	case SpecialCharacters:
		// todo: implement
		return false
	default:
		return false
	}
}

// Only succeeds if actual is a string containing only characters described by given options
// Only() defaults to empty string matching
// Only(Alpha|Numeric) succeeds if string contains only from alphabetic and numeric characters
// Available options are: Alpha, Numeric, Whitespace, Dots, Punctuation, SpecialCharacters
// TODO: special-characters are not supported yet
func Only(option StringOption) types.BeMatcher {
	if option == 0 {
		return EmptyString()
	}

	options := ExtractStringOptions(option)

	// We need stringified version of all options for the failure message
	optionsStr := make([]string, len(options), len(options))
	for i := range options {
		optionsStr[i] = options[i].String()
	}
	return psiString(func(actual any) (bool, error) {
		str := cast.AsString(actual)

		// empty string is not consider as any of string options
		if str == "" {
			return false, nil
		}

		// Check if it contains only letters
		for _, char := range str {
			// we're OK until there is a char that doesn't satisfy ANY option:
			var valid = false
			for _, opt := range options {
				valid = valid || validateStringOption(opt, char)
			}
			if !valid {
				return false, nil
			}
		}

		return true, nil
	}, fmt.Sprintf("contain only %s characters", strings.Join(optionsStr, "|")))
}

// EmptyString succeeds if actual is an empty string.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func EmptyString() types.BeMatcher {
	return psiString(gomega.BeEmpty(), "be an empty string")
}

// NonEmptyString succeeds if actual is not an empty string.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func NonEmptyString() types.BeMatcher {
	return psiString(gomega.Not(gomega.BeEmpty()), "be a non-empty string")
}

// Float succeeds if actual is a string representing a valid floating-point number.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func Float() types.BeMatcher {
	return psiString(func(actual any) (bool, error) {
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

	return psiString(func(actual any) (bool, error) {
		str := cast.AsString(actual)
		return cases.Title(lang).String(str) == str, nil
	}, "be a titled string")
}

// LowerCaseOnly succeeds if actual is a string containing only lowercase characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func LowerCaseOnly() types.BeMatcher {
	return psiString(func(actual any) (bool, error) {
		str := cast.AsString(actual)
		return strings.ToLower(str) == str, nil
	}, "be lower-case")
}

// UpperCaseOnly succeeds if actual is a string containing only uppercase characters.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func UpperCaseOnly() types.BeMatcher {
	return psiString(func(actual any) (bool, error) {
		str := cast.AsString(actual)
		return strings.ToUpper(str) == str, nil
	}, "be upper-case")
}

// ContainingSubstring succeeds if actual is a string containing only characters from a given set
func ContainingSubstring(substr string) types.BeMatcher {
	return psiString(func(actual any) (bool, error) {
		return strings.Contains(cast.AsString(actual), substr), nil
	}, fmt.Sprintf(`contain "%s" substring`, substr))
}

// ContainingOnlyCharacters succeeds if actual is a string containing only characters from a given set
func ContainingOnlyCharacters(characters string) types.BeMatcher {
	return psiString(func(actual any) (bool, error) {
		// empty string is not considered ContainsOf
		str := cast.AsString(actual)
		if str == "" || characters == "" {
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
	return psiString(func(actual any) (bool, error) {
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
	return psiString(func(actual any) (bool, error) {
		return wildcard.Match(pattern, cast.AsString(actual)), nil
	}, fmt.Sprintf(`match wildcard "%s"`, pattern))
}

// ValidEmail succeeds if actual is a valid email.
// Actual must be a string-like value (can be adjusted via SetStringFormat method).
func ValidEmail() types.BeMatcher {
	return psiString(func(actual any) (bool, error) {
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
	// Idea here is to switch from templating to regexp (Ugly, but ok for first attempt)
	// {{Name}} => (?P<Name>.+)
	variableRegex := regexp.MustCompile(`{{\s*([^}\s]+)\s*}}`)
	regexStr := variableRegex.ReplaceAllString(template, "(?P<$1>.+)")

	regex, err := regexp.Compile(regexStr)
	if err != nil {
		panic("invalid template: could not compile a regex from it: " + err.Error())
	}

	return psiString(func(actual any) (bool, error) {
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

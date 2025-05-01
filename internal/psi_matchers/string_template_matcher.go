package psi_matchers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/amberpixels/abu/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/format"
)

type Value struct {
	Name    string
	Matcher types.BeMatcher
}

type Values []*Value

// V creates a value used for replacing placeholders for templates in `MatchTemplate`
func V(name string, matching any) *Value {
	return &Value{Name: name, Matcher: Psi(matching)}
}

type StringTemplateMatcher struct {
	*MixinMatcherGomock

	regex  *regexp.Regexp
	values Values

	lastValue     *Value
	lastActual    any
	failedMessage string
}

var _ types.BeMatcher = &StringTemplateMatcher{}

func NewStringTemplateMatcher(template string, values ...*Value) *StringTemplateMatcher {
	// very quick and dirty check
	// We don't allow `{{Greeting}}{{Username}}`. We expect at least any separator between placeholders:
	if strings.Contains(template, "}}{{") {
		panic("invalid template: Placeholders can't be concatenated without separators")
	}

	// Idea here is to switch from templating to regexp (Ugly, but ok for first attempt)
	// {{Name}} => (?P<Name>.+)
	variableRegex := regexp.MustCompile(`{{\s*([^}\s]+)\s*}}`)
	regexStr := variableRegex.ReplaceAllString(template, "(?P<$1>.+)")

	regex, err := regexp.Compile(regexStr)
	if err != nil {
		panic("invalid template: could not compile a regex from it: " + err.Error())
	}

	return &StringTemplateMatcher{
		regex: regex, values: values,
	}
}

func (matcher *StringTemplateMatcher) Match(actual any) (success bool, err error) {
	match := matcher.regex.FindStringSubmatch(cast.AsString(actual))
	if len(match) != len(matcher.regex.SubexpNames()) {
		matcher.failedMessage = fmt.Sprintf(
			"initial mismatch: number of groups expected to be %d but not %d",
			len(match),
			len(matcher.regex.SubexpNames()),
		)
		return false, nil
	}

	results := make(map[string]string)
	for i, name := range matcher.regex.SubexpNames() {
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
	if len(matcher.values) == 0 {
		return true, nil
	}

	for _, v := range matcher.values {
		name := strings.ToLower(v.Name)
		result, ok := results[name]
		if !ok {
			return false, fmt.Errorf("var %s given but not met in actual value", name)
		}

		matched, err := v.Matcher.Match(result)
		if err != nil {
			return false, fmt.Errorf("var %s failed: %w", name, err)
		}
		matcher.lastValue = v
		matcher.lastActual = result

		if !matched {
			return false, nil
		}
	}

	return true, nil
}

func (matcher *StringTemplateMatcher) FailureMessage(actual any) string {
	if matcher.lastValue == nil {
		return format.Message(actual, "to match template:\n %s", matcher.failedMessage)
	}

	return format.Message(
		actual,
		fmt.Sprintf(
			"to match template on value %s:\n %s",
			matcher.lastValue.Name,
			matcher.lastValue.Matcher.FailureMessage(matcher.lastActual),
		),
	)
}

func (matcher *StringTemplateMatcher) NegatedFailureMessage(actual any) string {
	if matcher.lastValue == nil {
		return format.Message(actual, "not to match template:\n %s", matcher.failedMessage)
	}

	return format.Message(
		actual,
		fmt.Sprintf(
			"not to match template on value: %s:\n %s",
			matcher.lastValue.Name,
			matcher.lastValue.Matcher.NegatedFailureMessage(actual),
		),
	)
}

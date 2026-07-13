package be_struct

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/expectto/be/internal/psi" //nolint:staticcheck // should be moved to lintignore
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/gcustom"
)

// HavingField succeeds if the actual value is a struct of exactly type StructT
// and it has a field with the given name. If an expected value (raw value or
// matcher) is provided, the field's value must match it too.
//
// Example:
//
//	Expect(result).To(be_struct.HavingField[TestStruct]("Field1", "hello1"))
//	Expect(result).To(be_struct.HavingField[TestStruct]("Field2"))
//
// This is the generic, compile-time-typed variant: it pins the struct type at
// the call site. For everyday assertions prefer the root be.HaveField — it
// needs no type parameter, and supports dotted paths and "Method()" specs.
// The naming wobble (HaveField vs HavingField) is deliberate: it separates
// the default matcher from this typed variant.
func HavingField[StructT any](fieldName string, expectedValue ...any) types.BeMatcher {
	message := "have field " + fieldName
	if len(expectedValue) > 0 {
		message += fmt.Sprintf(" with value\n\t<%T>: %v", expectedValue[0], expectedValue[0])
	}

	structType := reflect.TypeFor[StructT]()

	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		val := reflect.ValueOf(actual)

		// Dereference pointer(s) to a struct
		for val.Kind() == reflect.Pointer {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return false, errors.New("actual value is not a struct")
		}
		// Wrong struct type for the [StructT] parameter is a clean non-match (so it
		// composes with Not); the previous code mutated the shared message here,
		// corrupting the matcher's description — that side-effect is removed.
		if val.Type() != structType {
			return false, nil
		}

		// Check if the field exists
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return false, nil
		}

		// No expected value: succeed as long as the field exists.
		if len(expectedValue) == 0 {
			return true, nil
		}

		// The expected value may itself be a matcher (be/gomega/gomock) — match the
		// field against it; otherwise compare by deep equality.
		if IsMatcher(expectedValue[0]) {
			return Psi(expectedValue[0]).Match(field.Interface())
		}
		return reflect.DeepEqual(field.Interface(), expectedValue[0]), nil
	}), message)
}

package be_struct

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega/gcustom"
)

// HavingField succeeds if the actual value is a struct and it has a field with the given name.
// If an expected value is provided, it also succeeds if the actual value's field has the same value.
//
// Example:
//
//	Expect(result).To(be_structs.HavingField[TestStruct]("Field1", "hello1"))
//	Expect(result).To(be_structs.HavingField[TestStruct]("Field2"))
func HavingField[StructT any](fieldName string, expectedValue ...any) types.BeMatcher {
	message := "have field " + fieldName
	if len(expectedValue) > 0 {
		message += fmt.Sprintf(" with value\n\t<%T>: %v", expectedValue[0], expectedValue[0])
	}

	structType := reflect.TypeFor[StructT]()

	return Psi(gcustom.MakeMatcher(func(actual any) (bool, error) {
		val := reflect.ValueOf(actual)

		// Dereference if it's a pointer to a struct
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return false, errors.New("actual value is not a struct")
		}
		if val.Type() != structType {
			// TODO: it doesn't work
			message = fmt.Sprintf("be type of %s", structType.String())
			return false, nil
		}

		// Check if the field exists
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return false, nil
		}

		// If an expected value is provided, compare the field's value with it
		if len(expectedValue) > 0 {
			expected := reflect.ValueOf(expectedValue[0])
			if !reflect.DeepEqual(field.Interface(), expected.Interface()) {
				return false, nil
			}
		}

		// If no value to compare, return true if the field exists
		return true, nil
	}), message)
}

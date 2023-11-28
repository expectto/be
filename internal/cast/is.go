package cast

import (
	"reflect"
)

// IsNil checks if the given input is a nil value.
func IsNil(a any) bool {
	if a == nil {
		return true
	}

	v := reflect.ValueOf(a)
	k := v.Kind()
	// check if v is OK to be used with IsNil
	if k == reflect.Ptr || k == reflect.Interface || k == reflect.Chan ||
		k == reflect.Func || k == reflect.Map || k == reflect.Slice {
		return v.IsNil()
	}

	return false
}

// IsStringish checks if the given input is a string or string-like value.
// To prevent code duplication, it employs panic recovery to handle type conversion
// and is designed for use in testing code, where panics are acceptable.
//
// Example Usage:
//
//	IsStringish("example") // Returns true
//	IsStringish([]byte("example")) // Returns true
//	IsStringish(CustomStringType("example")) // Returns true
//
// This function is suitable for scenarios where you want to quickly determine if
// a value can be treated as a string without handling detailed conversion errors.
func IsStringish(a any) (ok bool) {
	ok = true
	defer func() {
		if err := recover(); err != nil {
			ok = false
		}
	}()

	// here actually doesn't matter if we call AsBytes or AsString
	_ = AsBytes(a)
	return
}

// IsTime checks if the given input is a time.Time value (pointers and/or custom types are OK)
// To prevent code duplication, it employs panic recovery to handle type conversion
// and is designed for use in testing code, where panics are acceptable.
func IsTime(a any) (ok bool) {
	ok = true
	defer func() {
		if err := recover(); err != nil {
			ok = false
		}
	}()

	_ = AsTime(a)
	return
}

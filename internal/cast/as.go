package cast

import (
	"encoding/json"
	"fmt"
	reflect2 "github.com/expectto/be/internal/reflect"
	"reflect"
	"time"
)

// AsString converts the given input into a string or a string-like representation.
// It supports various input types, including actual strings, byte slices, JSON RawMessage, and custom string types.
//
// Input values may also be pointers.
//
// Note: If the input is []byte and contains a non-UTF-8 valid sequence, the resulting string may be invalid.
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	str = AsString("example") // Converts a string, returns "example"
//	str = AsString([]byte("byte_data")) // Converts a byte slice, returns "byte_data"
//	str = AsString(CustomStringType("example")) // Converts a custom string type
//
// This function is useful for converting diverse input types into a string representation,
// and it is designed to provide a convenient string conversion for various testing scenarios.
func AsString(a any) string {
	// First start with a type casting
	switch t := a.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	case json.RawMessage:
		return string(t)
	case *json.RawMessage: // shortcut without reflect
		return string(*t)

		// we intentionally do not support fmt.Stringer here
		// it must be handled manually
	}

	// Then fallback to reflect, in case we have custom string/[]byte types
	v := reflect.ValueOf(a)
	v = reflect2.IndirectDeep(v)

	if v.Kind() == reflect.String {
		return v.String()
	}
	if v.Kind() == reflect.Slice && v.Type().AssignableTo(reflect.TypeOf([]byte{})) {
		return string(v.Bytes())
	}

	panic(fmt.Sprintf("Expected a string-ish/[]byte-ish thing! Got <%T>", a))
}

// AsBytes converts the given input into a []byte or a []byte-like representation.
// It supports various input types, including byte slices, JSON RawMessage, and strings.
//
// Input values may also be pointers.
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	bytes = AsBytes([]byte("byte_data")) // Converts a byte slice, returns []byte with the same content
//	bytes = AsBytes(jsonRawMessage) // Converts a JSON RawMessage
//	bytes = AsBytes("example") // Converts a string, returns the corresponding []byte
//
// This function is useful for converting diverse input types into a []byte representation,
// and it is designed to provide a convenient []byte conversion for various testing scenarios.
func AsBytes(a any) []byte {
	// First start with a type casting
	switch t := a.(type) {
	case []byte:
		return t
	case json.RawMessage:
		return t
	case *json.RawMessage: // shortcut without reflect
		return *t
	case string:
		return []byte(t)
	}

	// Then fallback to reflect, in case we have custom string/[]byte types
	v := reflect.ValueOf(a)
	v = reflect2.IndirectDeep(v)

	if v.Kind() == reflect.Slice && v.Type().AssignableTo(reflect.TypeOf([]byte{})) {
		return v.Bytes()
	}
	if v.Kind() == reflect.String {
		return []byte(v.String())
	}

	panic(fmt.Sprintf("Expected []byte-ish/string-ish thing! Got <%T>", a))
}

// AsBool converts the given input into a bool.
// It supports various input types, including actual bool values and pointers to bool.
// Input values may also be pointers.
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	value = AsBool(true) // Converts a bool, returns true
//	value = AsBool(&boolValue) // Converts a pointer to a bool, returns the bool value
//
// This function is designed for converting different input types into bool values,
// and it is useful for various testing scenarios where boolean values are expected.
func AsBool(a any) bool {
	// First start with a type casting
	switch t := a.(type) {
	case bool:
		return t
	case *bool:
		return *t
	}

	// fallback to reflect
	v := reflect.ValueOf(a)
	v = reflect2.IndirectDeep(v)

	if v.Kind() == reflect.Bool {
		return v.Bool()
	}

	panic(fmt.Sprintf("Expected a bool! Got <%T>: %#v", a, a))
}

// AsInt converts the given input into an int.
// It supports various input types, including int values, float64 values (for integral floats),
// and pointers to int or float64.
// Input values may also be pointers.
//
// Note (1): Float64 values are converted to int only if they are integral floats (e.g., 42.0). Otherwise, use AsFloat.
// Note (2): Depending on the machine where the code is compiled, the resulting int may be of different sizes (e.g., int32).
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	intValue := AsInt(42) // Converts an int, returns 42
//	intValue := AsInt(&intValuePtr) // Converts a pointer to int, returns the int value
//	intValue := AsInt(42.0) // Converts an integral float, returns 42
//
// This function is designed for converting different input types into int values,
// and it is useful for various testing scenarios where integer values are expected.
func AsInt(a any) int {
	// First start with a type casting
	switch t := a.(type) {
	case int:
		return t
	case *int:
		return *t
	case int8:
		return int(t)
	case *int8:
		return int(*t)
	case int16:
		return int(t)
	case *int16:
		return int(*t)
	case int32:
		return int(t)
	case *int32:
		return int(*t)
	case int64:
		return int(t)
	case *int64:
		return int(*t)
	case uint:
		return int(t)
	case *uint:
		return int(*t)
	case uint8:
		return int(t)
	case *uint8:
		return int(*t)
	case uint16:
		return int(t)
	case *uint16:
		return int(*t)
	case uint32:
		return int(t)
	case *uint32:
		return int(*t)
	case uint64:
		return int(t)
	case *uint64:
		return int(*t)
	case float64:
		intResult := int(t)
		if float64(intResult) != t {
			panic("Expected an integral float")
		}
		return intResult
	case *float64:
		intResult := int(*t)
		if float64(intResult) != *t {
			panic("Expected an integral float")
		}
		return intResult
	case float32:
		intResult := int(t)
		if float32(intResult) != t {
			panic("Expected an integral float")
		}
		return intResult
	case *float32:
		intResult := int(*t)
		if float32(intResult) != *t {
			panic("Expected an integral float")
		}
		return intResult
	}

	// fallback to reflect
	v := reflect.ValueOf(a)
	v = reflect2.IndirectDeep(v)

	if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
		return int(v.Int())
	} else if v.Kind() >= reflect.Uint && v.Kind() <= reflect.Uint64 {
		return int(v.Uint())
	} else if v.Kind() >= reflect.Float32 && v.Kind() <= reflect.Float64 {
		intResult := int(v.Float())
		if float64(intResult) != v.Float() {
			panic("Expected an integer float")
		}
		return intResult
	}

	panic(fmt.Sprintf("Expected an integer number! Got <%T>: %#v", a, a))
}

// AsFloat converts the given input into a float64.
// It supports various input types, including float64 values and int values (converted to float64).
// Input values may also be pointers.
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	floatValue = AsFloat(3.14) // Converts a float64, returns 3.14
//	floatValue = AsFloat(&floatValuePtr) // Converts a pointer to float64, returns the float64 value
//	floatValue = AsFloat(42) // Converts an int to a float64, returns 42.0
//
// This function is designed for converting different input types into float64 values,
// and it is useful for various testing scenarios where floating-point values are expected.
func AsFloat(a any) float64 {
	// First start with a type casting
	switch t := a.(type) {
	case float64:
		return t
	case *float64:
		return *t
	case float32:
		return float64(t)
	case *float32:
		return float64(*t)
	case int:
		return float64(t)
	case *int:
		return float64(*t)
	case int8:
		return float64(t)
	case *int8:
		return float64(*t)
	case int16:
		return float64(t)
	case *int16:
		return float64(*t)
	case int32:
		return float64(t)
	case *int32:
		return float64(*t)
	case int64:
		return float64(t)
	case *int64:
		return float64(*t)
	case uint:
		return float64(t)
	case *uint:
		return float64(*t)
	case uint8:
		return float64(t)
	case *uint8:
		return float64(*t)
	case uint16:
		return float64(t)
	case *uint16:
		return float64(*t)
	case uint32:
		return float64(t)
	case *uint32:
		return float64(*t)
	case uint64:
		return float64(t)
	case *uint64:
		return float64(*t)
	}

	// fallback to reflect
	v := reflect.ValueOf(a)
	v = reflect2.IndirectDeep(v)

	if v.Kind() >= reflect.Float32 && v.Kind() <= reflect.Float64 {
		return v.Float()
	} else if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
		return float64(v.Int())
	} else if v.Kind() >= reflect.Uint && v.Kind() <= reflect.Uint64 {
		return float64(v.Uint())
	}

	panic(fmt.Sprintf("Expected a float number! Got <%T>: %#v", a, a))
}

// AsKind converts the given input into a reflect.Kind.
// It supports various input types, including reflect.Kind values and pointers to reflect.Kind.
// Input values may also be pointers.
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	kind = AsKind(reflect.Int) // Converts a reflect.Kind, returns reflect.Int
//	kind = AsKind(&kindPtr) // Converts a pointer to reflect.Kind, returns the reflect.Kind value
//
// This function is designed for converting different input types into reflect.Kind values,
// and it is useful for various testing scenarios where reflection is used.
func AsKind(a any) reflect.Kind {
	// First start with a type casting
	switch t := a.(type) {
	case reflect.Kind:
		return t
	case *reflect.Kind:
		return *t
	}

	// No reason to fallback to reflect here
	// as too small chance that it will be a custom reflect.Kind type
	// or a deeper pointer

	panic(fmt.Sprintf("Expected a reflect.Kind!  Got <%T>: %#v", a, a))
}

// AsSliceOfAny converts the given input into a []any.
// It supports various input types, including slices, arrays, and pointers to slices or arrays of any type.
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	anySlice := AsSliceOfAny([]string{"a", "b", "c"}) // Converts a string slice, returns []any{"a", "b", "c"}
//	anySlice := AsSliceOfAny(&intArrayPtr) // Converts a pointer to an array, returns []any with the array content
//
// This function is designed for converting different input types into a []any,
// and it is useful for various testing scenarios where a slice of arbitrary types is expected.
func AsSliceOfAny(v any) []any {
	// First start with a type casting
	switch t := v.(type) {
	case []any:
		return t
	}

	// Then fallback to reflect
	rv := reflect.ValueOf(v)
	rv = reflect2.IndirectDeep(rv)

	if rv.Kind() == reflect.Slice {
		slice := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			slice[i] = rv.Index(i).Interface()
		}
		return slice
	}
	// todo: support arrays

	panic(fmt.Sprintf("Expected a slice! Got <%T>: %#v", v, v))
}

func AsStrings(v any) []string {
	// First start with a type casting
	switch t := v.(type) {
	case []string:
		return t
	}

	// Then fallback to reflect
	rv := reflect.ValueOf(v)
	rv = reflect2.IndirectDeep(rv)
	st := reflect.TypeOf("")

	// todo: support arrays?
	if rv.Kind() == reflect.Slice {
		slice := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			if !rv.Index(i).Type().ConvertibleTo(st) {
				panic(fmt.Sprintf("expected a slice string! Got s[%d] <%T>: %#v", i, v, v))
			}

			slice[i] = rv.Index(i).String()
		}
		return slice
	}

	panic(fmt.Sprintf("Expected a slice! Got <%T>: %#v", v, v))
}

// AsTime converts the given input into a time.Time.
// It supports various input types, including time.Time values and pointers to time.Time.
//
// It panics if it's not possible to perform the conversion.
//
// Example Usage:
//
//	timestamp = AsTime(time.Now()) // Converts a time.Time, returns the current time
//	timestamp = AsTime(&timeValuePtr) // Converts a pointer to time.Time, returns the time.Time value
//
// This function is designed for converting different input types into time.Time values.
func AsTime(a any) time.Time {

	// First start with a type casting
	switch t := a.(type) {
	case time.Time:
		return t
	case *time.Time:
		return *t
	}

	// fallback to reflect
	v := reflect.ValueOf(a)
	v = reflect2.IndirectDeep(v)

	if v.CanConvert(reflect2.TypeFor[time.Time]()) {
		return v.Interface().(time.Time)
	}

	panic(fmt.Sprintf("Expected a time.Time value! Got <%T>: %#v", a, a))
}

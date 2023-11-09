package cast

import (
	"encoding/json"
	"fmt"
	reflect2 "github.com/expectto/be/internal/reflect"
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

// IsString checks if the given input is a string or string-like.
// To avoid duplicating type-checking logic, it provides extensive configuration options for
// customizing the type-checking behavior, making it a versatile utility for testing code.
// It supports both strict and non-strict mode checks, allowing you to precisely control
// which types are considered string-like. It also provides options for handling custom types,
// pointer dereferencing, and stringer interfaces.
//
// Example Usage:
//
//		// In a non-strict check, allows custom types, pointer dereferencing, and stringer interfaces
//		IsString("example", AllowCustomTypes(), AllowPointers(), AllowStringer()) // returns true
//
//		// In a strict check, only actual strings are accepted
//		isStringStrict := IsString(Strict())
//		IsString("example", Strict()) // Returns true
//	 IsString([]byte("example"), Strict()) // Returns false
func IsString(a any, opts ...optIsString) bool {
	// Even before computing the config,
	// if input is simply a string, return immediately
	_, ok := a.(string)
	if ok {
		return ok
	}

	// building a default config and override it with users options
	cfg := defaultConfigIsString.clone()
	for _, opt := range opts {
		opt(cfg)
	}

	// if it was a strict check, and simple casting failed, we can't continue
	if cfg.Strict && !ok {
		return false
	}

	// in non-strict mode we allow different string-ish types ([]byte, custom strings, fmt.Stringer().
	// Also we allow value to be hidden under pointer deeply
	// dedicated options can configure how non-strict mode works
	if cfg.AllowPointers && cfg.AllowDeepPointers && cfg.AllowCustomTypes && cfg.AllowStringer && cfg.AllowBytesConversion {
		return IsStringish(a)
	}

	if cfg.AllowBytesConversion {
		// First start with a type casting
		switch a.(type) {
		case []byte, json.RawMessage:
			return true
		}

		if cfg.AllowPointers {
			switch a.(type) {
			case *[]byte, *json.RawMessage:
				return true
			}
		}
	}

	if cfg.AllowStringer {
		if _, ok := a.(fmt.Stringer); ok {
			return true
		}
	}

	if cfg.AllowCustomTypes {
		v := reflect.ValueOf(a)
		if cfg.AllowPointers {
			if cfg.AllowDeepPointers {
				v = reflect2.IndirectDeep(v)
			} else {
				v = reflect.Indirect(v)
			}
		}

		if v.Kind() == reflect.String {
			return true
		}
		if v.Kind() == reflect.Slice && v.Type().AssignableTo(reflect.TypeOf([]byte{})) {
			return true
		}
	}

	return false
}

type configIsString struct {
	Strict bool

	// Non-strict mode:
	AllowCustomTypes     bool
	AllowBytesConversion bool
	AllowPointers        bool
	AllowDeepPointers    bool
	AllowStringer        bool
}

var defaultConfigIsString *configIsString

// clone is done via json round-trip marshalling
func (cis *configIsString) clone() *configIsString {
	// errors are omitted intentionally, as here we can't fail
	contents, _ := json.Marshal(cis)
	var clone configIsString
	_ = json.Unmarshal(contents, &clone)
	return &clone
}

func init() {
	// By default, it's not strict and allows all conversion options.
	SetDefaultIsStringConfig(
		AllowStringer(),
		AllowBytesConversion(),
		AllowCustomTypes(),
		AllowPointers(),
		AllowDeepPointers(),
	)
}

// SetDefaultIsStringConfig sets the default configuration for IsString checks.
func SetDefaultIsStringConfig(opts ...optIsString) {
	cfg := &configIsString{}
	for _, opt := range opts {
		opt(cfg)
	}
	if defaultConfigIsString == nil {
		defaultConfigIsString = &configIsString{}
	}
	*defaultConfigIsString = *cfg
}

type optIsString func(config *configIsString)

// Strict option enables strict mode. All other options are automatically ignored.
func Strict() optIsString {
	return func(cfg *configIsString) {
		cfg.Strict = true
	}
}

// AllowStringer option allows the use of fmt.Stringer for IsString checks.
func AllowStringer() optIsString {
	return func(cfg *configIsString) {
		cfg.AllowStringer = true
	}
}

// AllowCustomTypes option allows the use of custom string types for IsString checks.
func AllowCustomTypes() optIsString {
	return func(cfg *configIsString) {
		cfg.AllowCustomTypes = true
	}
}

// AllowBytesConversion option allows conversion from []byte to string for IsString checks.
func AllowBytesConversion() optIsString {
	return func(cfg *configIsString) {
		cfg.AllowBytesConversion = true
	}
}

// AllowPointers option allows checking of values under pointers for IsString checks.
func AllowPointers() optIsString {
	return func(cfg *configIsString) {
		cfg.AllowPointers = true
	}
}

// AllowDeepPointers option allows deep checking of values under pointers for IsString checks.
func AllowDeepPointers() optIsString {
	return func(cfg *configIsString) {
		cfg.AllowDeepPointers = true
	}
}

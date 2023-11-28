package cast

import (
	"encoding/json"
	reflect2 "github.com/expectto/be/internal/reflect"
	"reflect"
)

// IsString checks if the given input is a string or string-like.
// To avoid duplicating type-checking logic, it provides extensive configuration options for
// customizing the type-checking behavior, making it a versatile utility for testing code.
// It supports both strict and non-strict mode checks, allowing you to precisely control
// which types are considered string-like. It also provides options for handling custom types,
// pointer dereferencing..
//
// Example Usage:
//
//	// In a non-strict check, allows custom types, pointer dereferencing
//	IsString("example", AllowCustomTypes(), AllowPointers())) // returns true
//
//	// In a strict check, only actual strings are accepted
//	isStringStrict := IsString(Strict())
//	IsString("example", Strict()) // Returns true
//	IsString([]byte("example"), Strict()) // Returns false
func IsString(a any, opts ...optIsString) bool {
	// Even before computing the config,
	// if input is simply a string, return immediately
	_, ok := a.(string)
	if ok {
		return ok
	}

	// building a default config and override it with users options
	cfg := defaultIsStringConfig.clone()
	for _, opt := range opts {
		opt(cfg)
	}

	// if it was a strict check, and simple casting failed, we can't continue
	if cfg.IsStrict() && !ok {
		return false
	}

	// in allow-all mode we can simply call IsStringish
	if cfg.AllowsAll() {
		return IsStringish(a)
	}

	// We can still use type casting for simple cases, like AllowBytesConversion, AllowPointer:

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

	// Further we can only try reflect

	v := reflect.ValueOf(a)
	if cfg.AllowDeepPointers {
		v = reflect2.IndirectDeep(v)
	} else if cfg.AllowPointers {
		v = reflect.Indirect(v)
	}

	if v.Type() == reflect2.TypeFor[string]() {
		return true
	}

	if cfg.AllowCustomTypes {
		if v.Kind() == reflect.String {
			return true
		}

		if cfg.AllowBytesConversion {
			if v.Kind() == reflect.Slice && v.Type().AssignableTo(reflect.TypeOf([]byte{})) {
				return true
			}
		}
	}

	return false
}

// isStringConfig is a configuration for IsString check.
// An empty config (all flags=false) is considered "strict mode"
// `omitempty` tag is needed for marshalling for IsStrict() func
type isStringConfig struct {
	AllowCustomTypes     bool `json:"allow_custom_types,omitempty"`
	AllowBytesConversion bool `json:"allow_bytes_conversion,omitempty"`
	AllowPointers        bool `json:"allow_pointers,omitempty"`
	AllowDeepPointers    bool `json:"allow_deep_pointers,omitempty"`
}

var defaultIsStringConfig *isStringConfig

func init() {
	// no options given will lead to strict mode by default
	ConfigureIsStringConfig()
}

// clone is done via json round-trip marshalling
func (cis *isStringConfig) clone() *isStringConfig {
	// errors are omitted intentionally, as here we can't fail
	contents, _ := json.Marshal(cis)
	var clone isStringConfig
	_ = json.Unmarshal(contents, &clone)
	return &clone
}

// IsStrict returns true if all custom options are disabled
// IsString() in strict mode will return true only for actual `string` values
func (cis *isStringConfig) IsStrict() bool {
	// Strict mode is when all flags are false
	marshalled, _ := json.Marshal(cis)
	return string(marshalled) == "{}"
}

// AllowsAll returns true if all custom options are enabled
func (cis *isStringConfig) AllowsAll() bool {
	el := reflect.ValueOf(cis).Elem()
	var result = true
	for i := 0; i < el.NumField(); i++ {
		result = result && el.Field(i).Bool()

		if !result {
			break
		}
	}
	return result
}

// ConfigureIsStringConfig sets the default configuration for IsString checks.
func ConfigureIsStringConfig(opts ...optIsString) {
	cfg := &isStringConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if defaultIsStringConfig == nil {
		defaultIsStringConfig = &isStringConfig{}
	}
	*defaultIsStringConfig = *cfg
}

type optIsString func(config *isStringConfig)

// AllowCustomTypes option allows the use of custom string types for IsString checks.
func AllowCustomTypes() optIsString {
	return func(cfg *isStringConfig) { cfg.AllowCustomTypes = true }
}

// AllowBytesConversion option allows conversion from []byte to string for IsString checks.
func AllowBytesConversion() optIsString {
	return func(cfg *isStringConfig) { cfg.AllowBytesConversion = true }
}

// AllowPointers option allows checking of values under pointers for IsString checks.
func AllowPointers() optIsString {
	return func(cfg *isStringConfig) { cfg.AllowPointers = true }
}

// AllowDeepPointers option allows deep checking of values under pointers for IsString checks.
func AllowDeepPointers() optIsString {
	return func(cfg *isStringConfig) { cfg.AllowDeepPointers = true }
}

// AllowAll option allows all options (makes it the most non-strict)
func AllowAll() optIsString {
	return func(cfg *isStringConfig) {
		v := reflect.ValueOf(cfg).Elem()

		for i := 0; i < v.NumField(); i++ {
			v.Field(i).SetBool(true)
		}
	}
}

package be_json

import (
	"encoding/json"
	"fmt"
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/internal/cast"
	. "github.com/expectto/be/internal/psi"
	"github.com/expectto/be/internal/psi_matchers"
	"github.com/expectto/be/types"
	"github.com/onsi/gomega"
	"io"
)

type JsonInputType uint32

const (
	JsonAsBytes JsonInputType = 1 << iota
	JsonAsString
	JsonAsStringer
	JsonAsReader
	JsonAsObject
	JsonAsObjects
)

// Matcher is a JSON matcher. "JSON" here means a []byte with JSON data in it
// By default several input types are available: string(*) / []byte(*), fmt.Stringer, io.Reader
//   - custom string-based or []byte-based types are available as well
//
// To make it stricter and to specify which format JSON we should expect, you
// must pass one of transforms as first argument:
//   - JsonAsBytes/ JsonAsString / JsonAsStringer  / JsonAsReader (for string-like representation)
//   - JsonAsObject / JsonAsObjects (for map[string]any representation)
func Matcher(args ...any) types.BeMatcher {
	// Default input is ok to be any of these
	var inputMatcher types.BeMatcher = psi_matchers.NewAnyMatcher(
		// String-like inputs:
		be_reflected.AsBytes(), be_reflected.AsString(), be_reflected.AsStringer(), be_reflected.AsReader(),

		// Object-like inputs:
		// Here we accept map[string]any or []map[string]any
		be_reflected.AsObject(), be_reflected.AsObjects(),
	)

	// Check if first argument was given as a JsonAs* constant
	// that needs to be handled
	if len(args) > 0 {
		if t, ok := args[0].(JsonInputType); ok {
			inputMatchers := make([]types.BeMatcher, 0)
			if t&JsonAsBytes != 0 {
				inputMatchers = append(inputMatchers, be_reflected.AsBytes())
			}
			if t&JsonAsString != 0 {
				inputMatchers = append(inputMatchers, be_reflected.AsString())
			}
			if t&JsonAsStringer != 0 {
				inputMatchers = append(inputMatchers, be_reflected.AsStringer())
			}
			if t&JsonAsReader != 0 {
				inputMatchers = append(inputMatchers, be_reflected.AsReader())
			}
			if t&JsonAsObject != 0 {
				inputMatchers = append(inputMatchers, be_reflected.AsObject())
			}
			if t&JsonAsObjects != 0 {
				inputMatchers = append(inputMatchers, be_reflected.AsObjects())
			}

			// To avoid extra "Any" matching logic, let's simplify case when we have single input matcher
			if len(inputMatchers) == 1 {
				inputMatcher = inputMatchers[0]
			} else {
				inputMatcher = psi_matchers.NewAnyMatcher(cast.AsSliceOfAny(inputMatchers)...)
			}
			args = args[1:]
		}
	}

	// If no args (after handling JsonAs* constants)
	// then we just match if it's valid json
	if len(args) == 0 {
		// TODO: it must validate if it's actualy a valid JSON
		return inputMatcher
	}

	return &psi_matchers.AllMatcher{Matchers: []types.BeMatcher{
		inputMatcher,

		// JSON expects arguments to be matchers upon map[string]any
		// So let's perform a transform: raw => any
		WithFallibleTransform(func(actual any) any {
			// `actual` may be an io.Reader that is decoded directly
			if reader, ok := actual.(io.Reader); ok {
				var data any
				if err := json.NewDecoder(reader).Decode(&data); err != nil {
					return NewTransformError(fmt.Errorf("to read json: %w", err), actual)
				}
				if closer, ok := actual.(io.Closer); ok {
					_ = closer.Close()
				}

				return data
			}

			// convert `actual` into `any` (if `actual` is bytes/string):
			// it will end up `[]any` or `map[string]any` underneath it
			if cast.IsStringish(actual) {
				var data any
				if err := json.Unmarshal(cast.AsBytes(actual), &data); err != nil {
					return NewTransformError(fmt.Errorf("be a valid json: %w", err), actual)
				}

				return data
			}

			// no conversion is needed, `actual` will be checked via matchers directly
			return actual
		},

			// Applying given matchers to the raw JSON
			func() types.BeMatcher {
				// If we have just one arg then we match against it
				// If it's a string, we're remarshalling it into object
				if len(args) == 1 && !IsMatcher(args[0]) {
					var argData any
					if cast.IsStringish(args[0]) {
						if err := json.Unmarshal(cast.AsBytes(args[0]), &argData); err != nil {
							return psi_matchers.NewNeverMatcher(err)
						}
					} else {
						// todo: check if it's actually object|objects
						//       so we can nicer failure messages
						argData = args[0]
					}

					return psi_matchers.NewEqMatcher(argData)
				}

				return Psi(args...)
			}(),
		),
	}}
}

// HaveKeyValue is a facade to gomega.HaveKey & gomega.HaveKeyWithValue
func HaveKeyValue(key any, args ...any) types.BeMatcher {
	if len(args) == 0 {
		return Psi(gomega.HaveKey(key))
	}

	// todo: optimize for gomock messages ?
	// todo: should we optimize (check if len(args)==1,
	// 		 to reduce level of wrapping) ?
	return Psi(
		gomega.HaveKeyWithValue(key, Psi(args...)),
	)
}

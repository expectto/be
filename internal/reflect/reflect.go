// Package reflect contains helpers that extends standard reflect library
package reflect

import "reflect"

// IndirectDeep does reflect.Indirect deeply
func IndirectDeep(v reflect.Value) reflect.Value {
	for {
		if v.Kind() != reflect.Pointer {
			break
		}
		v = v.Elem()
	}
	return v
}

// LengthOf returns length of a given type
func LengthOf(a any) (int, bool) {
	if a == nil {
		return 0, false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Map, reflect.Array, reflect.String, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(a).Len(), true
	default:
		return 0, false
	}
}

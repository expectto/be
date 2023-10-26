// Package reflect contains helpers that extends standard reflect library
package reflect

import "reflect"

// TypeFor returns a reflect.Type for a given type
// Deprecated: Should go away when _I hope_ it will be implemented in reflect(go-1.22)
func TypeFor[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

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

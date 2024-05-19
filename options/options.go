// Package options declares options to be used in customizeable matchers
// Note: Options of ALL `be_*` matchers are stored here, in a separate package `options`.
//
//	 It's done with consideration that `options` package will be imported via `dot import`
//	 so inside your tests options are clear & short, without package name
//	e.g.: `Expect(myString).To(be_string.Only( Alpha | Numeric | Whitespace )`
package options

// Option represents an option for any customizable matcher
type Option int

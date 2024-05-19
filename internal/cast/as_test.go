package cast_test

import (
	"encoding/json"
	"github.com/expectto/be/internal/cast"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("As", func() {
	Context("AsString", func() {
		It("should return string for string", func() {
			Expect(cast.AsString("something")).To(Equal("something"))
		})

		It("should return string for empty string", func() {
			Expect(cast.AsString("")).To(Equal(""))
		})

		It("should return string for []byte", func() {
			Expect(cast.AsString([]byte("foobar"))).To(Equal("foobar"))
		})

		It("should return string for empty []byte", func() {
			Expect(cast.AsString([]byte{})).To(Equal(""))
		})

		It("should return string for CustomString", func() {
			type CustomString string
			Expect(cast.AsString(CustomString("foobar"))).To(Equal("foobar"))
		})

		It("should return string for json.RawMessage", func() {
			Expect(cast.AsString(json.RawMessage(`{"foo":"bar"}`))).To(Equal(`{"foo":"bar"}`))
		})

		It("should return string for *json.RawMessage", func() {
			msg := json.RawMessage(`{"foo":"bar"}`)
			Expect(cast.AsString(&msg)).To(Equal(`{"foo":"bar"}`))
		})

		It("should return string for a string under the pointer", func() {
			Expect(cast.AsString(new(string))).To(Equal(""))
		})

		It("should panic for non-stringish", func() {
			Expect(func() { cast.AsString(123) }).To(Panic())
		})
	})
})

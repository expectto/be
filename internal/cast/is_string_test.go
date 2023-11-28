package cast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/expectto/be/internal/cast"
)

var _ = Describe("IsString", func() {
	When("in strict mode", func() {
		It("should return true for string", func() {
			Expect(cast.IsString("something")).To(BeTrue())
		})
		It("should return true for empty string", func() {
			Expect(cast.IsString("")).To(BeTrue())
		})
		It("should return false for []byte", func() {
			Expect(cast.IsString([]byte("foobar"))).To(BeFalse())
		})
		It("should return false for empty []byte", func() {
			Expect(cast.IsString([]byte{})).To(BeFalse())
		})
		It("should return false for *string", func() {
			helloWorld := "hello world"
			Expect(cast.IsString(&helloWorld)).To(BeFalse())
		})
	})

	When("allowing all", func() {
		type CustomString string
		type CustomBytes []byte
		It("should return true for custom strings", func() {
			Expect(cast.IsString(CustomString("hello world"), cast.AllowAll())).To(BeTrue())
		})
		It("should return true for custom strings under pointer", func() {
			helloWorld := CustomString("hello world")
			ptrHelloWorld := &helloWorld
			Expect(cast.IsString(&helloWorld, cast.AllowAll())).To(BeTrue())
			Expect(cast.IsString(&ptrHelloWorld, cast.AllowAll())).To(BeTrue())
		})
		It("should return true for bytes", func() {
			Expect(cast.IsString([]byte("hello world"), cast.AllowAll())).To(BeTrue())
		})
		It("should return true for custom bytes", func() {
			helloWorld := CustomBytes("hello world")
			Expect(cast.IsString(helloWorld, cast.AllowAll())).To(BeTrue())
			Expect(cast.IsString(&helloWorld, cast.AllowAll())).To(BeTrue())

			ptrHelloWorld := &helloWorld
			Expect(cast.IsString(&ptrHelloWorld, cast.AllowAll())).To(BeTrue())
		})
	})

	When("allowing pointers", func() {
		It("should return true for string under the pointer", func() {
			Expect(cast.IsString(new(string), cast.AllowPointers())).To(BeTrue())
			Expect(cast.IsString(new(string), cast.AllowDeepPointers())).To(BeTrue())

			s := "hello"
			Expect(cast.IsString(&s, cast.AllowPointers())).To(BeTrue())
			Expect(cast.IsString(&s, cast.AllowDeepPointers())).To(BeTrue())
			ss := &s
			Expect(cast.IsString(&ss, cast.AllowDeepPointers())).To(BeTrue())
			Expect(cast.IsString(&ss, cast.AllowPointers())).To(BeFalse())
		})

		It("should return false for not-a-string under the pointer", func() {
			Expect(cast.IsString(new(int), cast.AllowPointers())).To(BeFalse())
		})
	})

})

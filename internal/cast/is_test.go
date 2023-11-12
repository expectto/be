package cast_test

import (
	"github.com/expectto/be/internal/cast"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type CustomStringer struct{ v string }

func (s CustomStringer) String() string {
	return s.v
}

var _ = Describe("Is", func() {
	Context("IsNil", func() {
		It("should return true for nil", func() {
			Expect(cast.IsNil(nil)).To(BeTrue())
		})
		It("should return true for typed nil", func() {
			var i *int
			Expect(cast.IsNil(i)).To(BeTrue())
		})
		It("should return true for interface nil", func() {
			var i interface{}
			Expect(cast.IsNil(i)).To(BeTrue())
		})

		It("should return false for non-nil pointer", func() {
			Expect(cast.IsNil(&struct{}{})).To(BeFalse())
		})
		It("should return false for non-nil map", func() {
			Expect(cast.IsNil(map[string]int{})).To(BeFalse())
		})
		It("should return false for non-nil func", func() {
			Expect(cast.IsNil(func() {})).To(BeFalse())
		})

		It("should return false for non-nil digit", func() {
			Expect(cast.IsNil(0)).To(BeFalse())
		})
		It("should return false for non-nil string", func() {
			Expect(cast.IsNil("")).To(BeFalse())
		})
	})

	Context("IsStringish", func() {
		When("considered stringish", func() {
			It("should return true for string", func() {
				Expect(cast.IsStringish("something")).To(BeTrue())
			})
			It("should return true for empty string", func() {
				Expect(cast.IsStringish("")).To(BeTrue())
			})
			It("should return true for []byte", func() {
				Expect(cast.IsStringish([]byte("foobar"))).To(BeTrue())
			})
			It("should return true for empty []byte", func() {
				Expect(cast.IsStringish([]byte{})).To(BeTrue())
			})
			It("should return true for stringer", func() {
				Expect(cast.IsStringish(CustomStringer{v: "hello-world"})).To(BeTrue())
			})
		})

		When("considered not stringish", func() {
			It("should return false for nil", func() {
				Expect(cast.IsStringish(nil)).To(BeFalse())
			})
			It("should return false for int", func() {
				Expect(cast.IsStringish(123)).To(BeFalse())
			})
			It("should return false for float", func() {
				Expect(cast.IsStringish(123.456)).To(BeFalse())
			})
			It("should return false for bool", func() {
				Expect(cast.IsStringish(true)).To(BeFalse())
			})
			It("should return false for complex", func() {
				Expect(cast.IsStringish(1 + 2i)).To(BeFalse())
			})
			It("should return false for struct", func() {
				Expect(cast.IsStringish(struct{}{})).To(BeFalse())
			})
			It("should return false for map", func() {
				Expect(cast.IsStringish(map[string]int{})).To(BeFalse())
			})
			It("should return false for func", func() {
				Expect(cast.IsStringish(func() {})).To(BeFalse())
			})
		})
	})

	Context("IsString", func() {
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
			It("should return false for stringer", func() {
				Expect(cast.IsString(CustomStringer{v: "hello-world"})).To(BeFalse())
			})
		})

		When("allowing pointers", func() {
			It("should return true for string under the pointer", func() {
				Expect(cast.IsString(new(string), cast.AllowPointers())).To(BeTrue())
				Expect(cast.IsString(new(string), cast.AllowDeepPointers())).To(BeTrue())
			})

			It("should return false for not-a-string under the pointer", func() {
				Expect(cast.IsString(new(int), cast.AllowPointers())).To(BeFalse())
			})

			It("should return true for string under the pointer", func() {
				x := new(string)
				*x = "hello-world"
				y := &x

				Expect(cast.IsString(y, cast.AllowDeepPointers())).To(BeTrue())
				Expect(cast.IsString(y, cast.AllowPointers())).To(BeFalse())
			})
		})

	})
})

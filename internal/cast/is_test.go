package cast_test

import (
	"github.com/expectto/be/internal/cast"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

	Context("IsInt", func() {
		It("should return true for int", func() {
			Expect(cast.IsInt(123)).To(BeTrue())
		})
		It("should return true for int8", func() {
			Expect(cast.IsInt(int8(123))).To(BeTrue())
		})
		It("should return true for int16", func() {
			Expect(cast.IsInt(int16(123))).To(BeTrue())
		})
		It("should return true for int32", func() {
			Expect(cast.IsInt(int32(123))).To(BeTrue())
		})
		It("should return true for int64", func() {
			Expect(cast.IsInt(int64(123))).To(BeTrue())
		})
		It("should return false for uint", func() {
			Expect(cast.IsInt(uint(123))).To(BeTrue())
		})

		It("should return false for float", func() {
			Expect(cast.IsInt(123.456)).To(BeFalse())
		})
		It("should return false for string", func() {
			Expect(cast.IsInt("123")).To(BeFalse())
		})
		It("should return false for bool", func() {
			Expect(cast.IsInt(true)).To(BeFalse())
		})
		It("should return false for nil", func() {
			Expect(cast.IsInt(nil)).To(BeFalse())
		})
		It("should return false for struct", func() {
			Expect(cast.IsInt(struct{}{})).To(BeFalse())
		})
	})
})

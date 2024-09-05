package be_struct_test

import (
	"github.com/expectto/be/be_struct"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeStructs", func() {
	It("should match a struct field", func() {
		type TestStruct struct {
			Field1 string
			Field2 int
		}

		var result = TestStruct{
			Field1: "hello1",
		}

		Expect(result).To(
			And(
				be_struct.HavingField[TestStruct]("Field1", "hello1"),
				be_struct.HavingField[TestStruct]("Field2"), // just ensure it exists
			))

		Expect(result).NotTo(
			be_struct.HavingField[TestStruct]("Field3"),
		)

		type WrongStruct struct {
			Field1 string
		}
		Expect(result).NotTo(
			be_struct.HavingField[WrongStruct]("Field1"),
		)
	})
})

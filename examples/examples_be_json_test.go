package examples_test

import (
	"encoding/json"

	"github.com/expectto/be"
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_reflected"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Foobar struct {
	Foo string `json:"foo"`
}

func (f Foobar) String() string {
	jj, _ := json.Marshal(f)
	return string(jj)
}

var _ = Describe("Examples on matching JSON", func() {
	Context("Valid JSON string", func() {

		It("should match against valid JSON string", func() {
			Expect(`{"foo":"bar"}`).To(be.JSON(
				be_json.JsonAsString,
				be_reflected.AsObject(),
				be_json.HaveKeyValue("foo", "bar"),
			))
		})

		// TODO: fixme
		PIt("should match against valid JSON struct", func() {
			f := Foobar{Foo: "bar"}
			Expect(f).To(be.JSON(
				be_json.JsonAsStringer,
				be_reflected.AsObject(),
				be_json.HaveKeyValue("foo", "bar"),
			))
		})
	})
})

package be_url_test

import (
	"github.com/expectto/be/be_json"
	"github.com/expectto/be/be_reflected"
	"github.com/expectto/be/be_url"
	"github.com/expectto/be/internal/testing/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/url"
)

var _ = Describe("URL", func() {
	Context("Function that returns a *url.URL", func() {

		It("should match against all parts of *url.URL", func() {
			u, err := url.Parse(`https://example.com/path/to/?foo=bar_123&v=1&payload={"hello":"world"}`)
			Expect(err).Should(Succeed())

			Expect(u).To(be_url.URL(
				be_url.WithHttps(),
				be_url.NotHavingPort(),
				be_url.HavingHostname("example.com"),
				be_url.HavingPath("/path/to/"),
				be_url.HavingRawQuery(HavePrefix("foo=bar_123&v=1")),
				be_url.HavingSearchParam("foo", ContainSubstring("bar_")),
				be_url.HavingSearchParam("v", be_reflected.AsNumericString()),
				be_url.HavingSearchParam("payload", be_json.Matcher()),
			))
		})

		It("should simply match separate url args from *url.URL", func() {
			u, err := url.Parse("https://example.com/?foo=bar&v=1")
			Expect(err).Should(Succeed())

			Expect(u).To(be_url.URL(
				be_url.HavingSearchParam("foo"),
			))
		})

		It("should ensure a given url is a valid *urlURL", func() {
			Expect("http://example.com/foo/bar").To(be_url.URL(
				be_url.TransformUrlFromString,
				be_url.HavingHostname("example.com"),
				be_url.NotHavingPort(),
				be_url.HavingPath("/foo/bar"),
				Not(be_url.HavingRawQuery()),
			))
		})
	})

	Context("with gomock", func() {
		var ctrl *gomock.Controller
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
		})
		AfterEach(func() {
			ctrl.Finish()
		})

		It("should match given url to the urler", func() {
			urler := mocks.NewMockUrler(ctrl)
			urler.EXPECT().SetUrl(
				be_url.URL(
					be_url.HavingHostname("example.com"),
					be_url.HavingScheme("http"),
					be_url.HavingPath("/foo/bar"),
				),
			)

			theUrl, _ := url.Parse("http://example.com/foo/bar")
			urler.SetUrl(theUrl)
		})
	})
})

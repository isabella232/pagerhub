package api_test

import (
	. "github.com/concourse/pagerhub/api"

	"net/http"
	"net/http/httptest"
	"strings"

	"io/ioutil"

	"github.com/concourse/pagerhub/api/apifakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GithubSignatureMiddleware", func() {
	var (
		fakeHandler *apifakes.FakeHandler

		secret string
		body   string

		m GithubSignatureMiddleware

		rw *httptest.ResponseRecorder
		r  *http.Request
	)

	BeforeEach(func() {
		fakeHandler = new(apifakes.FakeHandler)

		secret = "foo"
		body = "some-request-body"
		m = GithubSignatureMiddleware{
			GithubWebhookSecret: secret,
			Inner:               fakeHandler,
		}

		rw = httptest.NewRecorder()

		var err error
		r, err = http.NewRequest("POST", "/some-endpoint", strings.NewReader(body))
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns a Bad Request if the signature is missing", func() {
		m.ServeHTTP(rw, r)

		Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(0))
		Expect(rw.Result().StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("returns a Bad Request if the signature is wrong", func() {
		r.Header.Set("X-Hub-Signature", "sha1=wrong")
		m.ServeHTTP(rw, r)

		Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(0))
		Expect(rw.Result().StatusCode).To(Equal(http.StatusBadRequest))
	})

	Context("when the HMAC of the body matches the X-Hub-Signature header", func() {
		BeforeEach(func() {
			// ruby -ropenssl -e 'puts OpenSSL::HMAC.hexdigest(OpenSSL::Digest.new("sha1"), "foo", "some-request-body")'
			r.Header.Set("X-Hub-Signature", "sha1=1c41d6dfe1b29d0802d9b46d8c1136b8ad0c933b")
		})

		It("verifies the HMAC of the body matches the X-Hub-Signature header", func() {
			m.ServeHTTP(rw, r)

			Expect(fakeHandler.ServeHTTPCallCount()).To(Equal(1))
		})

		It("doesn't prevent inner readers from reading the body", func() {
			fakeHandler.ServeHTTPStub = func(rw http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(b)).To(Equal(body))
			}
		})
	})

})

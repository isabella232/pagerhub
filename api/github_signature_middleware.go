package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
)

//go:generate counterfeiter net/http.Handler
type GithubSignatureMiddleware struct {
	GithubWebhookSecret string
	Inner               http.Handler
}

func (g GithubSignatureMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	const signatureHeader = "X-Hub-Signature"
	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	s := r.Header.Get(signatureHeader)

	if len(s) != signatureLength || !strings.HasPrefix(s, signaturePrefix) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	signature := make([]byte, 20)
	hex.Decode(signature, []byte(s[5:]))

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !checkMAC(b, signature, []byte(g.GithubWebhookSecret)) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	g.Inner.ServeHTTP(rw, r)
}

func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

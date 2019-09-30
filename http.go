package toolbelt

import (
	"net/http"
	"time"
)

// DefaultHeaderTransport is an http.Transport that
// inserts user-agent and x-client-version headers.
type DefaultHeaderTransport struct {
	userAgent, version string
	baseTransport      http.RoundTripper
}

// NewDefaultHeaderTransport returns an initialied DefaultHeaderTransport that uses the http.DefaultTransport .
func NewDefaultHeaderTransport(userAgent, version string) *DefaultHeaderTransport {
	return &DefaultHeaderTransport{
		userAgent:     userAgent,
		version:       version,
		baseTransport: http.DefaultTransport,
	}
}

// SetBaseTransport allows to replace the base transport.
func (dht *DefaultHeaderTransport) SetBaseTransport(baseTransport http.RoundTripper) {
	dht.baseTransport = baseTransport
}

func (dht *DefaultHeaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", dht.userAgent)
	r.Header.Set("X-Client-Version", dht.version)
	return dht.baseTransport.RoundTrip(r)
}

const DefaultTimeout = 10 * time.Second

// NewHTTPClient returns an HTTP client initialized with a DefaultHeaderTransport and with DefaultTimeout.
func NewHTTPClient(userAgent, version string) *http.Client {
	cl := &http.Client{
		Transport: NewDefaultHeaderTransport(userAgent, version),
		Timeout: DefaultTimeout,
	}
	return cl
}

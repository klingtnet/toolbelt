package toolbelt

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHeaderTransport(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.UserAgent()
		if ua != "test-agent" {
			t.Fatalf("expected user-agent to be test-agent but was %q", ua)
		}
		clientVersion := r.Header.Get("x-client-version")
		if clientVersion != "test-version" {
			t.Fatalf("expected x-client-version to be test-version but was %q", clientVersion)
		}
	}))

	cl := NewHTTPClient("test-agent", "test-version")
	resp, err := cl.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestClientTimeout(t *testing.T) {
	cl := NewHTTPClient("test-agent", "test-version")
	if cl.Timeout != DefaultTimeout {
		t.Fatalf("expected client timeout to be %q but was %q", DefaultTimeout, cl.Timeout)
	}
}

package main

import (
	"net/http"
	"testing"
)

func TestSecurityHeaders_AllRoutes(t *testing.T) {
	srv := newTestServer(t)
	_, _ = srv.store.Create("a", "")

	cases := []struct {
		method, target string
		body           any
	}{
		{http.MethodGet, "/health", nil},
		{http.MethodGet, "/metrics", nil},
		{http.MethodGet, "/notes", nil},
		{http.MethodGet, "/notes/1", nil},
		{http.MethodGet, "/notes/999", nil},              // 404 from handler
		{http.MethodPost, "/notes", map[string]string{}}, // 400 from handler
		{http.MethodGet, "/no-such-route", nil},          // 404 from mux
		{http.MethodPut, "/notes", nil},                  // 405 from mux
	}
	want := map[string]string{
		"Content-Security-Policy": "default-src 'none'; frame-ancestors 'none'",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"Referrer-Policy":         "no-referrer",
	}
	for _, c := range cases {
		rec := do(t, srv, c.method, c.target, c.body)
		for k, v := range want {
			if got := rec.Header().Get(k); got != v {
				t.Errorf("%s %s (%d): %s = %q, want %q", c.method, c.target, rec.Code, k, got, v)
			}
		}
	}
}

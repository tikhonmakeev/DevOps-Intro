package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck_OK(t *testing.T) {
	ts := httptest.NewServer(newTestServer(t).Routes())
	defer ts.Close()
	if code := healthcheck(ts.URL + "/health"); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
}

func TestHealthcheck_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()
	if code := healthcheck(ts.URL + "/health"); code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
}

func TestHealthcheck_NoServer(t *testing.T) {
	ts := httptest.NewServer(http.NotFoundHandler())
	url := ts.URL + "/health"
	ts.Close()
	if code := healthcheck(url); code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
}

func TestLocalAddr(t *testing.T) {
	if got := localAddr(":8080"); got != "127.0.0.1:8080" {
		t.Fatalf("localAddr(:8080) = %q", got)
	}
	if got := localAddr("0.0.0.0:9000"); got != "0.0.0.0:9000" {
		t.Fatalf("localAddr(0.0.0.0:9000) = %q", got)
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandlerDefaultsToOneGIF(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	rootHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if !strings.Contains(res.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("content type = %q, want text/html", res.Header().Get("Content-Type"))
	}
	body := res.Body.String()
	if !strings.Contains(body, "<h1>") {
		t.Fatal("page is missing a heading")
	}
	if strings.Count(body, "<img") != 1 {
		t.Fatalf("expected 1 img tag for default count, got %d in %s", strings.Count(body, "<img"), body)
	}
}

func TestRootHandlerHonorsRequestedGIFCount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?count=3", nil)
	res := httptest.NewRecorder()

	rootHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Count(body, "<img") != 3 {
		t.Fatalf("expected 3 img tags for count=3, got %d in %s", strings.Count(body, "<img"), body)
	}
	if !strings.Contains(body, "name=\"count\"") || !strings.Contains(body, "value=\"3\"") {
		t.Fatal("page is missing the count field or selected value")
	}
}

func TestRootHandlerRejectsInvalidGIFCount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?count=42", nil)
	res := httptest.NewRecorder()

	rootHandler(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Count(body, "<img") != 1 {
		t.Fatalf("expected 1 img tag for invalid count, got %d in %s", strings.Count(body, "<img"), body)
	}
}

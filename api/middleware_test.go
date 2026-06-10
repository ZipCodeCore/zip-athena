package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitMiddlewareSkipsStaticContent(t *testing.T) {
	limiter := NewRateLimiter(0, 1)
	handler := RateLimitMiddleware(limiter)(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, requestPath := range []string{"/admin.html", "/js/utils/api.js", "/"} {
		req := httptest.NewRequest(http.MethodGet, requestPath, nil)
		req.RemoteAddr = "127.0.0.1:1234"
		rec := httptest.NewRecorder()

		handler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected static request %q to bypass rate limiting, got %d", requestPath, rec.Code)
		}
	}
}

func TestRateLimitMiddlewareStillLimitsAPIRequests(t *testing.T) {
	limiter := NewRateLimiter(0, 1)
	handler := RateLimitMiddleware(limiter)(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	firstReq := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	firstReq.RemoteAddr = "127.0.0.1:1234"
	firstRec := httptest.NewRecorder()
	handler(firstRec, firstReq)

	if firstRec.Code != http.StatusOK {
		t.Fatalf("expected first API request to succeed, got %d", firstRec.Code)
	}

	secondReq := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	secondReq.RemoteAddr = "127.0.0.1:1234"
	secondRec := httptest.NewRecorder()
	handler(secondRec, secondReq)

	if secondRec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second API request to be rate limited, got %d", secondRec.Code)
	}
}

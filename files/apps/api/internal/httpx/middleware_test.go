package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddlewarePreservesValidID(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RequestID(r) != "request-test" {
			t.Fatalf("request ID was not added to context")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Header.Set("X-Request-ID", "request-test")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Header().Get("X-Request-ID") != "request-test" {
		t.Fatalf("request ID header was not preserved")
	}
}

func TestIdempotencyMiddlewareReplaysResponse(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	calls := 0
	handler := IdempotencyMiddleware(store, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"created":true}`))
	}))
	for range 2 {
		request := httptest.NewRequest(http.MethodPost, "/api/v1/items", nil)
		request.Header.Set("Idempotency-Key", "item-1")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK || response.Body.String() != `{"created":true}` {
			t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
		}
	}
	if calls != 1 {
		t.Fatalf("expected one handler call, got %d", calls)
	}
}

package httpx

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func RequestID(r *http.Request) string {
	requestID, _ := r.Context().Value(requestIDKey).(string)
	return requestID
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if requestID == "" || len(requestID) > 128 {
			requestID = newRequestID()
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, requestID)))
	})
}

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		observed := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(observed, r)
		if observed.status == 0 {
			observed.status = http.StatusOK
		}
		logger.Info("http request", "request_id", RequestID(r), "method", r.Method, "path", r.URL.Path, "status", observed.status, "duration_ms", time.Since(started).Milliseconds())
	})
}

type IdempotencyRecord struct {
	Status  int
	Header  http.Header
	Body    []byte
	Expires time.Time
}

type IdempotencyStore interface {
	Get(key string) (IdempotencyRecord, bool)
	Put(key string, record IdempotencyRecord)
}

// MemoryIdempotencyStore is for local or single-process deployments. Replace
// it with a PostgreSQL-backed implementation for production clusters.
type MemoryIdempotencyStore struct {
	mu      sync.Mutex
	records map[string]IdempotencyRecord
}

func NewMemoryIdempotencyStore() *MemoryIdempotencyStore {
	return &MemoryIdempotencyStore{records: make(map[string]IdempotencyRecord)}
}

func (s *MemoryIdempotencyStore) Get(key string) (IdempotencyRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.records[key]
	if !ok || time.Now().After(record.Expires) {
		delete(s.records, key)
		return IdempotencyRecord{}, false
	}
	return record, true
}

func (s *MemoryIdempotencyStore) Put(key string, record IdempotencyRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[key] = record
}

func IdempotencyMiddleware(store IdempotencyStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if store == nil || !mutating(r.Method) {
			next.ServeHTTP(w, r)
			return
		}
		idempotencyKey := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		if idempotencyKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		if len(idempotencyKey) > 255 {
			JSON(w, http.StatusBadRequest, ErrorResponse{Error: ErrorBody{Code: "invalid_idempotency_key", Message: "Idempotency-Key must be 255 characters or fewer", RequestID: RequestID(r)}})
			return
		}
		key := r.Method + ":" + r.URL.Path + ":" + idempotencyKey
		if record, ok := store.Get(key); ok {
			w.Header().Set("X-Idempotent-Replay", "true")
			copyResponse(w, record.Status, record.Header, record.Body)
			return
		}
		captured := &captureWriter{header: make(http.Header)}
		next.ServeHTTP(captured, r)
		if captured.status < http.StatusInternalServerError {
			store.Put(key, IdempotencyRecord{Status: captured.status, Header: captured.header, Body: captured.body.Bytes(), Expires: time.Now().Add(24 * time.Hour)})
		}
		copyResponse(w, captured.status, captured.header, captured.body.Bytes())
	})
}

func mutating(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}

func copyResponse(w http.ResponseWriter, status int, header http.Header, body []byte) {
	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

type captureWriter struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func (w *captureWriter) Header() http.Header { return w.header }

func (w *captureWriter) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
}

func (w *captureWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.body.Write(body)
}

func newRequestID() string {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "request-unknown"
	}
	return hex.EncodeToString(value)
}

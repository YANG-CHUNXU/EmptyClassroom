package utils

import (
	"EmptyClassroom/logs"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	logs.Init(false)
}

func TestHttpGetWithHeaderHonorsContextCancellation(t *testing.T) {
	requestCanceled := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		close(requestCanceled)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, _, _, err := HttpGetWithHeader(ctx, server.URL, map[string]string{"X-Test": "1"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("HttpGetWithHeader() error = %v, want context deadline exceeded", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("HttpGetWithHeader() took too long to return: %v", elapsed)
	}

	select {
	case <-requestCanceled:
	case <-time.After(time.Second):
		t.Fatal("request context was not canceled on the server side")
	}
}

func TestHttpPostJSONReturnsOnCanceledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, _, err := HttpPostJson(ctx, server.URL, []byte(`{"hello":"world"}`))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("HttpPostJson() error = %v, want context canceled", err)
	}
}

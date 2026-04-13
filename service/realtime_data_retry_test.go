package service

import (
	"errors"
	"testing"
)

func TestShouldRetryRealtime(t *testing.T) {
	if shouldRetryRealtime(nil) {
		t.Fatalf("nil error should not retry")
	}

	if shouldRetryRealtime(errors.New("temporary")) != true {
		t.Fatalf("temporary errors should retry")
	}

	if shouldRetryRealtime(ErrLoginRejected) {
		t.Fatalf("login rejection should not retry")
	}

	if shouldRetryRealtime(errors.Join(ErrLoginRejected, errors.New("wrapped"))) {
		t.Fatalf("wrapped login rejection should not retry")
	}
}

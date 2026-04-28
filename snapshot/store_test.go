package snapshot

import (
	"context"
	"errors"
	"testing"
)

func TestNewDefaultStoreRequiresBlobTokenOnVercel(t *testing.T) {
	t.Setenv("VERCEL", "1")
	t.Setenv(BlobTokenEnvKey, "")

	store := NewDefaultStore()
	_, err := store.Load(context.Background())
	if !errors.Is(err, ErrSnapshotStoreUnavailable) {
		t.Fatalf("Load() error = %v, want ErrSnapshotStoreUnavailable", err)
	}
}

func TestNewDefaultStoreUsesLocalFileOutsideVercel(t *testing.T) {
	t.Setenv("VERCEL", "")
	t.Setenv(BlobTokenEnvKey, "")

	if _, ok := NewDefaultStore().(*fileStore); !ok {
		t.Fatalf("NewDefaultStore() should use file store outside Vercel when %s is empty", BlobTokenEnvKey)
	}
}

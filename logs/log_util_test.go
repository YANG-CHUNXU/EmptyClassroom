package logs

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestWithLogIDPreservesCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	ctx := WithLogID(parent)
	if got := GetLogIDFromContext(ctx); got == "" {
		t.Fatal("WithLogID() did not attach a log id")
	}

	cancel()

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("derived context did not observe parent cancellation")
	}

	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("ctx.Err() = %v, want context canceled", ctx.Err())
	}
}

func TestSetNewContextForGinContextUsesRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	parent, cancel := context.WithCancel(context.Background())
	c.Request = httptest.NewRequest("GET", "/api/get_data", nil).WithContext(parent)

	SetNewContextForGinContext(c)
	ctx := GetContextFromGinContext(c)
	logID := GetLogIDFromContext(ctx)
	if logID == "" {
		t.Fatal("SetNewContextForGinContext() did not store a log id")
	}
	if got := recorder.Header().Get("LogID"); got != logID {
		t.Fatalf("LogID header = %q, want %q", got, logID)
	}

	cancel()

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("gin context did not inherit request cancellation")
	}
}

package service

import (
	"EmptyClassroom/config"
	"EmptyClassroom/logs"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"context"
)

func init() {
	logs.Init(false)
	config.InitConfig()
}

func requireRealtimeIntegrationTests(t *testing.T) {
	t.Helper()

	if os.Getenv("RUN_REALTIME_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_REALTIME_INTEGRATION_TESTS=1 to run realtime integration tests")
	}
	if os.Getenv(LoginUsernameKey) == "" {
		t.Skip("set JW_USERNAME to run realtime integration tests")
	}
	if os.Getenv(LoginPasswordKey) == "" {
		t.Skip("set JW_PASSWORD to run realtime integration tests")
	}
}

func TestLogin(t *testing.T) {
	requireRealtimeIntegrationTests(t)

	if err := Login(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestQueryOne(t *testing.T) {
	requireRealtimeIntegrationTests(t)

	if err := Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := QueryOne(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestQueryAll(t *testing.T) {
	requireRealtimeIntegrationTests(t)

	if err := Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	ans, err := QueryAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	marshal, err := json.Marshal(ans)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}

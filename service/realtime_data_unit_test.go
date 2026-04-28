package service

import (
	"EmptyClassroom/service/model"
	"context"
	"errors"
	"net/http"
	"testing"
)

func withRealtimeRequestStubs(
	t *testing.T,
	login func(context.Context, string, map[string]string) (int, http.Header, []byte, error),
	query func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error),
) {
	t.Helper()

	previousLogin := realtimeLoginRequest
	previousQuery := realtimeQueryRequest
	previousToken := Token

	realtimeLoginRequest = login
	realtimeQueryRequest = query
	Token = ""

	t.Cleanup(func() {
		realtimeLoginRequest = previousLogin
		realtimeQueryRequest = previousQuery
		Token = previousToken
	})
}

func TestRealtimeLoginSuccess(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	var capturedURL string
	var capturedData map[string]string
	withRealtimeRequestStubs(t,
		func(ctx context.Context, url string, data map[string]string) (int, http.Header, []byte, error) {
			capturedURL = url
			capturedData = data
			return 200, nil, []byte(`{"code":"1","data":{"token":"test-token"}}`), nil
		},
		func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error) {
			t.Fatal("query request should not be called")
			return 0, nil, nil, nil
		},
	)

	if err := Login(context.Background()); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if capturedURL != LoginURL {
		t.Fatalf("Login() url = %q, want %q", capturedURL, LoginURL)
	}
	if capturedData["userNo"] != "user" {
		t.Fatalf("Login() userNo = %q, want %q", capturedData["userNo"], "user")
	}
	if capturedData["pwd"] == "" {
		t.Fatal("Login() should send encrypted password")
	}
	if Token != "test-token" {
		t.Fatalf("Token = %q, want %q", Token, "test-token")
	}
}

func TestRealtimeLoginRejected(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	withRealtimeRequestStubs(t,
		func(context.Context, string, map[string]string) (int, http.Header, []byte, error) {
			return 200, nil, []byte(`{"code":"0","Msg":"invalid credentials"}`), nil
		},
		func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error) {
			t.Fatal("query request should not be called")
			return 0, nil, nil, nil
		},
	)

	err := Login(context.Background())
	if !errors.Is(err, ErrLoginRejected) {
		t.Fatalf("Login() error = %v, want ErrLoginRejected", err)
	}
}

func TestRealtimeLoginHTTPStatusFailure(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	withRealtimeRequestStubs(t,
		func(context.Context, string, map[string]string) (int, http.Header, []byte, error) {
			return http.StatusBadGateway, nil, []byte(`{}`), nil
		},
		func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error) {
			t.Fatal("query request should not be called")
			return 0, nil, nil, nil
		},
	)

	if err := Login(context.Background()); err == nil {
		t.Fatal("Login() error = nil, want failure")
	}
}

func TestRealtimeLoginInvalidJSON(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	withRealtimeRequestStubs(t,
		func(context.Context, string, map[string]string) (int, http.Header, []byte, error) {
			return 200, nil, []byte(`not-json`), nil
		},
		func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error) {
			t.Fatal("query request should not be called")
			return 0, nil, nil, nil
		},
	)

	if err := Login(context.Background()); err == nil {
		t.Fatal("Login() error = nil, want JSON unmarshal failure")
	}
}

func TestRealtimeLoginRequestFailure(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	wantErr := errors.New("network down")
	withRealtimeRequestStubs(t,
		func(context.Context, string, map[string]string) (int, http.Header, []byte, error) {
			return 0, nil, nil, wantErr
		},
		func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error) {
			t.Fatal("query request should not be called")
			return 0, nil, nil, nil
		},
	)

	if err := Login(context.Background()); !errors.Is(err, wantErr) {
		t.Fatalf("Login() error = %v, want %v", err, wantErr)
	}
}

func TestRealtimeQueryOneSuccess(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	loginCalls := 0
	queryCalls := 0
	withRealtimeRequestStubs(t,
		func(context.Context, string, map[string]string) (int, http.Header, []byte, error) {
			loginCalls++
			return 200, nil, []byte(`{"code":"1","data":{"token":"query-token"}}`), nil
		},
		func(ctx context.Context, url string, data map[string]string, header map[string]string) (int, http.Header, []byte, error) {
			queryCalls++
			if url != QueryURL {
				t.Fatalf("QueryOne() url = %q, want %q", url, QueryURL)
			}
			if got := data["campusId"]; got != "01" {
				t.Fatalf("QueryOne() campusId = %q, want %q", got, "01")
			}
			if got := header["token"]; got != "query-token" {
				t.Fatalf("QueryOne() token header = %q, want %q", got, "query-token")
			}
			return 200, nil, []byte(`{"code":"1","data":[{"CLASSROOMS":"N101","NODETIME":"1-2","NODENAME":"上午"}]}`), nil
		},
	)

	got, err := QueryOne(context.Background(), 1)
	if err != nil {
		t.Fatalf("QueryOne() error = %v", err)
	}
	if loginCalls != 1 {
		t.Fatalf("Login() call count = %d, want 1", loginCalls)
	}
	if queryCalls != 1 {
		t.Fatalf("query call count = %d, want 1", queryCalls)
	}
	want := []model.JWClassInfo{{Classrooms: "N101", NodeTime: "1-2", NodeName: "上午"}}
	if len(got) != len(want) {
		t.Fatalf("QueryOne() len = %d, want %d", len(got), len(want))
	}
	if got[0] != want[0] {
		t.Fatalf("QueryOne() first item = %+v, want %+v", got[0], want[0])
	}
}

func TestRealtimeQueryOneSkipsQueryWhenLoginFails(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	queryCalled := false
	wantErr := errors.Join(ErrLoginRejected, errors.New("invalid credentials"))
	withRealtimeRequestStubs(t,
		func(context.Context, string, map[string]string) (int, http.Header, []byte, error) {
			return 0, nil, nil, wantErr
		},
		func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error) {
			queryCalled = true
			return 0, nil, nil, nil
		},
	)

	_, err := QueryOne(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("QueryOne() error = %v, want %v", err, wantErr)
	}
	if queryCalled {
		t.Fatal("QueryOne() should not call query endpoint when login fails")
	}
}

func TestRealtimeQueryOneDoesNotRetryTemporaryLoginFailure(t *testing.T) {
	t.Setenv(LoginUsernameKey, "user")
	t.Setenv(LoginPasswordKey, "password")

	loginCalls := 0
	queryCalled := false
	wantErr := errors.New("network down")
	withRealtimeRequestStubs(t,
		func(context.Context, string, map[string]string) (int, http.Header, []byte, error) {
			loginCalls++
			return 0, nil, nil, wantErr
		},
		func(context.Context, string, map[string]string, map[string]string) (int, http.Header, []byte, error) {
			queryCalled = true
			return 0, nil, nil, nil
		},
	)

	_, err := QueryOne(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("QueryOne() error = %v, want %v", err, wantErr)
	}
	if loginCalls != 1 {
		t.Fatalf("login calls = %d, want 1", loginCalls)
	}
	if queryCalled {
		t.Fatal("QueryOne() should not call query endpoint when login fails")
	}
}

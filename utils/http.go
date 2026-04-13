package utils

import (
	"EmptyClassroom/logs"
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	OutboundHTTPTimeoutEnvKey  = "OUTBOUND_HTTP_TIMEOUT"
	defaultOutboundHTTPTimeout = 15 * time.Second
)

var (
	outboundHTTPClientOnce sync.Once
	outboundHTTPClient     *http.Client
)

func OutboundHTTPClient() *http.Client {
	outboundHTTPClientOnce.Do(func() {
		outboundHTTPClient = &http.Client{
			Timeout: outboundHTTPTimeout(),
		}
	})
	return outboundHTTPClient
}

func outboundHTTPTimeout() time.Duration {
	raw := os.Getenv(OutboundHTTPTimeoutEnvKey)
	if raw == "" {
		return defaultOutboundHTTPTimeout
	}

	timeout, err := time.ParseDuration(raw)
	if err != nil || timeout <= 0 {
		log.Printf("[Warn] invalid %s=%q, falling back to %s", OutboundHTTPTimeoutEnvKey, raw, defaultOutboundHTTPTimeout)
		return defaultOutboundHTTPTimeout
	}

	return timeout
}

func HttpPostJson(ctx context.Context, url string, jsonStr []byte) (int, http.Header, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logs.CtxError(ctx, "http post json error: %v", err)
		return 0, nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := OutboundHTTPClient().Do(req)
	if err != nil {
		logs.CtxError(ctx, "http post json error: %v", err)
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	header := resp.Header
	body, _ := io.ReadAll(resp.Body)
	return statusCode, header, body, nil
}

func HttpPostForm(ctx context.Context, url string, data map[string]string) (int, http.Header, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		logs.CtxError(ctx, "http post form error: %v", err)
		return 0, nil, nil, err
	}
	q := req.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := OutboundHTTPClient().Do(req)
	if err != nil {
		logs.CtxError(ctx, "http post form error: %v", err)
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	header := resp.Header
	body, _ := io.ReadAll(resp.Body)
	return statusCode, header, body, nil
}

func HttpPostFormWithHeader(ctx context.Context, url string, data map[string]string, header map[string]string) (int, http.Header, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		logs.CtxError(ctx, "http post form with header error: %v", err)
		return 0, nil, nil, err
	}
	q := req.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := OutboundHTTPClient().Do(req)
	if err != nil {
		logs.CtxError(ctx, "http post form with header error: %v", err)
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)
	return statusCode, resp.Header, body, nil
}

func HttpGetWithHeader(ctx context.Context, url string, header map[string]string) (int, http.Header, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logs.CtxError(ctx, "http get with header error: %v", err)
		return 0, nil, nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := OutboundHTTPClient().Do(req)
	if err != nil {
		logs.CtxError(ctx, "http get with header error: %v", err)
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)
	return statusCode, resp.Header, body, nil
}

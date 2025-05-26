package traefik_forward_slash_redirector_test //nolint:revive,stylecheck

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	traefikforwardslashredirector "github.com/rjop-hccgt/traefik-forward-slash-redirector"
)

func TestForwardSlash_isFile(_ *testing.T) {
	infoLogger := log.Default()
	redirector := &traefikforwardslashredirector.ForwardSlash{
		InfoLogger: infoLogger,
	}

	isFile := redirector.IsFile("/")
	if isFile {
		log.Fatalf("'/' is not a file")
	}

	isFile = redirector.IsFile("/test.jpg")
	if !isFile {
		log.Fatalf("'/test.jpg' is a file")
	}
}

func TestForwardSlash_ServeHTTPPermanent(t *testing.T) {
	cfg := traefikforwardslashredirector.CreateConfig()
	cfg.Permanent = true
	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler, err := traefikforwardslashredirector.New(ctx, next, cfg, "forward-slash-redirector-test")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)
	assertPath(t, req, "/")
	assertHTTPResponseCode(t, recorder.Result(), 200)

	recorder = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/index.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)
	assertPath(t, req, "/index.jpg")
	assertHTTPResponseCode(t, recorder.Result(), 200)

	recorder = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/path", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)
	assertPath(t, req, "/path/")
	assertHTTPResponseCode(t, recorder.Result(), 301)
}

func TestForwardSlash_ServeHTTPTemporary(t *testing.T) {
	cfg := traefikforwardslashredirector.CreateConfig()
	cfg.Permanent = false
	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler, err := traefikforwardslashredirector.New(ctx, next, cfg, "forward-slash-redirector-test")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)
	assertPath(t, req, "/test/")
	assertHTTPResponseCode(t, recorder.Result(), 302)

	recorder = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/test?q=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)
	assertPath(t, req, "/test/?q=1")
	assertHTTPResponseCode(t, recorder.Result(), 302)
}

func assertPath(t *testing.T, req *http.Request, expected string) {
	t.Helper()
	log.Printf("Validating %v", req.URL.Path)
	if !strings.EqualFold(req.URL.Path, expected) {
		t.Errorf("invalid path: %s expected: %s", req.URL.Path, expected)
	}
}

func assertHTTPResponseCode(t *testing.T, response *http.Response, statusCode int) {
	t.Helper()
	if response.StatusCode != statusCode {
		t.Errorf("invalid responseCode: %d expected: %d", response.StatusCode, statusCode)
	}
}

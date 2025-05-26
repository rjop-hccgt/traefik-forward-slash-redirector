package traefik_forward_slash_redirector

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Config struct {
	Permanent bool `json:"permanent,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Permanent: false,
	}
}

type ForwardSlash struct {
	Next       http.Handler
	Permanent  bool
	Name       string
	InfoLogger *log.Logger
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	infoLogger := log.Default()

	return &ForwardSlash{
		Permanent:  config.Permanent,
		Next:       next,
		Name:       name,
		InfoLogger: infoLogger,
	}, nil
}

func (a *ForwardSlash) IsFile(relativePath string) bool {
	// Clean the path to handle cases like "/test/../file.jpg"
	cleanedPath := filepath.Clean(relativePath)
	a.InfoLogger.Printf("Cleaned path: %s", cleanedPath)

	// Get the base name of the path (e.g., "test.jpg" from "/dir/test.jpg")
	base := filepath.Base(cleanedPath)
	a.InfoLogger.Printf("Base Path path: %s", base)
	pos := strings.LastIndex(base, ".")
	return pos != -1
}

func (a *ForwardSlash) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, "/") {
		a.InfoLogger.Printf("Path has final '/': %s", req.URL.Path)
		a.Next.ServeHTTP(rw, req)
	} else {
		if a.IsFile(req.URL.Path) {
			a.InfoLogger.Printf("Path is a file: %s", req.URL.Path)
			a.Next.ServeHTTP(rw, req)
		} else {
			req.URL.Path += "/"
			if req.URL.RawQuery != "" {
				req.URL.Path += "?" + req.URL.RawQuery
			}
			if a.Permanent {
				http.Redirect(rw, req, req.URL.Path, http.StatusMovedPermanently)
			} else {
				http.Redirect(rw, req, req.URL.Path, http.StatusFound)
			}
		}
	}

}

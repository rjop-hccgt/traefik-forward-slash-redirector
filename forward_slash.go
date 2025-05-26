package traefik_forward_slash_redirector //nolint:revive

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	Permanent bool `json:"permanent,omitempty"`
}

// CreateConfig creates a base configuration.
func CreateConfig() *Config {
	return &Config{
		Permanent: false,
	}
}

// ForwardSlash the main plugin.
type ForwardSlash struct {
	Next       http.Handler
	Permanent  bool
	Name       string
	InfoLogger *log.Logger
}

// New creates a new ForwardSlash plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	infoLogger := log.New(os.Stdout, "ForwardSlash: ", log.Ldate|log.Ltime)

	return &ForwardSlash{
		Permanent:  config.Permanent,
		Next:       next,
		Name:       name,
		InfoLogger: infoLogger,
	}, nil
}

// IsFile checks whether a given path is a file.
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
	if strings.HasSuffix(req.URL.Path, "/") || a.IsFile(req.URL.Path) {
		a.InfoLogger.Printf("Path has final '/' or is a file: %s", req.URL.Path)
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

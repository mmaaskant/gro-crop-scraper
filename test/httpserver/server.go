package httpserver

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"testing"
)

//go:embed html
var html embed.FS

func StartTestHttpServer(t *testing.T) (baseUrl string) {
	fSys, err := fs.Sub(html, "html")
	if err != nil {
		t.Errorf("HTTP test server failed to access HTML files, error: %s", err)
	}
	port := os.Getenv("HTTP_TEST_SERVER_PORT")
	url := fmt.Sprintf("localhost:%s", port)
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(fSys)))
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != http.ErrServerClosed {
			t.Errorf("Test HTTP server error: %s", err)
		}
	}()
	return url
}

package server

import(
	"net/http"
	"github.com/Kamalpreet-singh007/url-shortener/internals/handler"
	"fmt"
)

func SetupRoutes(mux *http.ServeMux, urlHandler *handler.UrlHandler) *http.ServeMux {
	mux.HandleFunc("POST /api/shorten", urlHandler.Shorten)
	mux.HandleFunc("GET /{code}", urlHandler.Redirect)
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})
	return mux
}
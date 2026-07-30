package server

import(
	"net/http"
	"log"
	"time"
	"github.com/Kamalpreet-singh007/url-shortener/internals/store"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func (){
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		ctx := store.ContextWithCacheStatus(r.Context())
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		cacheStatus := store.GetCacheStatus(ctx)
		if cacheStatus != "" {
			log.Printf("%s %s -> %d (%v) cache=%s", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start), cacheStatus)
			return
		}
		log.Printf("%s %s -> %d (%v)", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}


func WrapHandler(handler http.Handler) http.Handler {
	wrappedHandler := loggingMiddleware(recoveryMiddleware(handler))
	return wrappedHandler
}
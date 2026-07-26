package main

import(
	"log"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"context"
	"github.com/joho/godotenv"
    _"github.com/jackc/pgx/v5/stdlib"
	"database/sql"


	"github.com/Kamalpreet-singh007/url-shortener/internals/store"
	"github.com/Kamalpreet-singh007/url-shortener/internals/handler"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode: 200,
		}

		next.ServeHTTP(wrapped, r)

		log.Printf("%s %s -> %d (%v)", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}
func main(){

	godotenv.Load()
	db_url := os.Getenv("DB_URL")	
	port := os.Getenv("PORT")


	
	// db connection ___________________
	db , err :=  sql.Open("pgx", db_url)
	if err != nil {
		log.Fatalf("could not open DB : %s", err)
	}
	defer db.Close()
	
	if err = db.Ping(); err != nil {
		log.Fatalf("could not connect to db: %s", err)
	}




	
	urlStore  := store.NewPostgresStore(db)
	UrlHandler := handler.NewUrlHandler(urlStore)
	
	log.Println("db connected succesfully")
	
	mux := http.NewServeMux()
	loggingMux := loggingMiddleware(mux)
	
    mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status":"ok"}`)
    })
	mux.HandleFunc("POST /api/shorten", UrlHandler.Shorten)
	mux.HandleFunc("GET /{code}", UrlHandler.Redirect)
	
	// server connection
	if port =="" {
		port ="8080"
	}

	srv := http.Server{
		Addr:         ":" + port,
		Handler:      loggingMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
	}

	log.Printf("server listening on http://localhost:%s",port)
	go func(){
		if err  := srv.ListenAndServe();err  != nil && err != http.ErrServerClosed{
			log.Fatalf("server error: %v", err)
		}
	}()


	quit := make(chan os.Signal,1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)		
	}
	log.Println("server exited properly")


}
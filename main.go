package main

import(
	"log"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"context"
	"syscall"
	
    _"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
func main(){

	// load env variables
	if err := godotenv.Load(); err != nil {
    	log.Printf("no .env file found, relying on OS environment: %v", err)
	}
	db_url := os.Getenv("DB_URL")	
	port := os.Getenv("PORT")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if db_url == "" {
    log.Fatal("DB_URL is not set")
	}

	if port == "" {
		port = "8080"
	}
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	if redisPassword == "" {
		log.Println("REDIS_PASSWORD is not set, assuming no password is required")
	}

	//redis connection

	rdb := redis.NewClient(&redis.Options{
    Addr:     redisAddr,
    Password: redisPassword,
    DB:       0,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
    	log.Fatal("could not connect to redis:", err)
	}

	log.Println("connected to redis successfully")

	// db connection ___________________
	db , err :=  sql.Open("pgx", db_url)
	if err != nil {
		log.Fatalf("could not open DB : %s", err)
	}
	defer db.Close()
	
	if err = db.Ping(); err != nil {
		log.Fatalf("could not connect to db: %s", err)
	}

	redisCache := store.NewRedisCache(rdb)
	urlStore  := store.NewPostgresStore(db)
	cacheStore := store.NewCacheStore(urlStore, redisCache)
	UrlHandler := handler.NewUrlHandler(cacheStore)

	log.Println("db connected succesfully")
	

	//mux and server setup
	mux := http.NewServeMux()

	handlerChain := loggingMiddleware(recoveryMiddleware(mux))

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
		Handler:      handlerChain,
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

	// graceful shutdown
	quit := make(chan os.Signal,1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)		
	}
	log.Println("server exited properly")


}
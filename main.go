package main

import(
	"log"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
	"log/slog"
    _"github.com/jackc/pgx/v5/stdlib"
	"database/sql"


	"github.com/Kamalpreet-singh007/url-shortener/internals/store"
	"github.com/Kamalpreet-singh007/url-shortener/internals/handler"
)

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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))


	
	urlStore  := store.NewPostgresStore(db)
	UrlHandler := handler.NewUrlHandler(urlStore, logger)
	
	log.Println("db connected succesfully")
	
	mux := http.NewServeMux()
	
	
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
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
	}

	log.Printf("server listening on http://localhost:%s",port)

	if err := srv.ListenAndServe();err  != nil{
		log.Fatalf("server error: %v", err)
	}


}
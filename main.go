package main

import(
	"log"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
	"database/sql"
    _"github.com/jackc/pgx/v5/stdlib"
)

func main(){
	godotenv.Load()
	db_url := os.Getenv("DB_URL")	
	port := os.Getenv("PORT")

	// server connection
	if port =="" {
		port ="8080"
	}
	// db connection ___________________
	db , err :=  sql.Open("pgx", db_url)
	if err != nil {
		log.Fatalf("could not open DB : %s", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("could not connect to db: %s", err)
	}

	log.Println("db connected succesfully")

	mux := http.NewServeMux()


    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status":"ok"}`)
    })

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
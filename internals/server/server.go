package server


import(
	"net/http"
	"github.com/Kamalpreet-singh007/url-shortener/internals/handler"
	"time"
	"context"
	"os"
	"os/signal"
	"syscall"
	"log"
	"github.com/Kamalpreet-singh007/url-shortener/internals/config"
)


func StartServer(urlHandler *handler.UrlHandler, cfg *config.Config) error {
	mux := http.NewServeMux()
	SetupRoutes(mux, urlHandler)
	wrappedHandler := WrapHandler(mux)
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrappedHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println("listening on port", cfg.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start server: %s", err)
		}
	}()
	quit := make(chan os.Signal,1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server: %s", err)
	}
	log.Println("server exited gracefully")
	return nil
}
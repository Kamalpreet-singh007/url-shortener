package main

import(
	"log"
	
	"github.com/Kamalpreet-singh007/url-shortener/internals/config"
	"github.com/Kamalpreet-singh007/url-shortener/internals/app"
)


func main(){

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %s", err)
	}

	err = app.Run(cfg)
	if err != nil {
		log.Fatalf("could not run app: %s", err)
	}
}
package config

import(	
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"fmt"
)

type Config struct {
	Port string
	Db_Url string
	Redis_Url string
	Redis_Password string
	Redis_DB int
}

func Load()(*Config, error){
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found, relying on OS environment: %v", err)
	}
	Port := os.Getenv("PORT")
	Db_Url := os.Getenv("DB_URL")
	Redis_Url := os.Getenv("REDIS_URL")
	Redis_Password := os.Getenv("REDIS_PASSWORD")
	Redis_DB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	if Port == "" {
		Port = "8080"
	}

	
	if Db_Url == "" {
		return nil,fmt.Errorf("DB_URL is not set")
	}
	if Redis_Url == "" {
		return nil,fmt.Errorf("REDIS_URL is not set")
	}


	return &Config{
		Port: Port,
		Db_Url: Db_Url,
		Redis_Url: Redis_Url,
		Redis_Password: Redis_Password,
		Redis_DB: Redis_DB,
	}, nil
}
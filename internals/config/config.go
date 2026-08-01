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
	Redis_Addr string
	Redis_Password string
	Redis_DB int
}

func Load()(*Config, error){
	app_env := os.Getenv("APP_ENV")
	if app_env != "docker" {
		if err := godotenv.Load(); err != nil {
			log.Printf("no .env file found, relying on OS environment")
		}
	}


	Port := os.Getenv("PORT")

	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_port := os.Getenv("POSTGRES_PORT")
	postgres_host := os.Getenv("POSTGRES_HOST")
	Db_Url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", postgres_user, postgres_password, postgres_host, postgres_port, postgres_db)

	Redis_Addr := os.Getenv("REDIS_ADDR")
	Redis_Password := os.Getenv("REDIS_PASSWORD")
	redis_db_raw := os.Getenv("REDIS_DB")
	Redis_DB := 0
	if redis_db_raw != "" {
		var err error
		Redis_DB, err = strconv.Atoi(redis_db_raw)
		if err != nil {
			return nil, fmt.Errorf("REDIS_DB must be an integer, got %q: %w", redis_db_raw, err)
		}
	}
	if Port == "" {
		Port = "8080"
	}
	if postgres_user == "" {
		return nil,fmt.Errorf("POSTGRES_USER is not set")
	}
	if postgres_password == "" {
		return nil,fmt.Errorf("POSTGRES_PASSWORD is not set")
	}
	if postgres_db == "" {
		return nil,fmt.Errorf("POSTGRES_DB is not set")
	}
	if postgres_port == "" {
		return nil,fmt.Errorf("POSTGRES_PORT is not set")
	}
	if postgres_host == "" {
		return nil,fmt.Errorf("POSTGRES_HOST is not set")
	}
	if Redis_Addr == "" {
		return nil,fmt.Errorf("REDIS_URL is not set")
	}


	return &Config{
		Port: Port,
		Db_Url: Db_Url,
		Redis_Addr: Redis_Addr,
		Redis_Password: Redis_Password,
		Redis_DB: Redis_DB,
	}, nil
}
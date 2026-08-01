package app

import(
	"github.com/Kamalpreet-singh007/url-shortener/internals/config"
	"github.com/Kamalpreet-singh007/url-shortener/internals/store"
	"github.com/Kamalpreet-singh007/url-shortener/internals/handler"
	"github.com/Kamalpreet-singh007/url-shortener/internals/server"
	"github.com/redis/go-redis/v9"
	_"github.com/jackc/pgx/v5/stdlib"

	"context"
	"log"
	"database/sql"

)

func Run(cfg *config.Config)error{
	log.Println("starting the application")
	db , err :=  sql.Open("pgx", cfg.Db_Url)
	if err != nil {
		log.Fatalf("could not open DB : %s", err)
	}
	defer db.Close()
	
	if err = db.Ping(); err != nil {
		log.Fatalf("could not connect to db: %s", err)
	}
	log.Println("db connected succesfully")

	rdb := redis.NewClient(&redis.Options{
    Addr:     cfg.Redis_Addr,
    Password: cfg.Redis_Password,
    DB:       cfg.Redis_DB,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
    	log.Fatal("could not connect to redis:", err)
	}
	log.Println("connected to redis successfully")

	postgresStore := store.NewPostgresStore(db)
	redisCache := store.NewRedisCache(rdb)
	cachedUrlStore := store.NewCachedUrlStore(postgresStore, redisCache)
	urlHandler := handler.NewUrlHandler(cachedUrlStore)
	
	return server.StartServer(urlHandler, cfg)


}
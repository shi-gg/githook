package main

import (
	"net/http"

	"github.com/shi-gg/githook/config"
	"github.com/shi-gg/githook/routes"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

var conf = config.Get()

var client = redis.NewClient(&redis.Options{
	Addr:     conf.Redis.Addr,
	Password: conf.Redis.Password,
	Username: conf.Redis.Username,
	DB:       conf.Redis.Db,
})

func main() {
	http.HandleFunc("POST /incoming/{id}", func(w http.ResponseWriter, r *http.Request) {
		routes.HandleIncoming(w, r, client)
	})

	http.HandleFunc("GET /create", routes.HandleCreate)

	http.ListenAndServe(":8080", nil)

	defer client.Close()
}

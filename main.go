package main

import (
	"html/template"
	"net/http"
	"os"

	"CacheLite/api"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

	router := mux.NewRouter()

	router.HandleFunc("/", handler)
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		api.GetUsers(client, w, r)
	}).Methods("GET")

	http.ListenAndServe(":8080", handlers.CORS(corsOrigins, corsMethods, corsHeaders)(router))

}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/index.html")
	t.Execute(w, nil)
}

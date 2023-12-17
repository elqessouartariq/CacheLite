package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"CacheLite/api"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	if mongoClient != nil {
		fmt.Println("Connected to MongoDB!")
	}

	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

	router := mux.NewRouter()

	router.HandleFunc("/", handler)

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users := api.GetUsers(client, w, r)
		json.NewEncoder(w).Encode(users)
	}).Methods("GET")

	router.HandleFunc("/posts-with-cache", func(w http.ResponseWriter, r *http.Request) {
		posts, err := api.GetPostsWithCache(client, mongoClient, w, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")

	router.HandleFunc("/posts-without-cache", func(w http.ResponseWriter, r *http.Request) {
		posts, _ := api.GetPostsWithoutCache(mongoClient, w, r)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")

	http.ListenAndServe(":8080", handlers.CORS(corsOrigins, corsMethods, corsHeaders)(router))

}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/index.html")
	t.Execute(w, nil)
}

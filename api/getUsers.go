package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var ctx = context.Background()

func GetUsers(client *redis.Client, w http.ResponseWriter, r *http.Request) {

	val, err := client.Get(ctx, "users").Result()
	if err == redis.Nil {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/users")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var users []User

		err = json.NewDecoder(resp.Body).Decode(&users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		usersJson, _ := json.Marshal(users)
		client.Set(ctx, "users", usersJson, 20*time.Second)

		json.NewEncoder(w).Encode(users)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		var users []User
		json.Unmarshal([]byte(val), &users)
		json.NewEncoder(w).Encode(users)
	}

}

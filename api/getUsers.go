package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"CacheLite/entities"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func GetUsers(client *redis.Client, w http.ResponseWriter, r *http.Request) []entities.User {
	var users []entities.User

	val, err := client.Get(ctx, "users").Result()
	if err == redis.Nil {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/users")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		usersJson, _ := json.Marshal(users)
		client.Set(ctx, "users", usersJson, 20*time.Second)

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		var users []entities.User
		json.Unmarshal([]byte(val), &users)
		json.NewEncoder(w).Encode(users)
	}

	return users

}

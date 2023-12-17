package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Post struct {
	ID       int       `json:"id"`
	UserID   int       `json:"userId"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID     int    `json:"id"`
	PostID int    `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

var postsFetched = false
var posts []Post

func GetPostsWithoutCache(mongoClient *mongo.Client, w http.ResponseWriter, r *http.Request) ([]Post, error) {
	postCollection := mongoClient.Database("proxy-db").Collection("posts")
	count, _ := postCollection.CountDocuments(context.Background(), bson.M{})

	if !postsFetched && count == 0 {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&posts)
		if err != nil {
			panic(err)
		}

		for i, post := range posts {
			var comments []Comment
			resp, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d/comments", post.ID))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&comments)
			if err != nil {
				panic(err)
			}

			posts[i].Comments = comments

			_, err = postCollection.InsertOne(context.Background(), posts[i])
			if err != nil {
				panic(err)
			}
		}

		postsFetched = true
	} else {
		query := bson.M{}
		cursor, err := postCollection.Find(context.Background(), query)
		if err != nil {
			panic(err)
		}

		if err = cursor.All(context.Background(), &posts); err != nil {
			panic(err)
		}
	}

	return posts, nil
}

func GetPostsWithCache(client *redis.Client, mongoClient *mongo.Client, w http.ResponseWriter, r *http.Request) ([]Post, error) {
	var posts []Post
	cachedPosts, err := client.Get(context.Background(), "posts").Result()

	if err == redis.Nil {
		fmt.Println("Cache miss")
		posts, err = GetPostsWithoutCache(mongoClient, w, r)

		if err != nil {
			panic(err)
		}

		postsJson, err := json.Marshal(posts)
		if err != nil {
			panic(err)
		}

		err = client.Set(context.Background(), "posts", postsJson, 20*time.Second).Err()
		if err != nil {
			panic(err)
		}

	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("Cache hit")

		err = json.Unmarshal([]byte(cachedPosts), &posts)
		if err != nil {
			panic(err)
		}

	}
	return posts, nil
}

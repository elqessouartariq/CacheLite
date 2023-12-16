# CachLite

CachLite is designed to fetch and cache data from a MongoDB database and an external API. It fetches posts, their associated comments, and user data, and caches them in Redis for faster subsequent access.

The `GetPosts` function fetches posts and their associated comments from a MongoDB database, while the `GetUsers` function fetches user data from an external API (https://jsonplaceholder.typicode.com/users in this case).

When a request is made to fetch posts, comments, or users, CachLite first checks if the data is available in the Redis cache. If it is, the cached data is returned. If not, CachLite fetches the data from the MongoDB database or the external API, stores it in the cache for future requests, and returns it in the response.

The data is cached for 20 seconds. After this time, the data is automatically removed from the cache, and the next request for data will fetch it from the MongoDB database or the external API and update the cache.

## Prerequisites

- Docker
- Go
- MongoDB
- Redis

## Installation

### MongoDB

You can install MongoDB using Docker with the following command:

```bash
docker run --name mongo -p 27017:27017 -d mongodb/mongodb-community-server:5.0-ubi8
```

### RedisStak

```bash
docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```

## Usage

To run CachLite, use the following command:

```bash
go run main.go
```

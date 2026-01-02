package pkg

import "github.com/redis/go-redis/v9"

type Cache struct {
	Client *redis.Client
}

func NewCache(addr string) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return Cache{
		Client: rdb,
	}
}

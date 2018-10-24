package session

import (
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

// NewSession returns pointer to the new session
func NewSession(addr string, pswd string) error {
	return initRedis(addr, pswd)
}

func initRedis(addr string, pswd string) error {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pswd, // no password set
		DB:       0,    // use default DB
	})

	_, err := client.Ping().Result()

	return err
}

// SetRecord to the redis, you can choose between persistent record (duration 0) and timed record (duration > 0)
func SetRecord(key string, value string, time time.Duration) error {
	return client.Set(key, value, 0).Err()
}

// RemoveRecord from the redis
func RemoveRecord(key string) error {
	return client.Del(key).Err()
}

// GetRecord from the redis
func GetRecord(key string) (string, error) {
	r := client.Get(key)
	return r.Val(), r.Err()
}

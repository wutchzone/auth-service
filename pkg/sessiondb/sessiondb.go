package sessiondb

import (
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

// NewSession returns pointer to the new sessiondb
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
	var err error
	err = client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	err = client.Set(value, key, 0).Err()
	return err
}

// RemoveRecord from the redis
func RemoveRecord(key string) error {
	var err error
	r, err := GetRecord(key)

	if err != nil {
		return err
	}

	err = client.Del(key).Err()
	if err != nil {
		return err
	}

	err = client.Del(r).Err()
	return err
}

// GetRecord from the redis
func GetRecord(key string) (string, error) {
	r := client.Get(key)

	if r.Err() != nil {

		return "", r.Err()
	}
	return r.Val(), nil
}

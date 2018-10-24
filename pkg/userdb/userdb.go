package userdb

import (
	"encoding/json"
	"strings"

	"github.com/go-redis/redis"
	"github.com/wutchzone/auth-service/pkg/user"
)

var client *redis.Client

// NewSession returns pointer to the new session
func NewSession(addr string, pswd string) error {
	err := initRedis(addr, pswd)
	if err != nil {
		return err
	}
	return nil
}

func initRedis(addr string, pswd string) error {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pswd, // no password set
		DB:       1,    // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		return err
	}
	return nil
}

// GetUser from DB
func GetUser(name string) (*user.User, error) {
	val, err := client.Get(name).Result()
	if err != nil {
		return nil, err
	}
	var user user.User
	json.NewDecoder(strings.NewReader(val)).Decode(&user)

	return &user, nil
}

//SaveUser to the DB
func SaveUser(user user.User) error {
	su, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return client.Set(user.Username, su, 0).Err()
}

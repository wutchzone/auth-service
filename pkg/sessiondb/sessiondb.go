package sessiondb

import (
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

type SessionDB struct {
	Client *redis.Client
}

// NewSession returns pointer to the new sessiondb
func NewSessionDB(addr string, pswd string, table int) (*SessionDB, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pswd,  // no password set
		DB:       table, // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &SessionDB{Client: client}, nil
}

// SetRecord to the redis, you can choose between persistent record (duration 0) and timed record (duration > 0)
func (d *SessionDB) SetRecord(key string, value interface{}, time time.Duration) error {
	if err := d.Client.Set(key, value, time).Err(); err != nil {
		return err
	}
	return nil
}

// RemoveRecord from the redis
func (d *SessionDB) RemoveRecord(key string) error {
	if err := d.Client.Del(key).Err(); err != nil {
		return err
	}
	return nil
}

// GetRecord from the redis
func (d *SessionDB) GetRecord(key string) (interface{}, error) {
	r := d.Client.Get(key)

	if r.Err() != nil {
		return nil, r.Err()
	}
	return r.Val(), nil
}

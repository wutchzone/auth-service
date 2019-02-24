package sessiondb

import (
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

type SessionDB struct {
	client *redis.Client
}

// NewSession returns pointer to the new sessiondb
func NewSessionDB(addr string, pswd string, table int) (*SessionDB, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pswd, // no password set
		DB:       table,    // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &SessionDB{client: client}, nil
}

// SetRecord to the redis, you can choose between persistent record (duration 0) and timed record (duration > 0)
func (d *SessionDB) setRecord(key string, value interface{}, time time.Duration) error {
	if err := d.client.Set(key, value, time).Err(); err != nil {
		return err
	}
	return nil
}

// RemoveRecord from the redis
func (d *SessionDB) removeRecord(key string) error {
	if err := d.client.Del(key).Err(); err != nil {
		return err
	}
	return nil
}

// GetRecord from the redis
func (d *SessionDB) getRecord(key string) (interface{}, error) {
	r := d.client.Get(key)

	if r.Err() != nil {

		return nil, r.Err()
	}
	return r.Val(), nil
}

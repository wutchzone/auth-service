package sessiondb

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/wutchzone/auth-service/pkg/accountdb"
	"time"

	"github.com/go-redis/redis"
)

var Singleton *SessionDB

type SessionDBConfiguration struct {
	Address string
	TableID int
}

type SessionDB struct {
	c *redis.Client
}

// NewSessionKey return unique identifier that can be used as a key for session DB
func NewSessionKey() string {
	return uuid.New().String()
}

// NewSessionItem return pointer to the newly created instance of the SessionItem
func NewSessionItem(role accountdb.Role) *SessionItem {
	return &SessionItem{
		Role: role,
	}
}

func GetInstance(conf SessionDBConfiguration) *SessionDB {
	if Singleton == nil {
		s := &SessionDB{}
		s.c = redis.NewClient(&redis.Options{
			Addr:     conf.Address,
			Password: "",           // no password set
			DB:       conf.TableID, // use default DB
		})

		// Check connection
		if _, err := s.c.Ping().Result(); err != nil {
			return nil
		}

		Singleton = s
	}

	return Singleton
}

// SetRecord to the redis, you can choose between persistent record (duration 0) and timed record (duration > 0)
func (d *SessionDB) SetRecord(key string, value SessionItem, time time.Duration) error {
	mar, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := d.c.Set(key, mar, time).Err(); err != nil {
		return err
	}
	return nil
}

// RemoveRecord from the redis
func (d *SessionDB) RemoveRecord(key string) error {
	if err := d.c.Del(key).Err(); err != nil {
		return err
	}
	return nil
}

// GetRecord from the redis
func (d *SessionDB) GetRecord(key string) (*SessionItem, error) {
	s := &SessionItem{}
	r := d.c.Get(key)

	if r.Err() != nil {
		return nil, errors.New("unable to get session record")
	}

	dat := r.Val()
	if err := json.Unmarshal([]byte(dat), &s); err != nil {
		return nil, err
	}

	return s, nil
}

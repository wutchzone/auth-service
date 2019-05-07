package sessiondb

import (
	"github.com/wutchzone/auth-service/pkg/accountdb"
	"testing"
	"time"
)

func TestNewSessionDB(t *testing.T) {
	d := GetInstance(SessionDBConfiguration{
		TableID: 1,
		Address: "localhost:6379",
	})

	if d == nil {
		t.Error("Failed to GetInstance of session DB. Expected value, got nil")
	}
	// Clean all data

	d.c.FlushAll()

	// Add record
	value := SessionItem{
		Role: accountdb.Guest,
	}
	key := "testKey"
	if err := d.SetRecord(key, value, 0); err != nil {
		t.Errorf("Error setting record, expected nil, got: %v", err)
		t.FailNow()
	}

	if got, err := d.GetRecord(key); err != nil || got.Role != value.Role {
		t.Errorf("Error getting value that was set. Expected: %v got %v ,with error %v", value, got, err)
	}

	// Remove record
	if err := d.RemoveRecord(key); err != nil {
		t.Errorf("Error removing key. Expected nil, got %v", err)
	}

	// Add timed record
	if err := d.SetRecord(key, value, time.Second*2); err != nil {
		t.Errorf("Error adding key key. Expected nil, got %v", err)
	}

	// Check for timed record before expire
	if got, err := d.GetRecord(key); err != nil || got.Role != value.Role {
		t.Errorf("Error getting error that was set. Expected: %v got %v with error %v", value, got, err)
	}

	// Check for timed record after expire
	time.Sleep(time.Second * 3)
	_, err := d.GetRecord(key)

	if err == nil {
		t.Errorf("User did not expired. Expected: error got nil")
	}
}

func TestNewSessionKey(t *testing.T) {
	u1 := NewSessionKey()
	u2 := NewSessionKey()

	if u1 == u2 {
		t.Error("Expected different values when calling NewSessionKey twice")
	}
}

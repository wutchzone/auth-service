package accountdb

import (
	"context"
	"testing"
)

func TestNewUser(t *testing.T) {
	if _, err := NewUser("test", "123456", "test@test.com", Guest); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if _, err := NewUser("", "123456", "test@test.com", Guest); err != nil {
		ex := "Name is not valid."
		if err.Error() != ex {
			t.Errorf("Expected %v, got %v", ex, err)
		}
	}

	if _, err := NewUser("test", "", "test@test.com", Guest); err != nil {
		ex := "Password is too short, 6 characters are required."
		if err.Error() != ex {
			t.Errorf("Expected %v, got %v", ex, err)
		}
	}

	if _, err := NewUser("test", "123456", "", Guest); err != nil {
		ex := "Email is not valid."
		if err.Error() != ex {
			t.Errorf("Expected %v, got %v", ex, err)
		}
	}
}

func TestUserConnection(t *testing.T) {
	// Make fresh instance of DB
	db := GetInstance(AccountConfiguration{
		AccoutCollectionName:  "accounts",
		Address:               "mongodb://localhost:27017",
		ServiceCollectionName: "services",
	})

	// Test connection
	if db == nil {
		t.Errorf("Expected value, got nil")
		t.FailNow()
	}

	db.db.Drop(context.Background())

	u1, _ := NewUser("test", "123456", "test@test.com", Guest)

	// Test if user not found
	if _, err := db.GetAccount(u1.Username); err == nil {
		t.Errorf("Expected error, got nil")
		t.Fail()
	}

	// Test create user
	if err := db.SaveUser(*u1); err != nil {
		t.Errorf("Error creating user. Expected nil, got %v", err)
		t.Fail()
	}
	ru1, err := db.GetAccount(u1.Username)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		t.Fail()
	}
	if ru1.Username != u1.Username && ru1.Email != u1.Email && ru1.Password != u1.Password && ru1.Role != u1.Role {
		t.Errorf("Returned user mismatched, expected the same user as saved")
		t.Fail()
	}

	// Test update user
	old := u1.Email
	u1.Email = "test2@test.com"
	if err := db.UpdateUser(*u1); err != nil {
		t.Errorf("Error user update. Expected nil, got %v", err)
		t.FailNow()
	}

	if ru2, err := db.GetAccount(u1.Username); err != nil {
		t.Errorf("Error user update. Expected nil, got %v", err)
		t.FailNow()
	} else if old == ru2.Email {
		t.Errorf("Error user update. Expected %v, got %v", old, u1.Email)
		t.FailNow()
	}

	// Test create user with same username
	if err := db.SaveUser(*u1); err == nil {
		t.Errorf("Error creating user. Username cannot be duplicit. Expected error, got nil")
		t.FailNow()
	}

	// Test delete user
	if err := db.DeleteUser(u1.Username); err != nil {
		t.Errorf("Expected nil, got %v", err)
		t.FailNow()
	}
	if _, err := db.GetAccount(u1.Username); err == nil {
		t.Errorf("Expected nil, got %v", err)
		t.FailNow()
	}
}

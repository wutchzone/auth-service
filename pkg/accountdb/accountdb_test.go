package accountdb

import "testing"

func TestNewUser(t *testing.T) {
	if _, err := NewUser("test", "123456", "test@test.com", "test"); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if _, err := NewUser("", "123456", "test@test.com", "test"); err != nil {
		ex := "Name is not valid."
		if err.Error() != ex {
			t.Errorf("Expected %v, got %v", ex, err)
		}
	}

	if _, err := NewUser("test", "", "test@test.com", "test"); err != nil {
		ex := "Password is too short, 6 characters are required."
		if err.Error() != ex {
			t.Errorf("Expected %v, got %v", ex, err)
		}
	}

	if _, err := NewUser("test", "123456", "", "test"); err != nil {
		ex := "Email is not valid."
		if err.Error() != ex {
			t.Errorf("Expected %v, got %v", ex, err)
		}
	}
}

func TestNewAccountDBConnection(t *testing.T) {
	db := NewAccountDBConnection("mongo:27017", "test", "test")
}

package accountdb

import (
	"testing"
)

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
	// Make fresh instance of DB
	// d, err := NewAccountDBConnection("mongodb://localhost:27017", "test", "test")
	// d.db.Drop(nil)

	// u1, _ := NewUser("test", "123456", "test@test.com", "test")
	// Test connection
	//if err != nil {
	//	t.Errorf("Expected nil, got %v", err)
	//	t.FailNow()
	//}

	//// Test if user not found
	//if err := d.GetAccount(u1.Username, &User{}); err == nil {
	//	t.Errorf("Expected error, got nil")
	//	t.FailNow()
	//}
	//
	//// Test create user
	//if err := d.SaveAccount(u1); err != nil {
	//	t.Errorf("Error creating user. Expected nil, got %v", err)
	//	t.FailNow()
	//}
	//if err := d.GetAccount(u1.Username, &User{}); err != nil {
	//	t.Errorf("Expected nil, got %v", err)
	//	t.FailNow()
	//}
	//
	//// Test update user
	//old := u1.Email
	//u1.Email = "test2@test.com"
	//if err := d.UpdateAccount(u1); err != nil {
	//	t.Errorf("Error user update. Expected nil, got %v", err)
	//	t.FailNow()
	//}
	//
	//u2 := User{}
	//if err := d.GetAccount(u1.Username, &u2); err != nil {
	//	t.Errorf("Error user update. Expected nil, got %v", err)
	//	t.FailNow()
	//} else if old == u2.Email {
	//	t.Errorf("Error user update. Expected %v, got %v", old, u1.Email)
	//	t.FailNow()
	//}
	//
	//// Test create user with same username
	////if err := d.saveAccount(u1); err == nil {
	////	t.Errorf("Error creating user. Username cannot be duplicit. Expected error, got nil")
	////	t.FailNow()
	////}
	//
	//// Test delete user
	//if err := d.DeleteAccount(u1); err != nil {
	//	t.Errorf("Expected nil, got %v", err)
	//	t.FailNow()
	//}
	//if err := d.GetAccount(u1.Username, &User{}); err == nil {
	//	t.Errorf("Expected nil, got %v", err)
	//	t.FailNow()
	//}

	// Test get all users
	// usrs := make([]user.User, 10)

	// cursor := d.GetAll()

	//itemsToSend := []user.User{}
	//itemRead := user.User{}
	//for cursor.Next(nil) {
	//	err := cursor.Decode(&itemRead)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	itemsToSend = append(itemsToSend, itemRead)
	//}

}

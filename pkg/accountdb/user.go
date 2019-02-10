package accountdb

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	Username string `bson:"name",json:"name"`
	Password string `bson:"password",json:"-"`
	Email    string `bson:"email",json:"email"`
	Role     string `bson:"role",json:"role"`
}

// NewUser instance
func NewUser(name string, password string, email string, role string) (*User, error) {
	hashedPswd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if name == "" {
		return nil, errors.New("Name is not valid.")
	}
	if len(password) < 6 {
		return nil, errors.New("Password is too short, 6 characters are required.")
	}
	if email == "" {
		return nil, errors.New("Email is not valid.")
	}

	return &User{
		Username: name,
		Password: string(hashedPswd),
		Email:    email,
		Role:     role,
	}, nil
}

// ComparePswdAndHash compares if given password is the crypted one
func (u *User) ComparePswdAndHash(password string) error {
	return bcrypt.CompareHashAndPassword([]byte((*u).Password), []byte(password))
}

package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	Username string
	Password string
	Email    string
	Role     Role
}

// Role enum for available roles
type Role int

const (
	Default      Role = 0
	DefaultUser  Role = 10
	ModerateUser Role = 20
	Admin        Role = 30
	SuperAdmin   Role = 40
)

// NewUser instance
func NewUser(name string, password string, email string, role Role) (*User, error) {
	hashedPswd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if name == "" {
		return nil, errors.New("Name is too short")
	}
	if password == "" {
		return nil, errors.New("Password is too short")
	}
	if email == "" {
		return nil, errors.New("Email is too short")
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

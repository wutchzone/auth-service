package parser

import (
	"github.com/wutchzone/auth-service/pkg/user"
)

type UserJSON struct {
	Name     string
	Password string
	Email    string
	Role     user.Role
}

package parser

import (
	"github.com/wutchzone/auth-service/pkg/user"
)

type API struct {
	From        string
	To          string
	Privilegies user.Role
}

type DB struct {
	Address  string
	Password string
}

type ConfigJSON struct {
	API       []API
	Port      int
	Sessiondb DB
	Userdb    DB
}

package parser

import (
	"github.com/wutchzone/auth-service/pkg/user"
)

type API struct {
	From        string
	To          string
	Privilegies user.Role
}

type ConfigJSON struct {
	API []API
}

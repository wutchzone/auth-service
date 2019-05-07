package sessiondb

import "github.com/wutchzone/auth-service/pkg/accountdb"

type SessionItem struct {
	Role accountdb.Role
}

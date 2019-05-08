package decoder

import "time"

type Authorize struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

package decoder

import "time"

type Authorize struct {
	Token      string
	Expiration time.Time
}

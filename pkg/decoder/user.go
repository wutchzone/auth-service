package decoder

type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

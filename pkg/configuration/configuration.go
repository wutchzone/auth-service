package configuration

// Route definition
type Route struct {
	From  string
	To    string
	Level int
}

// Role definition
type Role struct {
	Name  string // Assigned name
	Level int    // Access level
}

type DB struct {
	URL      string // URL to the connected DB
	Name     string // Name for connecting to the DB
	Password string // Password for connection to the DB
	Table    string // Used table name
}

type Configuration struct {
	DefaultPort int                        // Port where service will be listening
	Routes      []Route `json:"routes"`    // Routes definition
	Roles       []Role  `json:"roles"`     // Roles definition
	SessionDB   DB      `json:"sessiondb"` // DB for storing sessiondb information (Only Redis is currently supported)
	AccountDB   DB      `json:"accountdb"`   // Parameters for connecting to the DB where user data's are stored (Only MongoDB is currently supported)
}

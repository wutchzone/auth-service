package configuration

type StartupConfiguration struct {
	ID        string `bson:"name"`
	FirstBoot bool
}

func NewStartup() *StartupConfiguration {
	return &StartupConfiguration{
		ID: "startup",
	}
}

func (s *StartupConfiguration) Name() string {
	return s.ID
}

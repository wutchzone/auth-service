package accountdb

import "github.com/google/uuid"

type Service struct {
	ID    string `bson:"name"`
	Level int    `bson:"level"`
}

// NewService created pointer to the newly created instance
func NewService(level int) *Service {
	return &Service{
		ID:    uuid.New().String(),
		Level: level,
	}
}

func (s Service) Name() string {
	return s.ID
}

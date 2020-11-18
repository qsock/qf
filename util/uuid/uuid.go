package uuid

import (
	"github.com/google/uuid"
)

type UUID = uuid.UUID

func New() UUID {

	return uuid.New()
}

func NewString() string {
	return New().String()
}

func NewURN() string {
	return New().URN()
}

func NewUUID() (UUID, error) {
	return uuid.NewUUID()
}

func NewMD5(space UUID, data []byte) UUID {
	return uuid.NewMD5(space, data)
}

func NewRandom() (UUID, error) {
	return uuid.NewRandom()
}

func NewSHA1(space UUID, data []byte) UUID {
	return uuid.NewSHA1(space, data)
}

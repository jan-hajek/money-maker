package position

import (
	"github.com/satori/go.uuid"
	"strconv"
)

type IdGenerator interface {
	Generate() string
}

type UuidGenerator struct {
}

func (s *UuidGenerator) Generate() string {
	return uuid.Must(uuid.NewV4()).String()
}

type IncrementGenerator struct {
	id int
}

func (s *IncrementGenerator) Generate() string {
	s.id++
	return strconv.Itoa(s.id)
}

package adapter

import (
	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/oklog/ulid/v2"
)

var _ (port.IDGenerator) = (*ULIDGenerator)(nil)

type ULIDGenerator struct{}

func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{}
}

func (g *ULIDGenerator) Generate() (entity.ID, error) {
	return entity.ParseID(ulid.Make().String())
}

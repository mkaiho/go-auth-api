package crypto

import (
	"context"

	"github.com/mkaiho/go-auth-api/adapter"
	"golang.org/x/crypto/bcrypt"
)

var _ adapter.HashGenerator = (*BcryptoHashGenerator)(nil)

type BcryptoHashGenerator struct {
	cost int
}

func NewBcryptoHashGenerator() *BcryptoHashGenerator {
	return &BcryptoHashGenerator{
		cost: bcrypt.DefaultCost,
	}
}

func (g *BcryptoHashGenerator) Generate(ctx context.Context, value []byte) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), g.cost)
	if err != nil {
		return nil, err
	}
	return hashed, nil
}

func (g *BcryptoHashGenerator) Compare(ctx context.Context, hashed []byte, value []byte) error {
	err := bcrypt.CompareHashAndPassword(hashed, value)
	if err != nil {
		return err
	}
	return nil
}

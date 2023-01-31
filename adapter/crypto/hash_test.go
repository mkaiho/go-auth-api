package crypto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptoHashGenerator_Generate(t *testing.T) {
	const cost = bcrypt.DefaultCost
	const wantLength = 60

	type args struct {
		ctx   context.Context
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "return hashed value",
			args: args{
				ctx:   context.Background(),
				value: []byte("hello world"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &BcryptoHashGenerator{
				cost: cost,
			}
			got, err := g.Generate(tt.args.ctx, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("BcryptoHashGenerator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var gotLength = len(got)
				assert.Equal(t, wantLength, gotLength, "BcryptoHashGenerator.Generate() gotLength = %v, wantLength %v", gotLength, wantLength)
			}
		})
	}
}

func TestBcryptoHashGenerator_Compare(t *testing.T) {
	type fields struct {
		cost int
	}
	type args struct {
		ctx    context.Context
		hashed []byte
		value  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "return no error",
			args: args{
				ctx:    context.Background(),
				hashed: []byte("$2a$10$6TJH9Dhk9tHbR57kbC1ZaOC4gJEQ1lLO.kI5gdeEwgB7REWGxayoC"),
				value:  []byte("hello world"),
			},
			wantErr: false,
		},
		{
			name: "return error when compared values do not match",
			args: args{
				ctx:    context.Background(),
				hashed: []byte("$2a$10$6TJH9Dhk9tHbR57kbC1ZaOC4gJEQ1lLO.kI5gdeEwgB7REWGxayoC"),
				value:  []byte("hello world!"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &BcryptoHashGenerator{
				cost: tt.fields.cost,
			}
			if err := g.Compare(tt.args.ctx, tt.args.hashed, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("BcryptoHashGenerator.Compare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

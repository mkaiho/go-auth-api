package id

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestULIDGenerator_Generate(t *testing.T) {
	const wantLength = 26
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "return generated id",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &ULIDGenerator{}
			got, err := g.Generate(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ULIDGenerator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var gotLength = len(got)
				assert.Equal(t, wantLength, gotLength, "ULIDGenerator.Generate() gotLength = %v, wantLength %v", gotLength, wantLength)
			}
		})
	}
}

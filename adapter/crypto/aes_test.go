package crypto

import (
	"context"
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESCryptor_Encrypt(t *testing.T) {
	const key = "abcdefghijklmnopqrstuvwxyz012345"
	const blockSize = aes.BlockSize
	type args struct {
		ctx       context.Context
		plaintext []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "return ciphertext",
			args: args{
				ctx:       context.Background(),
				plaintext: []byte("hello world"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AESCryptor{
				key:       key,
				blockSize: blockSize,
			}
			got, err := c.Encrypt(tt.args.ctx, tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESCryptor.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var wantLength int = c.blockSize + len(tt.args.plaintext)
				var gotLength = len(got)
				assert.Equal(t, wantLength, gotLength, "AESCryptor.Encrypt() gotLength = %v, wantLength %v", gotLength, wantLength)
			}
		})
	}
}

func TestAESCryptor_Decrypt(t *testing.T) {
	const key = "abcdefghijklmnopqrstuvwxyz012345"
	const blockSize = aes.BlockSize
	type args struct {
		ctx        context.Context
		ciphertext []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "return plaintext",
			args: args{
				ctx:        context.Background(),
				ciphertext: []byte("\xa8L\xfe\x91\xedl\x18DZY8\x96ni\x81\x9a\xf4\xe0\xac0\xd4\xeb\xa3I/b\x0e"),
			},
			want:    []byte("hello world"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AESCryptor{
				key:       key,
				blockSize: blockSize,
			}
			got, err := c.Decrypt(tt.args.ctx, tt.args.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESCryptor.Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "AESCryptor.Decrypt() = %v, want %v", got, tt.want)
		})
	}
}

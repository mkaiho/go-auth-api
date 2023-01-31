package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/mkaiho/go-auth-api/adapter"
)

var _ adapter.Cryptor = (*AESCryptor)(nil)

type AESCryptor struct {
	key       string
	blockSize int
}

func NewAESCryptor(key string) *AESCryptor {
	return &AESCryptor{
		key:       key,
		blockSize: aes.BlockSize,
	}
}

func (c *AESCryptor) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	iv := make([]byte, c.blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	enc := cipher.NewCFBEncrypter(block, iv)
	enc.XORKeyStream(ciphertext, plaintext)

	return append(iv, ciphertext...), nil
}

func (c *AESCryptor) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	iv := []byte(ciphertext[:c.blockSize])
	text := []byte(ciphertext[c.blockSize:])

	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(text))
	enc := cipher.NewCFBDecrypter(block, iv)
	enc.XORKeyStream(plaintext, []byte(text))

	return plaintext, nil
}

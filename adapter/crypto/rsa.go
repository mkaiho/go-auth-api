package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"strings"
)

var (
	ErrInvalidRSAPrivateKeyFormat    = errors.New("indalid RSA private key format")
	ErrInvalidRSAPrivateKeyBlockType = errors.New("invalid block type")
)

type RSAPrivateKeyFormat int

const (
	RSAPrivateKeyFormatUnsupported RSAPrivateKeyFormat = iota
	RSAPrivateKeyFormatDer
	RSAPrivateKeyFormatPem
)

func (f RSAPrivateKeyFormat) String() string {
	return [...]string{
		"unsupported",
		"der",
		"pem",
	}[f]
}

func ParseRSAPrivateKeyFormat(v string) RSAPrivateKeyFormat {
	s := strings.ToLower(strings.ReplaceAll(v, "_", ""))
	switch s {
	default:
		return RSAPrivateKeyFormatUnsupported
	case "der":
		return RSAPrivateKeyFormatDer
	case "pem":
		return RSAPrivateKeyFormatPem
	}
}

type RSAPrivateKeyBlockType int

const (
	RSAPrivateKeyKeyTypeUnknown RSAPrivateKeyBlockType = iota
	RSAPrivateKeyBlockTypePKCS1
	RSAPrivateKeyBlockTypePKCS8
)

func (f RSAPrivateKeyBlockType) String() string {
	return [...]string{
		"",
		"RSA PRIVATE KEY",
		"PRIVATE KEY",
	}[f]
}

func ParseRSAPrivateKeyBlockType(v string) RSAPrivateKeyBlockType {
	switch v {
	default:
		return RSAPrivateKeyKeyTypeUnknown
	case "RSA PRIVATE KEY":
		return RSAPrivateKeyBlockTypePKCS1
	case "PRIVATE KEY":
		return RSAPrivateKeyBlockTypePKCS8
	}
}

var _ RSAKeyManager = (*rsaKeyManager)(nil)

type RSAKeyManager interface {
	GenerateRSAPrivateKey(bits int) (*rsa.PrivateKey, error)
	ReadPemFile(filename string) (*rsa.PrivateKey, error)
	ReadPemBytes(b []byte) (*rsa.PrivateKey, error)
	ConvertFormat(privateKey *rsa.PrivateKey, format RSAPrivateKeyFormat) ([]byte, error)
}

func NewRSAKeyManager() *rsaKeyManager {
	return &rsaKeyManager{
		random: rand.Reader,
	}
}

type rsaKeyManager struct {
	random io.Reader
}

func (m *rsaKeyManager) GenerateRSAPrivateKey(bits int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (m *rsaKeyManager) ReadPemFile(filename string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return m.ReadPemBytes(b)
}

func (m *rsaKeyManager) ReadPemBytes(b []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(b)

	if block == nil {
		key, err := x509.ParsePKCS1PrivateKey(b)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	if block.Type == RSAPrivateKeyBlockTypePKCS1.String() {
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	if block.Type == RSAPrivateKeyBlockTypePKCS8.String() {
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		var ok bool
		key, ok := keyInterface.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not RSA private key")
		}
		return key, nil
	}

	return nil, ErrInvalidRSAPrivateKeyBlockType
}

func (m *rsaKeyManager) ConvertFormat(privateKey *rsa.PrivateKey, format RSAPrivateKeyFormat) ([]byte, error) {
	der := x509.MarshalPKCS1PrivateKey(privateKey)
	if format == RSAPrivateKeyFormatDer {
		return der, nil
	}
	if format == RSAPrivateKeyFormatPem {
		var w bytes.Buffer
		if err := pem.Encode(&w, &pem.Block{Type: RSAPrivateKeyBlockTypePKCS1.String(), Bytes: der}); err != nil {
			return nil, err
		}
		b, err := io.ReadAll(&w)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return nil, ErrInvalidRSAPrivateKeyFormat
}

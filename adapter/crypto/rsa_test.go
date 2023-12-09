package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSAPrivateKeyFormat_String(t *testing.T) {
	tests := []struct {
		name string
		f    RSAPrivateKeyFormat
		want string
	}{
		{
			name: `return "der"`,
			f:    RSAPrivateKeyFormatDer,
			want: "der",
		},
		{
			name: `return "pem"`,
			f:    RSAPrivateKeyFormatPem,
			want: "pem",
		},
		{
			name: `return "unsupported"`,
			f:    RSAPrivateKeyFormatUnsupported,
			want: "unsupported",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.String()
			assert.Equalf(t, tt.want, got, "RSAPrivateKeyFormat.String() = %v, want %v", got, tt.want)
		})
	}
}

func TestParseRSAPrivateKeyFormat(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want RSAPrivateKeyFormat
	}{
		{
			name: `return enum of der format when value is "der"`,
			args: args{
				v: "der",
			},
			want: RSAPrivateKeyFormatDer,
		},
		{
			name: `return enum of der format when value is "pem"`,
			args: args{
				v: "pem",
			},
			want: RSAPrivateKeyFormatPem,
		},
		{
			name: `return enum of unsupported when value is invalid`,
			args: args{
				v: "invalid",
			},
			want: RSAPrivateKeyFormatUnsupported,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseRSAPrivateKeyFormat(tt.args.v)
			assert.Equalf(t, tt.want, got, "ParseRSAPrivateKeyFormat() = %v, want %v", got, tt.want)
		})
	}
}

func TestRSAPrivateKeyBlockType_String(t *testing.T) {
	tests := []struct {
		name string
		f    RSAPrivateKeyBlockType
		want string
	}{
		{
			name: `return "RSA PRIVATE KEY"`,
			f:    RSAPrivateKeyBlockTypePKCS1,
			want: "RSA PRIVATE KEY",
		},
		{
			name: `return "PRIVATE KEY"`,
			f:    RSAPrivateKeyBlockTypePKCS8,
			want: "PRIVATE KEY",
		},
		{
			name: `return empty`,
			f:    RSAPrivateKeyKeyTypeUnknown,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.String()
			assert.Equalf(t, tt.want, got, "RSAPrivateKeyBlockType.String() = %v, want %v", got, tt.want)
		})
	}
}

func TestParseRSAPrivateKeyBlockType(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want RSAPrivateKeyBlockType
	}{
		{
			name: `return enum of PKCS #1 format when value is "RSA PRIVATE KEY"`,
			args: args{
				v: "RSA PRIVATE KEY",
			},
			want: RSAPrivateKeyBlockTypePKCS1,
		},
		{
			name: `return enum of PKCS #8 format when value is "PRIVATE KEY"`,
			args: args{
				v: "PRIVATE KEY",
			},
			want: RSAPrivateKeyBlockTypePKCS8,
		},
		{
			name: `return enum of unknown when value is invalid`,
			args: args{
				v: "invalid",
			},
			want: RSAPrivateKeyKeyTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseRSAPrivateKeyBlockType(tt.args.v)
			assert.Equalf(t, tt.want, got, "ParseRSAPrivateKeyBlockType() = %v, want %v", got, tt.want)
		})
	}
}

func Test_rsaKeyManager_GenerateRSAPrivateKey(t *testing.T) {
	type fields struct {
		random io.Reader
	}
	type args struct {
		bits int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "return non empty value",
			fields: fields{
				random: strings.NewReader("a"),
			},
			args: args{
				bits: 2048,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &rsaKeyManager{
				random: tt.fields.random,
			}
			got, err := m.GenerateRSAPrivateKey(tt.args.bits)
			tt.assertion(t, err)
			if err == nil {
				assert.NotEmpty(t, got)
			}
		})
	}
}

func Test_rsaKeyManager_ReadPemFile(t *testing.T) {
	testPemFileName := "test.pem"
	type fields struct {
		random io.Reader
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "return rsa private key",
			args: args{
				filename: "test.pem",
			},
			fields: fields{
				random: strings.NewReader("test"),
			},
			assertion: assert.NoError,
		},
		{
			name: "return error",
			args: args{
				filename: "not_found.pem",
			},
			fields: fields{
				random: strings.NewReader("test"),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &rsaKeyManager{
				random: tt.fields.random,
			}

			dest, err := os.MkdirTemp("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dest)

			k, err := m.GenerateRSAPrivateKey(2048)
			if err != nil {
				t.Fatal(err)
			}
			b, err := m.ConvertFormat(k, RSAPrivateKeyFormatPem)
			if err != nil {
				t.Fatal(err)
			}
			pemFilePath := filepath.Join(dest, testPemFileName)
			if err := os.WriteFile(pemFilePath, b, 0600); err != nil {
				t.Fatal(err)
			}

			got, err := m.ReadPemFile(filepath.Join(dest, tt.args.filename))
			tt.assertion(t, err)
			if err == nil {
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_rsaKeyManager_ReadPemBytes(t *testing.T) {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	var pemBuf bytes.Buffer
	if err := pem.Encode(&pemBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}); err != nil {
		panic(err)
	}
	pemBytes, err := io.ReadAll(&pemBuf)
	if err != nil {
		panic(err)
	}

	type fields struct {
		random io.Reader
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "return rsa private key",
			fields: fields{
				random: strings.NewReader("test"),
			},
			args: args{
				b: pemBytes,
			},
			assertion: assert.NoError,
		},
		{
			name: "return error when value is invalid",
			fields: fields{
				random: strings.NewReader("test"),
			},
			args: args{
				b: []byte("invalid"),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &rsaKeyManager{
				random: tt.fields.random,
			}
			got, err := m.ReadPemBytes(tt.args.b)
			tt.assertion(t, err)
			if err == nil {
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_rsaKeyManager_ConvertFormat(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	type fields struct {
		random io.Reader
	}
	type args struct {
		privateKey *rsa.PrivateKey
		format     RSAPrivateKeyFormat
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "return bytes parsed as der format",
			fields: fields{
				random: strings.NewReader("test"),
			},
			args: args{
				privateKey: privateKey,
				format:     RSAPrivateKeyFormatDer,
			},
			assertion: assert.NoError,
		},
		{
			name: "return bytes parsed as pem format",
			fields: fields{
				random: strings.NewReader("test"),
			},
			args: args{
				privateKey: privateKey,
				format:     RSAPrivateKeyFormatPem,
			},
			assertion: assert.NoError,
		},
		{
			name: "return an error if the specified converted format is unsupported",
			fields: fields{
				random: strings.NewReader("test"),
			},
			args: args{
				privateKey: privateKey,
				format:     RSAPrivateKeyFormatUnsupported,
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrInvalidRSAPrivateKeyFormat)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &rsaKeyManager{
				random: tt.fields.random,
			}
			got, err := m.ConvertFormat(tt.args.privateKey, tt.args.format)
			tt.assertion(t, err)
			if err == nil {
				assert.NotNil(t, got)
			}
		})
	}
}

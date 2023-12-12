package adapter

import (
	"context"
	"crypto/rsa"
	"io"

	"github.com/mkaiho/go-auth-api/adapter/crypto"
	"github.com/mkaiho/go-auth-api/adapter/storage"
)

type KeyAccess struct {
	storageClient storage.Client
	rsaKeyManager crypto.RSAKeyManager
}

func NewKeyAccess(storageClient storage.Client, rsaKeyManager crypto.RSAKeyManager) *KeyAccess {
	return &KeyAccess{
		storageClient: storageClient,
		rsaKeyManager: rsaKeyManager,
	}
}

func (a *KeyAccess) ReadPrivateKey(ctx context.Context, path string) (*rsa.PrivateKey, error) {
	r, err := a.storageClient.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	key, err := a.rsaKeyManager.ReadPemBytes(b)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (a *KeyAccess) Save(ctx context.Context, path string, key *rsa.PrivateKey) error {
	b, err := a.rsaKeyManager.ConvertFormat(key, crypto.RSAPrivateKeyFormatPem)
	if err != nil {
		return err
	}
	err = a.storageClient.Save(ctx, path, storage.MimeTypePem, b)
	if err != nil {
		return err
	}

	return nil
}

func (a *KeyAccess) Remove(ctx context.Context, path string) error {
	err := a.storageClient.Remove(ctx, path)
	if err != nil {
		return err
	}

	return nil
}

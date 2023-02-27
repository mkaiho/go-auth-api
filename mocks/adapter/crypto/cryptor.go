// Code generated by mockery v2.19.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Cryptor is an autogenerated mock type for the Cryptor type
type Cryptor struct {
	mock.Mock
}

// Decrypt provides a mock function with given fields: ctx, ciphertext
func (_m *Cryptor) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	ret := _m.Called(ctx, ciphertext)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, []byte) []byte); ok {
		r0 = rf(ctx, ciphertext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctx, ciphertext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encrypt provides a mock function with given fields: ctx, plaintext
func (_m *Cryptor) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	ret := _m.Called(ctx, plaintext)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, []byte) []byte); ok {
		r0 = rf(ctx, plaintext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []byte) error); ok {
		r1 = rf(ctx, plaintext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCryptor interface {
	mock.TestingT
	Cleanup(func())
}

// NewCryptor creates a new instance of Cryptor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCryptor(t mockConstructorTestingTNewCryptor) *Cryptor {
	mock := &Cryptor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
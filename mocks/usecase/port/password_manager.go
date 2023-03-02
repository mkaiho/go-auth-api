// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/mkaiho/go-auth-api/entity"
	mock "github.com/stretchr/testify/mock"
)

// PasswordManager is an autogenerated mock type for the PasswordManager type
type PasswordManager struct {
	mock.Mock
}

// Compare provides a mock function with given fields: ctx, hashedPassword, password
func (_m *PasswordManager) Compare(ctx context.Context, hashedPassword entity.HashedPassword, password entity.Password) error {
	ret := _m.Called(ctx, hashedPassword, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.HashedPassword, entity.Password) error); ok {
		r0 = rf(ctx, hashedPassword, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Hash provides a mock function with given fields: ctx, value
func (_m *PasswordManager) Hash(ctx context.Context, value string) (entity.HashedPassword, error) {
	ret := _m.Called(ctx, value)

	var r0 entity.HashedPassword
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (entity.HashedPassword, error)); ok {
		return rf(ctx, value)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.HashedPassword); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Get(0).(entity.HashedPassword)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPasswordManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewPasswordManager creates a new instance of PasswordManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPasswordManager(t mockConstructorTestingTNewPasswordManager) *PasswordManager {
	mock := &PasswordManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

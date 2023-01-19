// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/mkaiho/go-auth-api/entity"
	interactor "github.com/mkaiho/go-auth-api/usecase/interactor"

	mock "github.com/stretchr/testify/mock"
)

// UserInteractor is an autogenerated mock type for the UserInteractor type
type UserInteractor struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, input
func (_m *UserInteractor) CreateUser(ctx context.Context, input interactor.CreateUserInput) (*entity.User, error) {
	ret := _m.Called(ctx, input)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, interactor.CreateUserInput) *entity.User); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interactor.CreateUserInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, input
func (_m *UserInteractor) GetUser(ctx context.Context, input interactor.GetUserInput) (*entity.User, error) {
	ret := _m.Called(ctx, input)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, interactor.GetUserInput) *entity.User); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interactor.GetUserInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserInteractor interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserInteractor creates a new instance of UserInteractor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserInteractor(t mockConstructorTestingTNewUserInteractor) *UserInteractor {
	mock := &UserInteractor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

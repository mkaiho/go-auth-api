// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	entity "github.com/mkaiho/go-auth-api/entity"
	mock "github.com/stretchr/testify/mock"
)

// IDGenerator is an autogenerated mock type for the IDGenerator type
type IDGenerator struct {
	mock.Mock
}

// Generate provides a mock function with given fields:
func (_m *IDGenerator) Generate() (entity.ID, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Generate")
	}

	var r0 entity.ID
	var r1 error
	if rf, ok := ret.Get(0).(func() (entity.ID, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() entity.ID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(entity.ID)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIDGenerator creates a new instance of IDGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDGenerator {
	mock := &IDGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

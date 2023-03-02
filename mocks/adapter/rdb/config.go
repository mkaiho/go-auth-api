// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	rdb "github.com/mkaiho/go-auth-api/adapter/rdb"
	mock "github.com/stretchr/testify/mock"
)

// Config is an autogenerated mock type for the Config type
type Config struct {
	mock.Mock
}

// GetDSN provides a mock function with given fields:
func (_m *Config) GetDSN() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetDriverName provides a mock function with given fields:
func (_m *Config) GetDriverName() rdb.DriverName {
	ret := _m.Called()

	var r0 rdb.DriverName
	if rf, ok := ret.Get(0).(func() rdb.DriverName); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(rdb.DriverName)
	}

	return r0
}

// GetMaxConns provides a mock function with given fields:
func (_m *Config) GetMaxConns() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

type mockConstructorTestingTNewConfig interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfig creates a new instance of Config. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfig(t mockConstructorTestingTNewConfig) *Config {
	mock := &Config{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

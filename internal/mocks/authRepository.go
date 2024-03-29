// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/hinccvi/go-ddd/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// AuthRepository is an autogenerated mock type for the Repository type
type AuthRepository struct {
	mock.Mock
}

// GetUserByUsername provides a mock function with given fields: ctx, username
func (_m *AuthRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	ret := _m.Called(ctx, username)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.User); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthRepository creates a new instance of AuthRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthRepository(t mockConstructorTestingTNewAuthRepository) *AuthRepository {
	mock := &AuthRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

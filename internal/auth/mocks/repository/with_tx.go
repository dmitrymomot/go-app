// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks_repository

import (
	context "context"

	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// TxQuerier is an autogenerated mock type for the TxQuerier type
type TxQuerier struct {
	mock.Mock
}

// BeginTx provides a mock function with given fields: ctx
func (_m *TxQuerier) BeginTx(ctx context.Context) (auth_repository.TxQuerier, error) {
	ret := _m.Called(ctx)

	var r0 auth_repository.TxQuerier
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (auth_repository.TxQuerier, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) auth_repository.TxQuerier); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(auth_repository.TxQuerier)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CleanUpTokens provides a mock function with given fields: ctx
func (_m *TxQuerier) CleanUpTokens(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CleanUpVerifications provides a mock function with given fields: ctx
func (_m *TxQuerier) CleanUpVerifications(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Commit provides a mock function with given fields:
func (_m *TxQuerier) Commit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: ctx, email
func (_m *TxQuerier) CreateUser(ctx context.Context, email string) (uuid.UUID, error) {
	ret := _m.Called(ctx, email)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (uuid.UUID, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) uuid.UUID); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteTokenByAccessTokenID provides a mock function with given fields: ctx, accessTokenID
func (_m *TxQuerier) DeleteTokenByAccessTokenID(ctx context.Context, accessTokenID uuid.UUID) error {
	ret := _m.Called(ctx, accessTokenID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, accessTokenID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTokenByRefreshTokenID provides a mock function with given fields: ctx, refreshTokenID
func (_m *TxQuerier) DeleteTokenByRefreshTokenID(ctx context.Context, refreshTokenID uuid.UUID) error {
	ret := _m.Called(ctx, refreshTokenID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, refreshTokenID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTokensByUserID provides a mock function with given fields: ctx, userID
func (_m *TxQuerier) DeleteTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUserByID provides a mock function with given fields: ctx, id
func (_m *TxQuerier) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteVerificationByID provides a mock function with given fields: ctx, id
func (_m *TxQuerier) DeleteVerificationByID(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindTokenByAccessTokenID provides a mock function with given fields: ctx, accessTokenID
func (_m *TxQuerier) FindTokenByAccessTokenID(ctx context.Context, accessTokenID uuid.UUID) (auth_repository.Token, error) {
	ret := _m.Called(ctx, accessTokenID)

	var r0 auth_repository.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (auth_repository.Token, error)); ok {
		return rf(ctx, accessTokenID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) auth_repository.Token); ok {
		r0 = rf(ctx, accessTokenID)
	} else {
		r0 = ret.Get(0).(auth_repository.Token)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, accessTokenID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTokenByID provides a mock function with given fields: ctx, id
func (_m *TxQuerier) FindTokenByID(ctx context.Context, id uuid.UUID) (auth_repository.Token, error) {
	ret := _m.Called(ctx, id)

	var r0 auth_repository.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (auth_repository.Token, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) auth_repository.Token); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(auth_repository.Token)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTokenByRefreshTokenID provides a mock function with given fields: ctx, refreshTokenID
func (_m *TxQuerier) FindTokenByRefreshTokenID(ctx context.Context, refreshTokenID uuid.UUID) (auth_repository.Token, error) {
	ret := _m.Called(ctx, refreshTokenID)

	var r0 auth_repository.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (auth_repository.Token, error)); ok {
		return rf(ctx, refreshTokenID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) auth_repository.Token); ok {
		r0 = rf(ctx, refreshTokenID)
	} else {
		r0 = ret.Get(0).(auth_repository.Token)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, refreshTokenID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserByEmail provides a mock function with given fields: ctx, email
func (_m *TxQuerier) FindUserByEmail(ctx context.Context, email string) (auth_repository.User, error) {
	ret := _m.Called(ctx, email)

	var r0 auth_repository.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (auth_repository.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) auth_repository.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(auth_repository.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserByID provides a mock function with given fields: ctx, id
func (_m *TxQuerier) FindUserByID(ctx context.Context, id uuid.UUID) (auth_repository.User, error) {
	ret := _m.Called(ctx, id)

	var r0 auth_repository.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (auth_repository.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) auth_repository.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(auth_repository.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindVerificationByID provides a mock function with given fields: ctx, id
func (_m *TxQuerier) FindVerificationByID(ctx context.Context, id uuid.UUID) (auth_repository.Verification, error) {
	ret := _m.Called(ctx, id)

	var r0 auth_repository.Verification
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (auth_repository.Verification, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) auth_repository.Verification); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(auth_repository.Verification)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefreshToken provides a mock function with given fields: ctx, arg
func (_m *TxQuerier) RefreshToken(ctx context.Context, arg auth_repository.RefreshTokenParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth_repository.RefreshTokenParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Rollback provides a mock function with given fields:
func (_m *TxQuerier) Rollback() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreOrUpdateVerification provides a mock function with given fields: ctx, arg
func (_m *TxQuerier) StoreOrUpdateVerification(ctx context.Context, arg auth_repository.StoreOrUpdateVerificationParams) (uuid.UUID, error) {
	ret := _m.Called(ctx, arg)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, auth_repository.StoreOrUpdateVerificationParams) (uuid.UUID, error)); ok {
		return rf(ctx, arg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auth_repository.StoreOrUpdateVerificationParams) uuid.UUID); ok {
		r0 = rf(ctx, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, auth_repository.StoreOrUpdateVerificationParams) error); ok {
		r1 = rf(ctx, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreToken provides a mock function with given fields: ctx, arg
func (_m *TxQuerier) StoreToken(ctx context.Context, arg auth_repository.StoreTokenParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth_repository.StoreTokenParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUserEmailByID provides a mock function with given fields: ctx, arg
func (_m *TxQuerier) UpdateUserEmailByID(ctx context.Context, arg auth_repository.UpdateUserEmailByIDParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth_repository.UpdateUserEmailByIDParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUserVerificationStatusByID provides a mock function with given fields: ctx, arg
func (_m *TxQuerier) UpdateUserVerificationStatusByID(ctx context.Context, arg auth_repository.UpdateUserVerificationStatusByIDParams) error {
	ret := _m.Called(ctx, arg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth_repository.UpdateUserVerificationStatusByIDParams) error); ok {
		r0 = rf(ctx, arg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTxQuerier interface {
	mock.TestingT
	Cleanup(func())
}

// NewTxQuerier creates a new instance of TxQuerier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTxQuerier(t mockConstructorTestingTNewTxQuerier) *TxQuerier {
	mock := &TxQuerier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

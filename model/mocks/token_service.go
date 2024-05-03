package mocks

import (
	"context"

	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"
	"github.com/stretchr/testify/mock"
)

// MockTokenService is a mock type for model.TokenService
type MockTokenService struct {
	mock.Mock
}

// NewPairForUser mocks concrete NewPairForUser
func (m *MockTokenService) NewPairForUser(ctx context.Context, u *model.Users, prevTokenID string) (*model.TokenPair, *errors.MathSheetsError) {
	ret := m.Called(ctx, u, prevTokenID)

	// first value passed to "Return"
	var r0 *model.TokenPair
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*model.TokenPair)
	}

	var r1 *errors.MathSheetsError

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*errors.MathSheetsError)
	}

	return r0, r1
}

// Signout mocks concrete Signout
func (m *MockTokenService) Signout(ctx context.Context, uid string, tokenString string) *errors.MathSheetsError {
	ret := m.Called(ctx, uid, tokenString)
	var r0 *errors.MathSheetsError

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
	}

	return r0
}

// ValidateIDToken mocks concrete ValidateIDToken
func (m *MockTokenService) ValidateIDToken(tokenString string) (*model.Users, string, *errors.MathSheetsError) {
	ret := m.Called(tokenString)

	// first value passed to "Return"
	var r0 *model.Users
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*model.Users)
	}

	var r1 string

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(string)
	}

	var r2 *errors.MathSheetsError

	if ret.Get(2) != nil {
		r2 = ret.Get(2).(*errors.MathSheetsError)
	}

	return r0, r1, r2
}

// ValidateRefreshToken mocks concrete ValidateRefreshToken
func (m *MockTokenService) ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, *errors.MathSheetsError) {
	ret := m.Called(refreshTokenString)

	var r0 *model.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.RefreshToken)
	}

	var r1 *errors.MathSheetsError

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*errors.MathSheetsError)
	}

	return r0, r1
}

// Signout mocks concrete Signout
func (m *MockTokenService) IsBlackedListed(ctx context.Context, uid string, tokenid string) *errors.MathSheetsError {
	ret := m.Called(ctx, uid, tokenid)
	var r0 *errors.MathSheetsError

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
	}

	return r0
}

package mocks

import (
	"context"

	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock type for model.UserRepository
type MockUserRepository struct {
	mock.Mock
}

// FindByID is mock of UserRepository FindByID
func (m *MockUserRepository) FindByID(ctx context.Context, uid string) (*model.Users, *errors.MathSheetsError) {
	ret := m.Called(ctx, uid)

	var r0 *model.Users
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Users)
	}

	var r1 *errors.MathSheetsError

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*errors.MathSheetsError)
	}

	return r0, r1
}

// Create is mock of UserRepository Create
func (m *MockUserRepository) Create(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	ret := m.Called(ctx, u)

	var r0 *errors.MathSheetsError

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
	}

	return r0
}

// FindByEmail is mock of UserRepository FindByEmail
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.Users, *errors.MathSheetsError) {
	ret := m.Called(ctx, email)

	var r0 *model.Users
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Users)
	}

	var r1 *errors.MathSheetsError

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*errors.MathSheetsError)
	}

	return r0, r1
}

// Update is mock of UserRepository Update
func (m *MockUserRepository) Update(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	ret := m.Called(ctx, u)

	var r0 *errors.MathSheetsError
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
	}

	return r0
}

// Delete in mock of UserRepository Delete
// TODO : create testing
func (m *MockUserRepository) Delete(ctx context.Context, id string) *errors.MathSheetsError {
	ret := m.Called(ctx, id)

	var r0 *errors.MathSheetsError
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
	}

	return r0
}

func (m *MockUserRepository) UpdateImage(ctx context.Context, u *model.Users, imageURL string) (*model.Users, *errors.MathSheetsError) {
	ret := m.Called(ctx, u, imageURL)

	var r0 *model.Users
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Users)
	}

	var r1 *errors.MathSheetsError
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*errors.MathSheetsError)
	}

	return r0, r1
}

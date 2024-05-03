package mocks

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"

	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock type for model.UserService
type MockUserService struct {
	mock.Mock
}

// Get is mock of UserService Get
func (m *MockUserService) Get(ctx context.Context, uid string) (*model.Users, *errors.MathSheetsError) {
	// args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret"
	ret := m.Called(ctx, uid)

	// first value passed to "Return"
	var r0 *model.Users
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*model.Users)
	}

	var r1 *errors.MathSheetsError
	// var tradeCVDErr MathSheetsError
	if ret.Get(1) != nil {

		r1 = ret.Get(1).(*errors.MathSheetsError)

	}

	return r0, r1
}

// Signup is a mock of UserService.Signup
func (m *MockUserService) Signup(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	ret := m.Called(ctx, u)

	var r0 *errors.MathSheetsError
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
	}

	return r0
}

// Signin is a mock of UserService.Signin
func (m *MockUserService) Signin(ctx context.Context, u *model.Users) (*model.Users, *errors.MathSheetsError) {
	ret := m.Called(ctx, u)

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

// UpdateDetails is a mock of UserService.UpdateDetails
func (m *MockUserService) UpdateDetails(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	ret := m.Called(ctx, u)

	var r0 *errors.MathSheetsError
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
	}

	return r0
}

// Delete is a mock of UserService.Delete
func (m *MockUserService) DeleteUser(ctx context.Context, id string) *errors.MathSheetsError {
	ret := m.Called(ctx, id)

	var r0 *errors.MathSheetsError
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*errors.MathSheetsError)
		log.Println(r0)
	}

	return r0
}

// Google Sign In is a mock of UserService.GoogleSignin
func (m *MockUserService) GoogleSignin(ctx context.Context, code string) (*model.Users, *errors.MathSheetsError) {
	ret := m.Called(ctx, code)

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

func (m *MockUserService) SetProfileImage(ctx context.Context, id string, imageFileHeader *multipart.FileHeader) (*model.Users, *errors.MathSheetsError) {
	ret := m.Called(ctx, id, imageFileHeader)

	var r0 *model.Users

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Users)
		log.Println(r0)
	}

	var r1 *errors.MathSheetsError

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*errors.MathSheetsError)
		log.Println(r1)
	}

	return r0, r1
}

// func (m *MockUserService) SetBusinessDocs(ctx context.Context, f *model.BusinessFiles, businessFileHeader *multipart.FileHeader) (*model.BusinessFiles, *errors.MathSheetsError) {
// 	ret := m.Called(ctx, f, businessFileHeader)

// 	var r0 *model.BusinessFiles
// 	if ret.Get(0) != nil {
// 		r0 = ret.Get(0).(*model.BusinessFiles)
// 		// log.Println(r0)
// 	}

// 	var r1 *errors.MathSheetsError
// 	if ret.Get(1) != nil {
// 		r1 = ret.Get(1).(*errors.MathSheetsError)
// 		// log.Println(r0)
// 	}

// 	return r0, r1
// }

// func (m *MockUserService) SignedURLFromFileName(ctx context.Context, objName string) (string, *errors.MathSheetsError) {
// 	ret := m.Called(ctx, objName)

// 	var r0 string
// 	if ret.Get(0) != nil {
// 		r0 = ret.Get(0).(string)
// 		// log.Println(r0)
// 	}

// 	var r1 *errors.MathSheetsError
// 	if ret.Get(1) != nil {
// 		r1 = ret.Get(1).(*errors.MathSheetsError)
// 		// log.Println(r0)
// 	}

// 	return r0, r1
// }

// func (m *MockUserService) ListUserFile(ctx context.Context, id string) ([]string, *errors.MathSheetsError) {
// 	ret := m.Called(ctx, id)

// 	var r0 []string
// 	if ret.Get(0) != nil {
// 		r0 = ret.Get(0).([]string)
// 	}

// 	var r1 *errors.MathSheetsError
// 	if ret.Get(1) != nil {
// 		r1 = ret.Get(1).(*errors.MathSheetsError)
// 	}

// 	return r0, r1
// }

package services

import (
	"context"
	"testing"

	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/model/mocks"
	"github.com/rudyjcruz831/mathSheets/util/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid := "123"

		mockUserResp := &model.Users{
			ID:        uid,
			Email:     "bob@bob.com",
			FirstName: "Bobby",
			LastName:  "Bobson",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, mathSheetErr := us.Get(ctx, uid)
		assert.Nil(t, mathSheetErr)
		// assert.NoError(t, err) // how to use this using mathSheet

		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid := "123"

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mathSheetErr := errors.NewInternalServerError("")
		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, mathSheetErr)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.Nil(t, u)
		// assert.
		assert.NotNil(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid := "123"

		mockUser := &model.Users{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		// We can use Run method to modify the user when the Create method is called.
		// We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.Users) // arg 0 is context, arg 1 is *User
				userArg.ID = uid
			}).Return(nil)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.Nil(t, err)

		// assert user now has a userID
		assert.Equal(t, uid, mockUser.ID)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.Users{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		// mockErr := apperrors.NewConflict("email", mockUser.Email)
		mockErr := errors.NewConflict("email", mockUser.Email)

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Return(mockErr)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		// assert error is error we response with in mock
		assert.Equal(t, err, mockErr)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignin(t *testing.T) {
	// setup valid email/pw combo with hashed password to test method
	// response when provided password is invalid
	email := "rudy@go.com"
	validPW := "password1!"
	hashedValidPW, _ := hashPassword(validPW)
	invalidPW := "howdyhodufus!"

	mockUserRepository := new(mocks.MockUserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid := "1"

		mockUser := &model.Users{
			Email:    email,
			Password: validPW,
		}

		mockUserResp := &model.Users{
			ID:       uid,
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		userSingIn, err := us.Signin(ctx, mockUser)

		assert.Nil(t, err)
		assert.Equal(t, userSingIn, mockUserResp)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})

	t.Run("Invalid email/password combination", func(t *testing.T) {
		uid := "2"

		mockUser := &model.Users{
			Email:    email,
			Password: invalidPW,
		}

		mockUserResp := &model.Users{
			ID:       uid,
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		userSignIn, err := us.Signin(ctx, mockUser)

		mockErr := errors.NewAuthorization("Invalid email and password combination")

		assert.NotNil(t, err)
		assert.Nil(t, userSignIn)
		assert.Equal(t, err, mockErr)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})
}

func TestDeleteUser(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockUserRepository := new(mocks.MockUserRepository)

		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		uid := "3"

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid,
		}

		mockUserRepository.On("Delete", mockArgs...).Return(nil)
		ctx := context.TODO()
		err := us.DeleteUser(ctx, uid)

		assert.Nil(t, err)

	})

	t.Run("Failure", func(t *testing.T) {
		mockUserRepository := new(mocks.MockUserRepository)

		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		uid := "3"

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uid,
		}
		mockmathSheetErr := errors.NewNotFound("User not found: ", uid)
		mockUserRepository.On("Delete", mockArgs...).Return(mockmathSheetErr)
		ctx := context.TODO()
		err := us.DeleteUser(ctx, uid)

		assert.NotNil(t, err)
	})
}

func TestGoogleSignin(t *testing.T) {

}

func TestSetProfileImage(t *testing.T) {

}

func TestSetBusinessDocs(t *testing.T) {

}

func TestListUserFile(t *testing.T) {

}

func TestSignedURLFromFileName(t *testing.T) {

}

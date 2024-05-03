package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/model/mocks"
	"github.com/rudyjcruz831/mathSheets/util/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPairForUser(t *testing.T) {
	var idExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600
	priv, _ := ioutil.ReadFile("../rsa_keys_tokens/rsa_private_dev.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := ioutil.ReadFile("../rsa_keys_tokens/rsa_public_dev.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	mockTokenRepository := new(mocks.MockTokenRepository)
	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TSConfig{
		TokenRepository:       mockTokenRepository,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         secret,
		IDExpirationsSecs:     idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid := "123"
	u := &model.Users{
		ID:       uid,
		Email:    "rudy@rudy.com",
		Password: "blarghedymcblarghface",
	}

	uidErrorCase := "1212334"
	uErrorCase := &model.Users{
		ID:       uidErrorCase,
		Email:    "failure@failure.com",
		Password: "blarghedymcblarghface",
	}
	prevID := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		context.Background(),
		u.ID,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		context.Background(),
		uidErrorCase,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPrevIDArguments := mock.Arguments{
		context.Background(),
		u.ID,
		prevID,
	}

	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mathSheetErrOut := errors.NewInternalServerError("")
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(mathSheetErrOut)
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(nil)

	t.Run("Returns a token pair with proper values", func(t *testing.T) {

		ctx := context.Background()
		tokenPair, mathSheetErr := tokenService.NewPairForUser(ctx, u, prevID) // replace "" with prevID from setup
		assert.Nil(t, mathSheetErr)

		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken.SS)

		// decode the Base64URL encoded string
		// simpler to use jwt library which is already imported
		idTokenClaims := &idTokenCustomClaims{}

		_, err := jwt.ParseWithClaims(tokenPair.IDToken.SS, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			u.ID,
			u.Email,
			u.FirstName,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.ID,
			idTokenClaims.User.Email,
			idTokenClaims.User.FirstName,
			// idTokenClaims.User.ImageURL,
			// idTokenClaims.User.Website,
		}

		fmt.Printf("actualIDClaims: %v\n", actualIDClaims)

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // password should never be encoded to json

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &refreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken.SS, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken.SS)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, u.ID, refreshTokenClaims.UID)
		// TODO Try to get better test for this times
		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {

		ctx := context.Background()
		_, err := tokenService.NewPairForUser(ctx, uErrorCase, "")
		fmt.Printf("mathSheetErr: %v\n", err)
		assert.NotNil(t, err) // should return an error

		// SetRefreshToken should be called with setErrorArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		// DeleteRefreshToken should not be since SetRefreshToken causes method to return
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for prevID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairForUser(ctx, u, "")
		assert.Nil(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Prev token not in repository", func(t *testing.T) {
		ctx := context.Background()
		uid := "123"
		u := &model.Users{
			ID: uid,
		}

		tokenIDNotInRepo := "not_in_token_repo"

		deleteArgs := mock.Arguments{
			ctx,
			u.ID,
			tokenIDNotInRepo,
		}

		mockError := errors.NewAuthorization("Invalid refresh token")
		mockTokenRepository.
			On("DeleteRefreshToken", deleteArgs...).
			Return(mockError)

		_, err := tokenService.NewPairForUser(ctx, u, tokenIDNotInRepo)
		assert.NotNil(t, err)

		assert.Equal(t, mockError, err)
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteArgs...)
		mockTokenRepository.AssertNotCalled(t, "SetRefreshToken")
	})
}

func TestSignout(t *testing.T) {
	mockTokenRepository := new(mocks.MockTokenRepository)
	tokenService := NewTokenService(&TSConfig{
		TokenRepository: mockTokenRepository,
	})
	var idExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600

	priv, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := ioutil.ReadFile("../rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	t.Run("No error", func(t *testing.T) {
		uidSuccess := "123"
		u := &model.Users{
			ID:       uidSuccess,
			Email:    "rudy@go.com",
			Password: "avalidpassword",
		}
		// token will be valid for 15 minutes
		idToken, err := generateIDToken(u, privKey, idExp)
		assert.NoError(t, err)

		mockTokenRepository := new(mocks.MockTokenRepository)
		// instantiate a common token service to be used by all tests
		tokenService1 := NewTokenService(&TSConfig{
			TokenRepository:       mockTokenRepository,
			PrivKey:               privKey,
			PubKey:                pubKey,
			RefreshSecret:         secret,
			IDExpirationsSecs:     idExp,
			RefreshExpirationSecs: refreshExp,
		})

		mockDeleteUserRefreshToken := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uidSuccess,
		}
		mockTokenRepository.On("DeleteUserRefreshTokens", mockDeleteUserRefreshToken...).Return(nil)

		mockClaim, err := validateIDToken(idToken.SS, pubKey)
		assert.NoError(t, err)

		mockIssuedAt := time.Unix(mockClaim.IssuedAt, 0)
		mockExpiresIn := time.Unix(mockClaim.ExpiresAt, 0)
		expIn := mockExpiresIn.Sub(mockIssuedAt)

		mockTokenBlackedListed := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uidSuccess,
			mockClaim.Id,
			expIn,
		}

		mockTokenRepository.On("TokenBlackedListed", mockTokenBlackedListed...).Return(nil)

		// maybe not the best approach to depend on utility method

		assert.NotNil(t, idToken)
		// tokenString := "tokenString"
		ctx := context.Background()
		mathSheetErr := tokenService1.Signout(ctx, uidSuccess)

		assert.Nil(t, mathSheetErr)
		mockTokenRepository.AssertCalled(t, "DeleteUserRefreshTokens", mockDeleteUserRefreshToken...)
		mockTokenRepository.AssertCalled(t, "TokenBlackedListed", mockTokenBlackedListed...)
	})

	// t.Run("Error from DeleteUserRefreshTokens", func(t *testing.T) {
	// 	uidError := "123"
	// 	tokenString := "abc"
	// 	mockTokenRepository.
	// 		On("DeleteUserRefreshTokens", mock.AnythingOfType("*context.emptyCtx"), uidError).
	// 		Return(errors.NewInternalServerError(""))

	// 	ctx := context.Background()
	// 	mathSheetErr := tokenService.Signout(ctx, uidError, tokenString)

	// 	assert.NotNil(t, mathSheetErr)
	// })
	t.Run("Error from TokenBlackedListed", func(t *testing.T) {
		uidError := "123"
		// tokenString := "abc"

		issuedAt := time.Unix(600, 0)
		expiresIn := time.Unix(100, 0)
		expIn := expiresIn.Sub(issuedAt)

		mockTokenBlackedListed := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			uidError,
			"tokenString",
			expIn,
		}
		mockTokenRepository.
			On("DeleteUserRefreshTokens", mock.AnythingOfType("*context.emptyCtx"), uidError).
			Return(nil)

		mockTokenRepository.
			On("TokenBlackedListed", mockTokenBlackedListed...).
			Return(errors.NewInternalServerError(""))

		ctx := context.Background()
		mathSheetErr := tokenService.Signout(ctx, uidError)

		assert.NotNil(t, mathSheetErr)
	})
}

func TestValidateIDToken(t *testing.T) {
	var idExp int64 = 15 * 60

	priv, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := ioutil.ReadFile("../rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TSConfig{
		PrivKey:           privKey,
		PubKey:            pubKey,
		IDExpirationsSecs: idExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid := "123"
	u := &model.Users{
		ID:       uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		idToken, _ := generateIDToken(u, privKey, idExp)

		uFromToken, _, mathSheetErr := tokenService.ValidateIDToken(idToken.SS)
		// assert.NoError(t, err)
		assert.Nil(t, mathSheetErr)

		assert.ElementsMatch(
			t,
			[]interface{}{u.Email, u.FirstName, u.ID},
			[]interface{}{uFromToken.Email, uFromToken.FirstName, uFromToken.ID},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		idToken, _ := generateIDToken(u, privKey, -1) // expires one second ago

		expectedErr := errors.NewAuthorization("Unable to verify user from idToken")

		_, _, err := tokenService.ValidateIDToken(idToken.SS)
		assert.Equal(t, err.Message, expectedErr.Message)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		idToken, _ := generateIDToken(u, privKey, -1) // expires one second ago

		expectedErr := errors.NewAuthorization("Unable to verify user from idToken")

		_, _, err := tokenService.ValidateIDToken(idToken.SS)
		assert.Equal(t, err.Message, expectedErr.Message)
	})

	// TODO - Add other invalid token types
}

func TestValidateRefreshToken(t *testing.T) {
	var refreshExp int64 = 3 * 24 * 2600
	secret := "anotsorandomtestsecret"

	tokenService := NewTokenService(&TSConfig{
		RefreshSecret:         secret,
		RefreshExpirationSecs: refreshExp,
	})

	uid := "123"
	u := &model.Users{
		ID:       uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		testRefreshToken, _ := generateRefreshToken(u.ID, secret, refreshExp)

		validatedRefreshToken, mathSheetErr := tokenService.ValidateRefreshToken(testRefreshToken.SS)
		assert.Nil(t, mathSheetErr)

		assert.Equal(t, u.ID, validatedRefreshToken.UID)
		assert.Equal(t, testRefreshToken.SS, validatedRefreshToken.SS)
		assert.Equal(t, u.ID, validatedRefreshToken.UID)
	})

	t.Run("Expired token", func(t *testing.T) {
		testRefreshToken, _ := generateRefreshToken(u.ID, secret, -1)

		expectedErr := errors.NewAuthorization("Unable to verify user from refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken.SS)
		assert.Equal(t, err, expectedErr)
	})
}

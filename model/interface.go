package model

import (
	"bytes"
	"context"
	"time"

	"github.com/rudyjcruz831/mathSheets/util/errors"
)

// UserService defines methods the handler layer expects
type UserService interface {
	Get(ctx context.Context, id string) (*Users, *errors.MathSheetsError)
	Signup(ctx context.Context, u *Users) *errors.MathSheetsError
	Signin(ctx context.Context, u *Users) (*Users, *errors.MathSheetsError)
	// UpdateDetails(ctx context.Context, u *Users) *errors.MathSheetsError
	DeleteUser(ctx context.Context, id string) *errors.MathSheetsError
	GoogleSignin(ctx context.Context, code string) (*Users, *errors.MathSheetsError)
	CreatePDF(ctx context.Context, grade string, subject string) (bytes.Buffer, *errors.MathSheetsError)
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*Users, *errors.MathSheetsError)
	Create(ctx context.Context, u *Users) *errors.MathSheetsError
	FindByEmail(ctx context.Context, email string) (*Users, *errors.MathSheetsError)
	Update(ctx context.Context, u *Users) *errors.MathSheetsError
	Delete(ctx context.Context, id string) *errors.MathSheetsError
	// UpdateImage(ctx context.Context, u *Users, imageURL string) (*Users, *errors.MathSheetsError)
}

// TokenService defines methods handler layer expects to interact
// with in regards to producing JWT as string
type TokenService interface {
	NewPairForUser(ctx context.Context, u *Users, prevTokenID string) (*TokenPair, *errors.MathSheetsError)
	Signout(ctx context.Context, uid string) *errors.MathSheetsError
	ValidateIDToken(tokenString string) (*Users, string, *errors.MathSheetsError)
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, *errors.MathSheetsError)
	// IsBlackedListed(ctx context.Context, uid string, tokenid string) *errors.MathSheetsError
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.MathSheetsError
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) *errors.MathSheetsError
	DeleteUserRefreshTokens(ctx context.Context, userID string) *errors.MathSheetsError
	// TokenBlackedListed(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.MathSheetsError
}

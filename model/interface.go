package model

import (
	"context"

	"github.com/rudyjcruz831/mathSheets/util/errors"
)

// UserService defines methods the handler layer expects
type UserService interface {
	Get(ctx context.Context, id string) (*Users, *errors.MathSheetsError)
	Signup(ctx context.Context, u *Users) *errors.MathSheetsError
	Signin(ctx context.Context, u *Users) (*Users, *errors.MathSheetsError)
	// UpdateDetails(ctx context.Context, u *Users) *errors.MathSheetsError
	// DeleteUser(ctx context.Context, id string) *errors.MathSheetsError
	// GoogleSignin(ctx context.Context, code string) (*Users, *errors.MathSheetsError)
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*Users, *errors.MathSheetsError)
	Create(ctx context.Context, u *Users) *errors.MathSheetsError
	FindByEmail(ctx context.Context, email string) (*Users, *errors.MathSheetsError)
	Update(ctx context.Context, u *Users) *errors.MathSheetsError
	Delete(ctx context.Context, id string) *errors.MathSheetsError
	UpdateImage(ctx context.Context, u *Users, imageURL string) (*Users, *errors.MathSheetsError)
}

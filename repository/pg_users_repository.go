package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"
	"gorm.io/gorm"
)

// PGUserRepository is data/repository implementation
// of service layer UserRepository
type pGUserRepository struct {
	DB *gorm.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &pGUserRepository{
		DB: db,
	}
}

// FindByID fetches user by id
func (r *pGUserRepository) FindByID(ctx context.Context, id string) (*model.Users, *errors.MathSheetsError) {
	// panic("Create function in Pg user repository")
	u := &model.Users{}

	result := r.DB.First(&u, id)
	if result.Error != nil {
		tradeCVDErr := errors.NewNotFound("id", id)
		return nil, tradeCVDErr
	}

	return u, nil
}

// Create reaches out to database postrges using gorm api
func (r *pGUserRepository) Create(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	// panic("Create function in Pg user repository")
	uid, _ := uuid.NewRandom()
	// query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"
	u.CreatedOn = time.Now()
	u.UpdatedAt = time.Now()
	u.ID = uid.String()
	if result := r.DB.FirstOrCreate(&u, u); result.Error != nil {
		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, result.Error)
		tradeCVDErr := errors.NewConflict("email", u.Email)
		return tradeCVDErr
	}

	return nil
}

// FindByEmail fetches user by email
func (r *pGUserRepository) FindByEmail(ctx context.Context, email string) (*model.Users, *errors.MathSheetsError) {
	// panic("FindByEmail in pGUserRepository")

	u := &model.Users{}

	// using gorm to hit postgresDB using email
	if result := r.DB.Where("email = ?", email).First(u); result.Error != nil {
		log.Printf("Db error : %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			tradeCVDErr := errors.NewNotFound("email", email)
			return nil, tradeCVDErr
		} else {
			tradeCVDErr := errors.NewInternalServerError("")
			return nil, tradeCVDErr
		}

	}

	return u, nil
}

// Update updates user information
func (r *pGUserRepository) Update(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	userInRepo := &model.Users{}
	result := r.DB.Where("email = ?", u.Email).First(userInRepo)
	if result.Error != nil {
		log.Printf("Db error: %v", result.Error)
		tradeCVDErr := errors.NewNotFound("email", u.Email)
		return tradeCVDErr
	}

	userInRepo.FirstName = u.FirstName
	userInRepo.LastName = u.LastName

	result = r.DB.Save(userInRepo)
	if result.Error != nil {
		log.Printf("Db error %v", result.Error)
		tradeCVDErr := errors.NewInternalServerError("")
		return tradeCVDErr
	}

	return nil
}

// Delete deletes user information
func (r *pGUserRepository) Delete(ctx context.Context, id string) *errors.MathSheetsError {
	u := &model.Users{}
	results := r.DB.Delete(u, id)
	if results.Error != nil {
		return errors.NewInternalServerError("")
	}

	return nil
}

func (r *pGUserRepository) UpdateImage(ctx context.Context, u *model.Users, imageURL string) (*model.Users, *errors.MathSheetsError) {

	// must be instantiated to scan into ref using `GetContext`
	userInRepo := &model.Users{}

	result := r.DB.Where("email = ?", u.Email).First(userInRepo)
	if result.Error != nil {
		log.Printf("Db error: %v", result.Error)
		tradeCVDErr := errors.NewNotFound("email", u.Email)
		return nil, tradeCVDErr
	}
	//update the user Image
	userInRepo.Image = imageURL
	result = r.DB.Save(userInRepo)
	if result.Error != nil {
		log.Printf("Db error %v", result.Error)
		tradeCVDErr := errors.NewInternalServerError("")
		return nil, tradeCVDErr
	}
	return userInRepo, nil
}

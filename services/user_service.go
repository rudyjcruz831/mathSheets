package services

import (
	"context"
	"log"

	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"
)

// userService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository model.UserRepository
	// DocsRepository  model.DocsRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository model.UserRepository
	// DocsRepository  model.DocsRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
		// DocsRepository:  c.DocsRepository,
	}
}

// Get retrieves a user based on their ID
func (s *userService) Get(ctx context.Context, id string) (*model.Users, *errors.MathSheetsError) {
	u, err := s.UserRepository.FindByID(ctx, id)
	return u, err
}

// Signup reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userService) Signup(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	pw, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", u.Email)
		return errors.NewInternalServerError("")
	}
	u.Password = pw
	u.Role = "user"
	if mathShtErr := s.UserRepository.Create(ctx, u); mathShtErr != nil {
		log.Printf("UserRepository return error: %v", mathShtErr)
		return mathShtErr
	}

	// If we get around to adding events, we'd Publish it here
	// err := s.EventsBroker.PublishUserUpdated(u, true)

	// if err != nil {
	// 	return nil, apperrors.NewInternal()
	// }

	return nil
}

// Signin reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userService) Signin(ctx context.Context, u *model.Users) (*model.Users, *errors.MathSheetsError) {
	// panic("Sign In Method not implemented")
	uFetched, MathShtErr := s.UserRepository.FindByEmail(ctx, u.Email)

	// Will return NotAuthorized to client to omit details of why
	if MathShtErr != nil {
		log.Printf("FindByEmail return error : %v", MathShtErr)
		return nil, errors.NewAuthorization("Invalid email and password combination")
	}

	match, err := comparePasswords(uFetched.Password, u.Password)

	if err != nil {
		log.Printf("comparePassword return error %v", err)
		return nil, errors.NewInternalServerError("")
	}

	if !match {
		log.Println("Match was false return error")
		return nil, errors.NewAuthorization("Invalid email and password combination")
	}

	return uFetched, nil
}

// Update Details reaches out
func (s *userService) UpdateDetails(ctx context.Context, u *model.Users) *errors.MathSheetsError {
	// Update user in UserRepository
	MathShtErr := s.UserRepository.Update(ctx, u)

	if MathShtErr != nil {
		return MathShtErr
	}

	// // Publish user updated nats streaming server // kafca
	// err = s.EventsBroker.PublishUserUpdated(u, false)
	// if err != nil {
	// 	return errors.NewInternal()
	// }

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) *errors.MathSheetsError {
	// panic("Delete user service")
	mathSheetsErr := s.UserRepository.Delete(ctx, id)
	return mathSheetsErr
}

// Using google to sign in user if no user will be created calling the repo
func (s *userService) GoogleSignin(ctx context.Context, code string) (*model.Users, *errors.MathSheetsError) {
	// panic("Google Sing In")
	// TODO - testing for this service
	_, u, MathShtErr := auth(code)
	if MathShtErr != nil {
		// c.JSON(MathShtErr.Status, MathShtErr)
		return nil, MathShtErr
	}

	uFetched, MathShtErr := s.UserRepository.FindByEmail(ctx, u.Email)
	if MathShtErr != nil {
		if MathShtErr.Status == 404 {
			u.Role = "user"
			if MathShtErr = s.UserRepository.Create(ctx, u); MathShtErr != nil {
				return nil, MathShtErr
			}
			return u, nil
		} else {
			return nil, MathShtErr
		}
	}

	return uFetched, nil
}

func (s *userService) CreatePDF(ctx context.Context, grade string, subject string) *errors.MathSheetsError {
	// here we are going to try to qurey OpenAI API
	panic("CreatedPDF")
}

// func (s *userService) SetProfileImage(ctx context.Context, id string, imageFileHeader *multipart.FileHeader) (*model.Users, *errors.MathSheetsError) {
// 	// TODO - testing for this service also I will need to come in here and make sure if I am even adding an image to user
// 	u, err := s.UserRepository.FindByID(ctx, id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	objName, err := objNameFromURL(u.Image)

// 	if err != nil {
// 		return nil, err
// 	}

// 	imageFile, err1 := imageFileHeader.Open()
// 	if err1 != nil {
// 		log.Printf("Failed to open image file: %v\n", err1)
// 		return nil, errors.NewInternalServerError("")
// 	}

// 	// Upload user's image to ImageRepository
// 	// Possibly received updated imageURL
// 	imageURL, MathShtErr := s.ImageRepository.UpdateProfile(ctx, objName, imageFile)

// 	if MathShtErr != nil {
// 		log.Printf("Unable to upload image to cloud provider: %v\n", MathShtErr)
// 		return nil, MathShtErr
// 	}

// 	updatedUser, err := s.UserRepository.UpdateImage(ctx, u, imageURL)

// 	if err != nil {
// 		log.Printf("Unable to update imageURL: %v\n", err)
// 		return nil, err
// 	}

// 	return updatedUser, nil
// }

// This function has to do with the image as well
// func objNameFromURL(imageURL string) (string, *errors.MathSheetsError) {
// 	// if user doesn't have imageURL - create one
// 	// otherwise, extract last part of URL to get cloud storage object name
// 	if imageURL == "" {
// 		objID, _ := uuid.NewRandom()
// 		return objID.String(), nil
// 	}

// 	// split off last part of URL, which is the image's storage object ID
// 	urlPath, err := url.Parse(imageURL)

// 	if err != nil {
// 		log.Printf("Failed to parse objectName from imageURL: %v\n", imageURL)
// 		return "", errors.NewInternalServerError("")
// 	}

// 	// get "path" of url (everything after domain)
// 	// then get "base", the last part
// 	return path.Base(urlPath.Path), nil
// }

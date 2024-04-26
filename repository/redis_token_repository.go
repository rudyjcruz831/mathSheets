package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"
)

// redisTokenRepository is data/repository implementation
// of service layer TokenRepository
type redisTokenRepository struct {
	Redis *redis.Client
}

// NewTokenRepository is a factory for initializing User Repositories
func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// SetRefreshToken stores a refresh token with an expiry time
func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.MathSheetsError {
	// panic("Set Refresh actuall iplmentation to interact with redis")
	// We'll store userID with token id so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return errors.NewInternalServerError("")
	}
	return nil
}

// 2024/04/24 22:14:47 Refresh token to redis for userID: dcd914f9-3406-4a45-940d-97f980fb07f9:11205908-b238-49ce-93e5-b7c9d552d8a8 does not exist
// 2024/04/24 22:14:47 Could not delete previous refreshToken for uid: dcd914f9-3406-4a45-940d-97f980fb07f9, tokenID: 11205908-b238-49ce-93e5-b7c9d552d8a8
// 2024/04/24 22:14:47 Failed to create tokens for user: &{ID:dcd914f9-3406-4a45-940d-97f980fb07f9 Email:rudy5@go.com Username:rudy5
// Password:f04fb19f15b9d68c1c610eb29ab0de4f18e13f7c7f721cfe37c5e101580b3782.2d41e61336a7c94fbb601e717f47bd32739662e7ab7b01eac24a9b57e56189a7 FirstName: LastName: Role:user CreatedOn:2024-04-13 04:17:57.190229 +0000 UTC UpdatedAt:2024-04-13 04:17:57.190229 +0000 UTC DeletedAt:<nil>}. Error: &{0001-01-01 00:00:00 +0000 UTC 401 AUTHORIZATION Invalid refresh token }

// DeleteRefreshToken used to delete old  refresh tokens
// Services my access this to revolve tokens
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) *errors.MathSheetsError {
	// panic("jlkjflkajsf")
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return errors.NewInternalServerError("")
	}

	// // Val returns count of deleted keys
	// // If no value was deleted, the refresh token is invalid
	// fmt.Print(result.Val(), "\n")
	// if result.Val() < 1 {
	// 	log.Printf("Refresh token to redis for userID: %s:%s does not exist\n", userID, tokenID)
	// 	return errors.NewAuthorization("Invalid refresh token")
	// }

	return nil
}

// DeleteUserRefreshTokens looks for all tokens beginning with
// userID and scans to delete them in a non-blocking fashion
func (r *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) *errors.MathSheetsError {
	pattern := fmt.Sprintf("%s*", userID)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", iter.Val())
			failCount++
		}
	}

	// check last value
	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token token iter Error: %s\n", iter.Val())
	}

	if failCount > 0 {
		return errors.NewInternalServerError("")
	}

	return nil
}

// TokenBlackedListed add tokens to redis to when signed out
// func (r *redisTokenRepository) TokenBlackedListed(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.MathSheetsError {
// 	// panic("")
// 	// We'll store userID with token id so we can scan (non-blocking)
// 	// over the user's tokens and delete them in case of token leakage
// 	key := fmt.Sprintf("%s:%s", userID, tokenID)
// 	if err := r.Redis.Set(ctx, key, 1, expiresIn).Err(); err != nil {
// 		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
// 		return errors.NewInternalServerError("")
// 	}
// 	return nil
// }

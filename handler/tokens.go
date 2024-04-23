package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (h *Handler) Tokens(c *gin.Context) {
	fmt.Print("Tokens handler\n")
	var req tokenReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	//verify refresh JWT
	refreshToken, mathSheetsErr := h.TokenService.ValidateRefreshToken(req.RefreshToken)
	if mathSheetsErr != nil {
		log.Printf("Failed to validate refresh token: %v\n", mathSheetsErr)
		c.JSON(mathSheetsErr.Status, mathSheetsErr)
		return
	}

	// get up-to-date user
	u, mathSheetsErr := h.UserService.Get(ctx, refreshToken.UID)

	if mathSheetsErr != nil {
		log.Printf("Failed to get user: %v\n", mathSheetsErr)
		c.JSON(mathSheetsErr.Status, mathSheetsErr)
		return
	}

	tokens, mathSheetsErr := h.TokenService.NewPairForUser(ctx, u, refreshToken.ID.String())

	if mathSheetsErr != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", u, mathSheetsErr)
		c.JSON(mathSheetsErr.Status, mathSheetsErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokens,
	})
}

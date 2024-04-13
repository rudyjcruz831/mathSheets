package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/mathSheets/model"
)

type signinReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Signin used to authenticate extant user
func (h *Handler) SignIn(c *gin.Context) {
	var req signinReq

	if ok := bindData(c, &req); !ok {
		fmt.Println("binding data unsuccessful")
		return
	}

	// inject user request to user to use for User service layer
	u := &model.Users{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	uNew, mathSheetErr := h.UserService.Signin(ctx, u)

	if mathSheetErr != nil {
		log.Printf("Failed to sign in user: %v\n", mathSheetErr)
		// mathSheetErr := errors.NewInternalServerError("Sign in error")
		c.JSON(mathSheetErr.Status, mathSheetErr)
		return
	}

	tokens, mathSheetErr := h.TokenService.NewPairForUser(ctx, uNew, "")

	if mathSheetErr != nil {
		log.Printf("Failed to create tokens for user: %v\n", mathSheetErr)
		c.JSON(mathSheetErr.Status, mathSheetErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}

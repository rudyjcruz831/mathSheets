package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/mathSheets/model"
)

func (h *Handler) SignOut(c *gin.Context) {
	// TODO: need to make sure how I am adding the token to blacklist
	// here I am using a user and tokenstring from the middleware
	// this is not correct
	// TODO: make sure that i am getting the user becasue as how it is now it will not work
	// I have check this in middleware/auth_user.go
	user := c.MustGet("user").(*model.Users)

	ctx := c.Request.Context()
	if mathSheetsErr := h.TokenService.Signout(ctx, user.ID); mathSheetsErr != nil {
		c.JSON(mathSheetsErr.Status, mathSheetsErr)
		return
	}

	// Respond with success message or redirect to home page
	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed out"})
}

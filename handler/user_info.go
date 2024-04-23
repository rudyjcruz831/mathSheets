package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"
)

// type userInfoReq struct {
// 	ID string `json:"id" binding:"required"`
// }

// Me handler calls services for getting
// a user's details
func (h *Handler) UserInfo(c *gin.Context) {
	// var req userInfoReq
	// uid := user.(*model.Users).ID

	// A *model.User will eventually be added to context in middleware
	user, exists := c.Get("user")

	// This shouldn't happen, as our middleware ought to throw an error.
	// This is an extra safety measure
	// We'll extract this logic later as it will be common to all handler
	// methods which require a valid user
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		tradeCVDErr := errors.NewInternalServerError("Unable to extract user from request context for unknown reasons")
		c.JSON(tradeCVDErr.Status, tradeCVDErr)
		return
	}
	uid := user.(*model.Users).ID
	// log.Printf("")
	log.Printf("uid: %s\n", uid)

	// use the Request Context
	ctx := c.Request.Context()
	u, tradeCVDErr1 := h.UserService.Get(ctx, uid)

	if tradeCVDErr1 != nil {
		log.Printf("Unable to find user: %v\n%v", uid, *tradeCVDErr1)
		tradeCVDErr := errors.NewNotFound("userID:", uid)
		c.JSON(tradeCVDErr.Status, tradeCVDErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

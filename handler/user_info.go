package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/mathSheets/util/errors"
)

type userInfoReq struct {
	ID string `json:"id" binding:"required"`
}

// Me handler calls services for getting
// a user's details
func (h *Handler) UserInfo(c *gin.Context) {
	var req userInfoReq
	// uid := user.(*model.Users).ID
	uid := req.ID
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

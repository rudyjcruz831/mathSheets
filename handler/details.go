package handler

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rudyjcruz831/mathSheets/model"
// )

// // omitempty must be listed first (tags evaluated sequentially, it seems)
// type detailsReq struct {
// 	FirstName   string `json:"first_name" binding:"omitempty,max=50"`
// 	LastName    string `json:"last_name" binding:"omitempty,max=50"`
// 	Email       string `json:"email" binding:"required,email"`
// 	PhoneNumber string `json:"phone_number" binding:"omitempty"`
// 	Image       string `json:"image" `
// }

// // // @Details godoc
// // // @Summary for details handler
// // // @Description This handler updates the details of the user
// // // @Tags -
// // // @ID -
// // // @Accept  json {
// // //  }
// // // @Produce  json {
// // // }
// // // @Success 200 {user} user.response
// // // @Failure 400 {object} tradeCVDErr.
// // // @Router /actions [get]
// func (h *Handler) Details(c *gin.Context) {
// 	authUser := c.MustGet("user").(*model.Users)

// 	var req detailsReq

// 	if ok := bindData(c, &req); !ok {
// 		return
// 	}

// 	// Should be returned with current imageURL
// 	u := &model.Users{
// 		ID:        authUser.ID,
// 		FirstName: req.FirstName,
// 		LastName:  req.LastName,
// 		Email:     req.Email,
// 		// PhoneNumber: req.PhoneNumber,
// 		Image: req.Image,
// 		// Business:     req.Business,
// 		// BusinessName: req.BusinessName,
// 		// Street:       req.Street,
// 		// State:        req.State,
// 		// Country:      req.Country,
// 		// Zipcode:      req.Zipcode,
// 		// City:         req.City,
// 	}

// 	ctx := c.Request.Context()

// 	// log.Printf("error : %v", tradeCVDErr)
// 	tradeCVDErr := h.UserService.UpdateDetails(ctx, u)
// 	if tradeCVDErr != nil {
// 		log.Printf("Failed to update user: %v\n", tradeCVDErr)
// 		c.JSON(tradeCVDErr.Status, tradeCVDErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"user": u,
// 	})
// }

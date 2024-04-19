package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//

type createPdfReq struct {
	Grade   string `json:"grade"`
	Subject string `json:"subject"`
}

func (h *Handler) CreatePdf(c *gin.Context) {
	var req createPdfReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	//inject req info to User service to get math sheet
	mathSheetErr := h.UserService.CreatePDF(ctx, req.Grade, req.Subject)
	if mathSheetErr != nil {
		log.Printf("Failed to create PDF: %v\n", mathSheetErr)
		c.JSON(mathSheetErr.Status, mathSheetErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": "created",
	})
}

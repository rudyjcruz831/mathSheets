package handler

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//

type createPdfReq struct {
	Grade   string `json:"grade"`
	Subject string `json:"subject"`
}

func (h *Handler) CreatePDF(c *gin.Context) {
	var req createPdfReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	//inject req info to User service to get math sheet
	buf, mathSheetErr := h.UserService.CreatePDF(ctx, req.Grade, req.Subject)
	if mathSheetErr != nil {
		log.Printf("Failed to create PDF: %v\n", mathSheetErr)
		c.JSON(mathSheetErr.Status, mathSheetErr)
		return
	}

	// Convert bytes to base64 string
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Set the content-type header
	c.Header("Content-Type", "application/pdf")
	c.JSON(http.StatusCreated, gin.H{
		"pdf": base64String,
	})
}

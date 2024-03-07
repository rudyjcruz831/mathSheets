package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/mathSheets/handler"
)

func inject() (*gin.Engine, error) {
	log.Println("Injecting data source...")

	/*
		Repository Layer
	*/
	// TODO : add this when creating the AWS s3 bucket to store worksheets created history
	// workSheetBucketName := os.Getenv("AWS_FILE_BUCKET")
	// userRepository := repository.NewUserRepository(d.DB)
	// imageRepository := repository.NewWorkSheetRepository(d.StorageClient, workSheetBucketName)

	/*
		Service Layer
	*/

	// userServcie := services.NewUserService(&services.USConfig{
	// 	UserRepository:      userRepository,
	// 	WorkSheetRepository: workSheetRepository,
	// })

	// initialize gin.Engine
	router := gin.Default()

	// read in MathSheets API url
	baseURL := os.Getenv("MATHSHEET_API_URL")

	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		R:                router,
		BaseURL:          baseURL,
		TimeoutDurations: time.Duration(time.Duration(ht) * time.Second),
		MaxBodyBytes:     1024 * 1024 * 1024,
	})

	return router, nil
}

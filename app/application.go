package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartApp() {
	// fmt.Println("Hello World")
	log.Println("Start server...")

	// initialze data source
	// TODO: connect DB like for example S3 buckest to store mathsheets created before and may be answers
	ds, err := initDS()
	if err != nil {
		log.Fatalf("Error intilizing data source: %v\n", err)
	}

	//injection other services and add any env variables

	router, err := inject(ds) // remove ds for now becuse we don't have no db for now [router, err := inject(ds)]
	if err != nil {
		log.Fatalf("Failure to inject data source: %v\n", err)
	}

	//grabbing port from env for running server local or other host
	port := os.Getenv("USER_API_PORT")
	//if port env is empty the make it default 50052
	if port == "" {
		port = "50052"
	}

	// Graceful server shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		// running the server on localhost with given port
		//------------------------- This CODE HERE is for TLS-------------------------
		// if err := srv.ListenAndServeTLS("./server.crt", "./server.pem"); err != nil && err != http.ErrServerClosed{
		// 	log.Fatalf("Failed to intialized server: %v\n", err)
		// }

		// This code is without TLS
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to intialized server: %v\n", err)
		}
	}()

	fmt.Printf("Listining on port %s\n", port)

	//wait for kill signal of channel
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This block until a signal is passed into the quit channel
	<-quit
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data source
	// if err := ds.Close(); err != nil {
	// 	log.Fatalf("Problem occured graceful shutting down data source: %v\n", err)
	// }

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

}

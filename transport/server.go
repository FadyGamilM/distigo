package transport

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpServer(r *gin.Engine, addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func RunServer(server *http.Server) {
	if err := server.ListenAndServe(); err != nil {
		log.Printf("error trying to run the server ➜ %v\n", err)
	}
}

func ShutdownGracefully(server *http.Server) {
	// Define a context with timeout.
	// This timeout is the time available for the server to finish
	// whatever requests are running currently before being forced to shut down.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("The server is shutting down now ...")

	// Start closing all connections.
	// If the context that is passed is expired, we will receive an error.
	if err := server.Shutdown(ctx); err != nil {
		// This will be done only if there is an error while trying to shut down the server.
		log.Printf("Server is forced to shutdown ➜ %v\n", err)
	}

	// Log a message indicating successful shutdown.
	log.Println("Server has gracefully shutdown ... ")
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"paradox_server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Paradox Server Running")
	})

	router.POST("/upload", handlers.HandleUpload)

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Printf("Fatal: Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

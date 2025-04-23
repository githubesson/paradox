package routes

import (
	"database/sql"
	"paradox_server/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUploadRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.POST("/upload", func(c *gin.Context) {
		handlers.HandleUpload(c, db)
	})
}

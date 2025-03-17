package routes

import (
	"book-service/internal/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	group := router.Group("api/v1")

	api.BookRoutes(group.Group("/books"), db)

	return router

}

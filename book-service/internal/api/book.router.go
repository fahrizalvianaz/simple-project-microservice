package api

import (
	"book-service/internal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BookRoutes(router *gin.RouterGroup, db *gorm.DB) {
	bookRepository := internal.NewBookRepository(db)
	bookService := internal.NewBookService(bookRepository)
	bookHandler := NewBookHandler(bookService)

	router.POST("/add", bookHandler.Create)
	router.GET("/:id", bookHandler.FindByID)

}

package api

import (
	users "bookstore-framework/internal"

	pkg "github.com/fahrizalvianaz/shared-middleware"
	"github.com/fahrizalvianaz/shared-middleware/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UsersRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepository := users.NewUserRepository(db)
	jwtGenerator := &pkg.Claims{}
	userService := users.NewUserService(userRepository, jwtGenerator)
	userHandler := NewUserHandler(userService)

	router.POST("/register", userHandler.RegisterHandler)
	router.POST("/login", userHandler.LoginHandler)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuth())
	protected.GET("/profile", userHandler.GetProfile)
}

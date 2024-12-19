package routes

import (
	"backend/internal/users/delivery"
	"backend/internal/users/repository"
	"backend/internal/users/usecase"
	"backend/middleware"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Setup User Repository dan Usecase
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewUserHandler(userUsecase)

	// Setup routes untuk User
	router.POST("/api/login", userHandler.Login)

	// Routes dengan autentikasi JWT
	auth := router.Group("/api")
	auth.Use(middleware.JWTMiddleware(db)) // Menggunakan JWT Middleware
	{
		auth.GET("/users", userHandler.GetAllUsers)
		auth.POST("/users/delete", userHandler.DeleteUser)
		auth.POST("/users", userHandler.CreateUser)

	}
	for _, route := range router.Routes() {
		log.Printf("Route %s %s", route.Method, route.Path)
	}
}

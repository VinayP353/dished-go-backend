package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/dished-go-backend/internal/config"
	"github.com/yourusername/dished-go-backend/internal/database"
	"github.com/yourusername/dished-go-backend/internal/handlers"
	"github.com/yourusername/dished-go-backend/internal/middleware"
	"github.com/yourusername/dished-go-backend/internal/models"
	"github.com/yourusername/dished-go-backend/internal/repository"
	"github.com/yourusername/dished-go-backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.CookProfile{}, &models.Cook{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	redisClient, err := database.NewRedisClient(cfg.GetRedisAddr(), cfg.Redis.Password)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, redisClient)
	userHandler := handlers.NewUserHandler(userService)

	cookRepo := repository.NewCookRepository(db)
	cookService := service.NewCookService(cookRepo, redisClient)
	cookHandler := handlers.NewCookHandler(cookService)

	cookProfileRepo := repository.NewCookProfileRepository(db)
	cookProfileService := service.NewCookProfileService(cookProfileRepo, redisClient)
	cookProfileHandler := handlers.NewCookProfileHandler(cookProfileService)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		cooks := api.Group("/cooks")
		{
			cooks.POST("", cookHandler.CreateCook)
			cooks.POST("/login", cookHandler.Login)
			cooks.GET("", cookHandler.GetAllCooks)
			cooks.GET("/:id", cookHandler.GetCook)
			cooks.PUT("/:id", cookHandler.UpdateCook)
			cooks.DELETE("/:id", cookHandler.DeleteCook)
		}

		profiles := api.Group("/cook-profiles")
		{
			profiles.POST("", cookProfileHandler.CreateProfile)
			profiles.GET("", cookProfileHandler.GetAllProfiles)
			profiles.GET("/:id", cookProfileHandler.GetProfile)
			profiles.PUT("/:id", cookProfileHandler.UpdateProfile)
			profiles.DELETE("/:id", cookProfileHandler.DeleteProfile)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

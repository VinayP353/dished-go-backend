// @title           Dished API
// @version         1.0
// @description     Backend API for the Dished chef platform.
// @host            localhost:8081
// @BasePath        /api/v1

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yourusername/dished-go-backend/docs"
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

	if err := db.AutoMigrate(&models.User{}, &models.ChefProfile{}, &models.Chef{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	redisClient, err := database.NewRedisClient(cfg.GetRedisAddr(), cfg.Redis.Password)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, redisClient)
	userHandler := handlers.NewUserHandler(userService)

	chefRepo := repository.NewChefRepository(db)
	chefProfileRepo := repository.NewChefProfileRepository(db)
	chefService := service.NewChefService(chefRepo, chefProfileRepo, redisClient)
	chefHandler := handlers.NewChefHandler(chefService)

	chefProfileService := service.NewChefProfileService(chefProfileRepo, redisClient)
	chefProfileHandler := handlers.NewChefProfileHandler(chefProfileService)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

		auth := api.Group("/auth")
		{
			auth.POST("/register", chefHandler.Register)
			auth.POST("/login", chefHandler.Login)
		}

		chefs := api.Group("/chefs")
		{
			chefs.GET("/usernames", chefHandler.GetUsernames)
			chefs.GET("", chefHandler.GetAllChefs)
			chefs.GET("/:id", chefHandler.GetChef)
			chefs.PUT("/:id", chefHandler.UpdateChef)
			chefs.PUT("/:id/profile", chefHandler.UpdateProfile)
			chefs.DELETE("/:id", chefHandler.DeleteChef)
		}

		profiles := api.Group("/chef-profiles")
		{
			profiles.POST("", chefProfileHandler.CreateProfile)
			profiles.GET("", chefProfileHandler.GetAllProfiles)
			profiles.GET("/:id", chefProfileHandler.GetProfile)
			profiles.PUT("/:id", chefProfileHandler.UpdateProfile)
			profiles.DELETE("/:id", chefProfileHandler.DeleteProfile)
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"log"
	"testwire/config"
	"testwire/internal/repository"
	"testwire/internal/wire"
	"testwire/logs"
	"testwire/routes"

	_ "testwire/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title My API
// @version 1.0
// @description API for my project
// @host localhost:8080
// @BasePath /
func main() {
	logs.Init()

	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	config.AppConfig = *appConfig
	router.Use(gin.Recovery(), gin.Logger())
	config.Connect(appConfig)

	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	repository.SeedRolesAndPermissions()

	routes.AuthRoute(*app.AuthController, router)
	routes.UserRoute(*app.UserController, app.Middleware, router)
	routes.ProductRoute(*app.ProductController, app.Middleware, router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}

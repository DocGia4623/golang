package main

import (
	"log"
	"testwire/config"
	"testwire/internal/middleware"
	"testwire/internal/repository"
	"testwire/internal/wire"
	"testwire/routes"

	_ "testwire/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

// @title My API
// @version 1.0
// @description API for my project
// @host localhost:8080
// @BasePath /
func main() {
	router := gin.New()
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	router.Use(gin.Recovery(), gindump.Dump(), middleware.Logger())
	config.Connect(appConfig)

	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("❌ Failed to initialize app: %v", err)
	}

	repository.SeedRolesAndPermissions()

	routes.AuthRoute(*app.AuthController, router)
	routes.UserRoute(*app.UserController, app.Middleware, router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}

type LogMessage struct {
	Level   string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Service string `json:"service"`
}

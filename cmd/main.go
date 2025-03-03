package main

import (
	"log"
	"testwire/config"
	"testwire/internal/middleware"
	"testwire/internal/repository"
	"testwire/internal/wire"
	"testwire/routes"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

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
	router.Run(":8080")
}

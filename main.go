package main

import (
	"github.com/joho/godotenv"
	"os"
	"tuxiaocao/configs"
	"tuxiaocao/middleware"
	"tuxiaocao/pkg/logger"
	routes2 "tuxiaocao/routes"
	"tuxiaocao/setup"
	"tuxiaocao/utils"

	"github.com/gofiber/fiber/v2"

	_ "tuxiaocao/docs" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		panic("config is not load" + err.Error())
	}
	// Define Fiber config.
	config := configs.FiberConfig()
	setup.InitAll("logger", "mysql")
	// Define a new Fiber service with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for service.
	// Routes.
	routes2.PublicRoutes(app) // Register a public routes for service.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}

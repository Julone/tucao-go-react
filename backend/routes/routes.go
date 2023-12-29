package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	controllers2 "tuxiaocao/routes/controllers"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(app *fiber.App) {
	// Create routes group.
	pubRoute := app.Group("/api/v1")
	// Routes for GET method:
	pubRoute.Get("/products", controllers2.Getproducts)   // get list of all products
	pubRoute.Get("/product/:id", controllers2.Getproduct) // get one product by ID
	// Routes for POST method:
	pubRoute.Post("/user/sign/up", controllers2.UserSignUp) // register app new user
	pubRoute.Post("/user/sign/in", controllers2.UserSignIn) // auth, return Access & Refresh tokens

	// Create routes group.
	route := app.Group("/api/v1")
	// Routes for POST method:
	route.Post("/product", controllers2.Createproduct)     // create app new product
	route.Post("/user/sign/out", controllers2.UserSignOut) // de-authorization user
	route.Post("/token/renew", controllers2.RenewTokens)   // renew Access & Refresh tokens
	// Routes for PUT method:
	route.Put("/product", controllers2.Updateproduct) // update one product by ID
	// Routes for DELETE method:
	route.Delete("/product", controllers2.Deleteproduct) // delete one product by ID

	// Create routes group.
	swagRoute := app.Group("/swagger")
	// Routes for GET method:
	swagRoute.Get("*", swagger.HandlerDefault) // get one user by ID
	app.Use("*", NotFoundHandler)
}

// NotFoundHandler func for describe 404 Error route.
func NotFoundHandler(c *fiber.Ctx) error {
	// Register new special route.
	// Return HTTP 404 status and JSON response.
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": true,
		"msg":   "sorry, endpoint is not found",
	})
}

package main

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/notifier"
	"anti-shoplifting/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	database.ConnectDynamoDB()
	database.AutoMigrate()
	notifier.InitializeFirebaseAdminSDK()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")
}

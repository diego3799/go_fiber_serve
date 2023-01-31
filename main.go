package main

import (
	"fun_server/db"
	"fun_server/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}
	db.InitDB()
	PORT := os.Getenv("PORT")
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")

	userRoutes := api.Group("/users")
	routes.UserRoutes(&userRoutes)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Working!!!")
	})
	app.Listen(PORT)
}

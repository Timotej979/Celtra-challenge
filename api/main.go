package main

import (
	//"github.com/Timotej979/Celtra-challenge/api/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Send a string back for GET calls to the endpoint "/"
	app.Get("/healthz", func(c *fiber.Ctx) error {
		err := c.SendString("API is running!")
		return err
	})

	app.Listen(":3000")
}

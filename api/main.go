package main

import (
	// Globally available packages
	"log"

	"github.com/gofiber/fiber/v2"

	// Locally available packages
	"github.com/Timotej979/Celtra-challenge/api/config"
)

func main() {

	// Get the environment variables
	envVars, err := config.GetEnvVars()
	if err != nil {
		log.Fatalf("Error getting environment variables: %v", err)
	}

	// Extract the port from cmd

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		err := c.SendString("API is running!")
		return err
	})

	app.Listen(":3000")
}

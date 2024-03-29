package main

import (
	// Globally available packages
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"

	// Locally available packages
	"github.com/Timotej979/Celtra-challenge/api/config"
)

func main() {

	// Set up the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get the environment variables
	envVars, err := config.GetEnvVars()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting environment variables")
	}

	// Print the environment variables
	log.Info().Interface("envVars", envVars).Msg("environment variables")

	// Extract the port from cmd

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		err := c.SendString("API is running!")
		return err
	})

	app.Listen(":3000")
}

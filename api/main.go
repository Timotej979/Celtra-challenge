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

	// Get the environment variables
	envVars, err := config.GetEnvVars()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting environment variables")
	}

	// Setup the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	// Set the log level
	switch envVars.AppConfig {
	case "dev":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "prod":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		log.Fatal().Str("appConfig", envVars.AppConfig).Msg("invalid app config: must be 'dev' or 'prod'")
		os.Exit(1)
	}

	// Log the environment variables
	log.Info().Interface("envVars", envVars).Msg("environment variables")

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		err := c.SendString("API is running!")
		return err
	})

	app.Listen(":3000")
}

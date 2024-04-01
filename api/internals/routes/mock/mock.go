package mockRoutes

import (
	"github.com/rs/zerolog"

	"github.com/gofiber/fiber/v2"

	"github.com/Timotej979/Celtra-challenge/api/dal"
	mockHandler "github.com/Timotej979/Celtra-challenge/api/internals/handlers/mock"
)

// SetupRoutes sets up the routes for the Instagram API
func SetupRoutes(router fiber.Router, dalConfig dal.DALConfig, defaultLogger zerolog.Logger) {
	// Create the Instagram handler
	mock := router.Group("/mock")

	// Initialize the Instagram handler
	mockHandlerInstance := mockHandler.NewMockHandler(dalConfig, defaultLogger)

	// Health check of the Instagram API
	mock.Get("/healthz", mockHandlerInstance.Healthz)

	// Get the Instagram user data
	mock.Get("/:accountID/data", mockHandlerInstance.GetRandomAccountData)

}

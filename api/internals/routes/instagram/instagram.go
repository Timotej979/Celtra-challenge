package instagramRoutes

import (
	"github.com/rs/zerolog"

	"github.com/gofiber/fiber/v2"

	"github.com/Timotej979/Celtra-challenge/api/dal"
	instagramHandler "github.com/Timotej979/Celtra-challenge/api/internals/handlers/instagram"
)

// SetupRoutes sets up the routes for the Instagram API
func SetupRoutes(router fiber.Router, dal *dal.DAL, defaultLogger zerolog.Logger) {
	// Create the Instagram handler
	instagram := router.Group("/instagram")

	// Initialize the Instagram handler
	instagramHandlerInstance := instagramHandler.NewInstagramHandler(dal, defaultLogger)

	// Health check of the Instagram API
	instagram.Get("/healthz", instagramHandlerInstance.Healthz())

	// Get the Instagram user data
	instagram.Get("/:accountID/description", instagramHandlerInstance.GetInstagramUserDescription())

}

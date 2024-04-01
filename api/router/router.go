package router

import (
	"github.com/rs/zerolog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"github.com/Timotej979/Celtra-challenge/api/dal"
	mockRoutes "github.com/Timotej979/Celtra-challenge/api/internals/routes/mock"
)

func SetupRouter(app *fiber.App, dalConfig dal.DALConfig, defaultLogger zerolog.Logger) {

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/mock-api/v1", logger.New())

	mockRoutes.SetupRoutes(api, dalConfig, defaultLogger)
}

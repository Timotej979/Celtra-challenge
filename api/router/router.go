package router

import (
	"github.com/rs/zerolog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"github.com/Timotej979/Celtra-challenge/api/dal"
	"github.com/Timotej979/Celtra-challenge/api/internals/routes/instagram"
)

func SetupRouter(app *fiber.App, dal *dal.DAL, defaultLogger zerolog.Logger) {

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/social-media-api/v1", logger.New())

	instagram.SetupRoutes(api, dal, defaultLogger)
}

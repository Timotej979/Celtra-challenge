package mockHandlers

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"

	"github.com/Timotej979/Celtra-challenge/api/dal"
)

type CoincapHandler struct {
	// DAL is the Data Access Layer
	dal *dal.DAL
}

func NewMockHandler(dalConfig dal.DALConfig, logger zerolog.Logger) *CoincapHandler {
	// Assign the logger to the handler
	log.Logger = logger

	// Create the DAL
	dalInstance, err := dal.NewDAL(&dalConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating the DAL")
	}

	// Connect to the database
	err = dalInstance.DbDriver.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to the database")
	}

	// Migrate the database
	err = dalInstance.DbDriver.Migrate()
	if err != nil {
		log.Fatal().Err(err).Msg("error migrating the database")
	}

	return &CoincapHandler{
		dal: dalInstance,
	}
}

// Healthz is a handler for the Instagram API health check
func (h *CoincapHandler) Healthz(c *fiber.Ctx) error {
	log.Info().Msg("Checking mock API health...")
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *CoincapHandler) GetRandomAccountData(c *fiber.Ctx) error {
	log.Info().Msg("Getting random account data...")

	// Get the accountID from the URL
	accountID := c.Params("accountID")

	// Generate data from string
	var data string
	if len(accountID)%2 == 0 {
		data = "even"
	} else {
		data = "odd"
	}

	// Write to the DB through DAL
	err := h.dal.DbDriver.InsertUserData(accountID, data)
	if err != nil {
		log.Error().Err(err).Msg("error writing to the database")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "error writing to the database"})
	}

	return c.JSON(fiber.Map{"status": "ok", "accountID": accountID, "timestamp": time.Now().Unix(), "data": data})
}

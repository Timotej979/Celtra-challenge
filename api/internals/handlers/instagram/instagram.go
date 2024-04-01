package instagramHandlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"

	"github.com/Timotej979/Celtra-challenge/api/dal"
)

type InstagramHandler struct {
	dal *dal.DAL
}

func NewInstagramHandler(dal *dal.DAL, logger zerolog.Logger) *InstagramHandler {
	// Assign the logger to the handler
	log.Logger = logger

	return &InstagramHandler{
		dal: dal,
	}
}

// Healthz is a handler for the Instagram API health check
func (h *InstagramHandler) Healthz(c *fiber.Ctx) error {
	log.Info().Msg("Instagram API health check")
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *InstagramHandler) GetInstagramUserDescription(c *fiber.Ctx) error {
	// Get the account ID from the URL
	accountID := c.Params("accountID")

	// Make a web request to the Instagram API like so: https://www.instagram.com/leomessi/?__a=1&__d=dis
	// Make a web request to the Instagram API
	reqURL := fmt.Sprintf("https://www.instagram.com/%s/?__a=1&__d=dis", accountID)
	resp, err := http.Get(reqURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the response body
	data := string(body)
	log.Info().Msg("Instagram API response: " + data)

	// TODO: Parse the Instagram user data

	// TODO: Save the Instagram user data to the database

	// TODO: Return the Instagram user data in JSON format

	return c.JSON(fiber.Map{"data": data})
}

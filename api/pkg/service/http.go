package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type HttpService struct {
	app *fiber.App
}

func NewHttpService() *HttpService {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New())

	app.Get("/api/v1/status", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})

	return &HttpService{
		app: app,
	}
}

func (h *HttpService) Run(addr string) error {
	return h.app.Listen(addr)
}

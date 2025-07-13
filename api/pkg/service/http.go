package service

import (
	"github.com/agelito/rinha-de-backend-2025/api/pkg/handler"
	"github.com/agelito/rinha-de-backend-2025/api/pkg/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

type HttpService struct {
	app *fiber.App
}

func NewHttpService(payments *handler.PaymentsHandler) *HttpService {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New())

	app.Use(healthcheck.New())

	app.Post("/payments", func(c *fiber.Ctx) error {
		var payment model.Payment

		if err := c.BodyParser(&payment); err != nil {
			log.Error("invalid request body", "error", err)
			return c.Status(fiber.StatusBadRequest).SendString("invalid payment")
		}

		if err := payments.Payment(&payment); err != nil {
			log.Error("error handling payment", "error", err)

			switch err {
			case handler.PaymentFailed:
				return c.SendStatus(fiber.StatusUnprocessableEntity)
			case handler.PaymentTimeout:
				return c.SendStatus(fiber.StatusGatewayTimeout)
			}

			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	app.Get("/payments-summary", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("{}")
	})

	return &HttpService{
		app: app,
	}
}

func (h *HttpService) Run(addr string) error {
	return h.app.Listen(addr)
}

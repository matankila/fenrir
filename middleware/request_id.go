package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewUuidMid() fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := uuid.New()
		c.Request().Header.Set(fiber.HeaderXRequestID, u.String())
		return c.Next()
	}
}

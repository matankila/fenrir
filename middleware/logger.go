package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matankila/fenrir/config"
	"go.uber.org/zap"
)

func NewLoggingMid(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := config.RequestInfo{
			Method: c.Method(),
			Url:    string(c.Request().RequestURI()),
			Ip:     c.IP(),
		}

		logger.Info("start",
			             zap.Any("requestInfo", req),
			             zap.String("uid", c.Get(fiber.HeaderXRequestID)))
		defer logger.Info("finish",
								zap.Any("requestInfo", req),
								zap.String("uid", c.Get(fiber.HeaderXRequestID)))

		return c.Next()
	}
}

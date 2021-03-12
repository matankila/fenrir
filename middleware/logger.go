package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/logger"
	"go.uber.org/zap"
)

func NewLoggingMid() fiber.Handler {
	return func(c *fiber.Ctx) error {
		l := logger.GetLogger(logger.Default)
		req := config.RequestInfo{
			Method: c.Method(),
			Url:    string(c.Request().RequestURI()),
			Ip:     c.IP(),
		}

		l.Info("start",
			zap.Any("requestInfo", req),
			zap.String("uid", c.Get(fiber.HeaderXRequestID)))
		defer l.Info("finish",
			zap.Any("requestInfo", req),
			zap.String("uid", c.Get(fiber.HeaderXRequestID)))

		return c.Next()
	}
}

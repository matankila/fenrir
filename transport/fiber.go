package transport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/endpoints"
	"github.com/matankila/fenrir/logger"
	"github.com/matankila/fenrir/middleware"
	"go.uber.org/zap"
	"sync"
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			return &config.RequestInfo{}
		},
	}
)

func InitApp(ep endpoints.Endpoints) *fiber.App {
	l := logger.GetLogger(logger.Default)
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			req := pool.Get().(*config.RequestInfo)
			req.Method = c.Method()
			req.Ip = c.IP()
			req.Url = string(c.Request().RequestURI())
			l.Error(err.Error(),
				zap.String("uid", c.Get(fiber.HeaderXRequestID)),
				zap.Any("requestInfo", req))
			pool.Put(req)
			return c.Status(code).JSON(fiber.Map{"error": true, "msg": err.Error()})
		},
	})
	app.Use(middleware.NewUuidMid())
	app.Use(middleware.NewLoggingMid())
	app.Post("/validate", ep.Validate)
	app.Get("/health", ep.Health)

	return app
}

package main

import (
	"fmt"
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/endpoints"
	"github.com/matankila/fenrir/service"
	"github.com/matankila/fenrir/transport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

// set logger
func initLogger() *zap.Logger {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(config.LogLvl)); err != nil {
		panic(err)
	}

	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	c := zap.NewProductionConfig()
	c.Level = lvl
	c.InitialFields = map[string]interface{}{"loggerName": config.LoggerName}
	c.OutputPaths = []string{config.Output}
	c.EncoderConfig = ec
	logger, err := c.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func main() {
	logger := initLogger()
	defer logger.Sync()

	errs := make(chan error)
	// handle signals
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <- c)
	}()
	// server
	go func() {
		s := service.NewService(logger)
		ep := endpoints.MakeEndpoints(s)
		app := transport.InitApp(ep, logger)
		errs <- app.Listen(":" + config.Port)
	}()

	logger.Error((<-errs).Error())
}
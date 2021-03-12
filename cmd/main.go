package main

import (
	"fmt"
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/endpoints"
	"github.com/matankila/fenrir/logger"
	"github.com/matankila/fenrir/reloader"
	"github.com/matankila/fenrir/service"
	"github.com/matankila/fenrir/transport"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// done signals all loggers to be synced/ flushed
func handleGracefulShutdown(done chan struct{}) {
	close(done)
	time.Sleep(5 * time.Second)
}

func main() {
	errs := make(chan error)
	done := logger.InitFactory()
	l := logger.GetLogger(logger.Default)
	defer handleGracefulShutdown(done)

	// handle config loader, and reload it on change
	r := reloader.New(&config.FallBackConf, config.ConfigPolicyPath)
	r.Run()

	// handle signals
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// server
	go func() {
		s := service.NewService()
		ep := endpoints.MakeEndpoints(s)
		app := transport.InitApp(ep)
		errs <- app.Listen(":" + config.Port)
	}()

	l.Error((<-errs).Error())
}

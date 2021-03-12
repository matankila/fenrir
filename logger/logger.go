package logger

import (
	"github.com/matankila/fenrir/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerFactory map[string]*zap.Logger
type Logger interface {
	String() string
}
type WatcherLogger struct{}
type DefaultLogger struct{}

var (
	l       = loggerFactory{}
	Default = DefaultLogger{}
	Watcher = WatcherLogger{}
	done    = make(chan struct{})
)

func (d DefaultLogger) String() string {
	return "default"
}

func (w WatcherLogger) String() string {
	return "watcher"
}

func initLogger(loggerName string) *zap.Logger {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(config.LogLvl)); err != nil {
		panic(err)
	}

	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	c := zap.NewProductionConfig()
	c.Level = lvl
	c.InitialFields = map[string]interface{}{"loggerName": loggerName}
	c.OutputPaths = []string{config.Output}
	c.EncoderConfig = ec
	logger, err := c.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func InitFactory() chan struct{} {
	l[Default.String()] = initLogger(config.DefaultLoggerName)
	l[Watcher.String()] = initLogger(config.WatcherLoggerName)

	// waits for channel to be closed and sync all loggers
	go func() {
		<-done
		for _, v := range l {
			v.Sync()
		}
	}()

	return done
}

func GetLogger(logger Logger) *zap.Logger {
	return l[logger.String()]
}

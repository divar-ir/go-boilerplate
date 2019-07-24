package main

import (
	"fmt"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

// LoggingConfig the loggers's configuration structure
type LoggingConfig struct {
	Level         string
	SentryEnabled bool
}

const (
	sentryDSN = "http://xxxx:xxxx@sentry.cafebazaar.ir/112"
)

// ConfigureLogging handlerconfig logger based on the given configuration
func provideLogger(config *Config) (*logrus.Logger, error) {
	logger := logrus.New()
	if config.Logging.Level != "" {
		level, err := logrus.ParseLevel(config.Logging.Level)
		if err != nil {
			return nil, err
		}
		logger.SetLevel(level)
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})

	if config.Logging.SentryEnabled {

		hook, err := logrus_sentry.NewAsyncSentryHook(sentryDSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

		if err != nil {
			fmt.Println(err)
			panic("failed to setup raven!")
		}

		hook.StacktraceConfiguration.Enable = true

		logger.AddHook(hook)
	}

	return logger, nil
}

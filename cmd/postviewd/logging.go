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
func configureLogging(config *LoggingConfig) error {
	if config.Level != "" {
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return err
		}
		logrus.SetLevel(level)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})

	if config.SentryEnabled {

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

		logrus.AddHook(hook)
	}

	return nil
}

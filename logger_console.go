package kimg

import (
	"log"
	"os"
)

type kimgConsoleLogger struct {
	*KimgBaseLogger
}

// NewKimgConsoleLogger create a console logger instance.
func NewKimgConsoleLogger(config *KimgConfig) (KimgLogger, error) {
	return &kimgConsoleLogger{
		KimgBaseLogger: &KimgBaseLogger{
			log:   log.New(os.Stdout, "", log.LstdFlags),
			level: config.Logger.Level,
		},
	}, nil
}

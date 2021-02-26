package kimg

import (
	"log"
	"os"
)

type kimgFileLogger struct {
	logFile *os.File
	*KimgBaseLogger
}

// NewKimgFileLogger create a file based logger instance.
func NewKimgFileLogger(config *KimgConfig) (KimgLogger, error) {
	logFile, err := os.OpenFile(config.Logger.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &kimgFileLogger{
		logFile: logFile,
		KimgBaseLogger: &KimgBaseLogger{
			log:   log.New(logFile, "", log.LstdFlags),
			level: config.Logger.Level,
		},
	}, nil
}

package kimg

import (
	"fmt"
	"log"
)

// KimgLogger is a interface to provide logger in kimg.
type KimgLogger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// KimgBaseLogger base logger struct hold golang logger.
type KimgBaseLogger struct {
	log   *log.Logger
	level string
}

// NewKimgLogger create a logger instance according to logger mode in config.
func NewKimgLogger(config *KimgConfig) (KimgLogger, error) {
	switch config.Logger.Mode {
	case "console":
		log.Println("[INFO] logger [console] used")
		return NewKimgConsoleLogger(config)
	case "file":
		log.Println("[INFO] logger [file] used")
		return NewKimgFileLogger(config)
	default:
		log.Printf("unsupported logger mode :%s\n", config.Logger.Mode)
		return nil, nil
	}
}

// Debug log
func (logger *KimgBaseLogger) Debug(format string, v ...interface{}) {
	if logger.level == "debug" {
		logger.log.Output(2, fmt.Sprintf("[DEBUG] %s", fmt.Sprintf(format, v...)))
	}
}

// Info log
func (logger *KimgBaseLogger) Info(format string, v ...interface{}) {
	if logger.level == "debug" || logger.level == "info" {
		logger.log.Output(2, fmt.Sprintf("[INFO] %s", fmt.Sprintf(format, v...)))
	}
}

// Warn log
func (logger *KimgBaseLogger) Warn(format string, v ...interface{}) {
	if logger.level != "error" {
		logger.log.Output(2, fmt.Sprintf("[WARN] %s", fmt.Sprintf(format, v...)))
	}
}

// Error log
func (logger *KimgBaseLogger) Error(format string, v ...interface{}) {
	if logger.level == "error" {
		logger.log.Output(2, fmt.Sprintf("[ERROR] %s", fmt.Sprintf(format, v...)))
	}
}

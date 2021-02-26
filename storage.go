package kimg

import (
	"errors"
	"log"
)

// KimgStorage is a interface to provide storage in kimg.
type KimgStorage interface {
	Set(req *KimgRequest, data []byte) error
	Get(req *KimgRequest) ([]byte, error)
	Del(req *KimgRequest) error
}

// KimgBaseStorage base storage struct hold kimg context.
type KimgBaseStorage struct {
	ctx *KimgContext
}

// Debug log
func (storage *KimgBaseStorage) Debug(format string, v ...interface{}) {
	storage.ctx.Logger.Debug(format, v...)
}

// Info log
func (storage *KimgBaseStorage) Info(format string, v ...interface{}) {
	storage.ctx.Logger.Info(format, v...)
}

// Warn log
func (storage *KimgBaseStorage) Warn(format string, v ...interface{}) {
	storage.ctx.Logger.Warn(format, v...)
}

// Error log
func (storage *KimgBaseStorage) Error(format string, v ...interface{}) {
	storage.ctx.Logger.Error(format, v...)
}

// NewKimgStorage create a storage instance according to storage mode in config.
func NewKimgStorage(ctx *KimgContext) (KimgStorage, error) {
	switch ctx.Config.Storage.Mode {
	case "file":
		log.Println("[INFO] storage [file] used")
		return NewKimgFileStorage(ctx)
	default:
		log.Printf("[WARN] unsupported storage mode :%s\n", ctx.Config.Storage.Mode)
		return nil, errors.New("No Available Storage")
	}
}

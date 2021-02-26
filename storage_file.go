package kimg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type kimgFileStorage struct {
	rootDir string
	mtx     sync.RWMutex
	*KimgBaseStorage
}

// NewKimgFileStorage create a file based storage instance.
func NewKimgFileStorage(ctx *KimgContext) (KimgStorage, error) {
	return &kimgFileStorage{
		rootDir: ctx.Config.Storage.Root,
		KimgBaseStorage: &KimgBaseStorage{
			ctx: ctx,
		},
	}, nil
}

func (storage *kimgFileStorage) Set(req *KimgRequest, data []byte) error {
	imageDir, imageFile := storage.imageDirAndFile(req)

	storage.mtx.Lock()
	defer storage.mtx.Unlock()

	err := os.MkdirAll(imageDir, 0755)
	if err != nil {
		storage.Warn("MkdirAll %s, err: %s", imageDir, err)
		return err
	}

	err = ioutil.WriteFile(imageFile, data, 0644)
	if err != nil {
		storage.Warn("WriteFile %s, err: %s", imageFile, err)
		return err
	}

	storage.Debug("kimgFileStorage Set dir:%s, file: %s, size: %d", imageDir, imageFile, len(data))

	return nil
}

func (storage *kimgFileStorage) Get(req *KimgRequest) ([]byte, error) {
	_, imageFile := storage.imageDirAndFile(req)

	storage.mtx.RLock()
	defer storage.mtx.RUnlock()

	data, err := ioutil.ReadFile(imageFile)
	if err != nil {
		storage.Warn("ReadFile %s, err: %s", imageFile, err)
		return nil, err
	}

	storage.Debug("kimgFileStorage Get file: %s, size: %d", imageFile, len(data))

	return data, nil
}

func (storage *kimgFileStorage) Del(req *KimgRequest) error {
	imageDir, _ := storage.imageDirAndFile(req)

	storage.mtx.Lock()
	defer storage.mtx.Unlock()

	err := os.RemoveAll(imageDir)
	if err != nil {
		storage.Warn("RemoveAll %s, err: %s", imageDir, err)
		return err
	}

	storage.Debug("kimgFileStorage Del dir: %s", imageDir)

	return nil
}

func (storage *kimgFileStorage) imageDirAndFile(req *KimgRequest) (string, string) {
	l1, _ := strconv.ParseUint(req.Md5[:3], 16, 0)
	l2, _ := strconv.ParseUint(req.Md5[3:6], 16, 0)
	d1 := strconv.FormatUint(l1/4, 10)
	d2 := strconv.FormatUint(l2/4, 10)

	imageDir := filepath.Join(storage.rootDir, d1, d2, req.Md5)
	imageFile := ""

	if req.Origin {
		imageFile = filepath.Join(imageDir, "origin")
	} else {
		imageFile = filepath.Join(imageDir, req.Key())
	}

	return imageDir, imageFile
}

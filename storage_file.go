package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type kimgFileStorage struct {
	rootDir string
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

func (storage *kimgFileStorage) Release() {}

func (storage *kimgFileStorage) Set(req *KimgRequest, data []byte) error {
	imageDir, imageFile := storage.imageDirAndFile(req)

	err := os.MkdirAll(imageDir, 0755)
	if err != nil {
		storage.Warn("MkdirAll %s, err: %s", imageDir, err)
		return err
	}

	size := len(data)
	f, err := os.OpenFile(imageFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		storage.Warn("OpenFile %s, err: %s", imageFile, err)
		return err
	}
	defer f.Close()
	if n, err := f.Write(data); err != nil {
		storage.Warn("Write %s, err: %s", imageFile, err)
		return err
	} else if n < size {
		storage.Warn("Write %s, wrote: %d, need: %d", imageFile, n, size)
		return io.ErrShortWrite
	}

	storage.Debug("kimgFileStorage Set dir:%s, file: %s, size: %d", imageDir, imageFile, size)

	return err
}

func (storage *kimgFileStorage) Get(req *KimgRequest) ([]byte, error) {
	_, imageFile := storage.imageDirAndFile(req)

	f, err := os.OpenFile(imageFile, os.O_RDONLY, 0755)
	defer f.Close()
	if err != nil {
		storage.Debug("OpenFile %s, err: %s", imageFile, err)
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		storage.Warn("Stat %s, err: %s", imageFile, err)
		return nil, err
	}
	size := int(fi.Size())
	if size <= 0 {
		storage.Warn("Stat %s, size is 0", imageFile)
		return nil, io.ErrShortBuffer
	}

	data := make([]byte, size)
	if n, err := f.Read(data); err != nil {
		storage.Warn("Read %s, err: %s", imageFile, err)
		return nil, err
	} else if n < size {
		storage.Warn("Read %s, read: %d, need: %d", imageFile, n, size)
		return nil, io.ErrShortBuffer
	}

	storage.Debug("kimgFileStorage Get file: %s, size: %d", imageFile, size)

	return data, nil
}

func (storage *kimgFileStorage) Del(req *KimgRequest) error {
	imageDir, _ := storage.imageDirAndFile(req)
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

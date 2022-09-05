package kimg

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type kimgMinioStorage struct {
	client *minio.Client
	bucket string
	*KimgBaseStorage
}

// NewKimgMinioStorage create a minio based storage instance.
func NewKimgMinioStorage(ctx *KimgContext) (KimgStorage, error) {
	client, err := minio.New(ctx.Config.Storage.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ctx.Config.Storage.Minio.AccessKeyID, ctx.Config.Storage.Minio.SecretAccessKey, ""),
		Secure: ctx.Config.Storage.Minio.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	context := context.Background()
	err = client.MakeBucket(context, ctx.Config.Storage.Minio.Bucket, minio.MakeBucketOptions{ObjectLocking: false})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context, ctx.Config.Storage.Minio.Bucket)
		if !exists || errBucketExists != nil {
			return nil, err
		}
	}

	return &kimgMinioStorage{
		client: client,
		bucket: ctx.Config.Storage.Minio.Bucket,
		KimgBaseStorage: &KimgBaseStorage{
			ctx: ctx,
		},
	}, nil
}

func (storage *kimgMinioStorage) Set(req *KimgRequest, data []byte) error {
	imageDir, imageFile := storage.imageDirAndFile(req)

	_, err := storage.client.PutObject(context.Background(), storage.bucket, imageFile, bytes.NewReader(data), -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	storage.Debug("kimgMinioStorage Set dir:%s, file: %s, size: %d", imageDir, imageFile, len(data))

	return nil
}

func (storage *kimgMinioStorage) Get(req *KimgRequest) ([]byte, error) {
	_, imageFile := storage.imageDirAndFile(req)

	object, err := storage.client.GetObject(context.Background(), storage.bucket, imageFile, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	data, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}

	storage.Debug("kimgMinioStorage Get file: %s, size: %d", imageFile, len(data))

	return data, nil
}

func (storage *kimgMinioStorage) Del(req *KimgRequest) error {
	imageDir, _ := storage.imageDirAndFile(req)

	err := storage.client.RemoveObject(context.Background(), storage.bucket, imageDir, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	storage.Debug("kimgMinioStorage Del dir: %s", imageDir)

	return nil
}

func (storage *kimgMinioStorage) imageDirAndFile(req *KimgRequest) (string, string) {
	imageDir := req.Md5
	imageFile := ""

	if req.Origin {
		imageFile = filepath.Join(imageDir, "origin")
	} else {
		imageFile = filepath.Join(imageDir, req.Key())
	}

	return imageDir, imageFile
}

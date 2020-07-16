package minio

import (
	"github.com/minio/minio-go/v6"
	"github.com/wmsx/store_api/setting"
	"io"
)

const (
	AvatarBulk = "avatar"
)

var (
	minioClient *minio.Client
)

func SetUp(minioSetting setting.Minio) (err error) {
	if minioClient, err = minio.New(
		minioSetting.Endpoint,
		minioSetting.AccessKey,
		minioSetting.SecretAccessKey,
		false); err != nil {
		return
	}
	return nil
}

func UploadFile(bulk, objectName string, size int64, reader io.Reader) (err error) {
	_, err = minioClient.PutObject(bulk, objectName, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

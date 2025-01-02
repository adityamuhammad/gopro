package config

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinIOClient *minio.Client

func InitMinIO() error {
	var err error
	MinIOClient, err = minio.New("minio-gopro:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	return err
}

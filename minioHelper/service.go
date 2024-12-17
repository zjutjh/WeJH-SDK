package minioHelper

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
)

// PutObject 用于上传对象
func (ms *Service) PutObject(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}
	_, err := (*ms).Client.PutObject(context.Background(), (*ms).Bucket, objectKey, reader, size, opts)
	if err != nil {
		return "", err
	}
	return (*ms).Domain + (*ms).Bucket + "/" + objectKey, nil
}

// GetObjectKeyFromUrl 从 Url 中提取 ObjectKey
// 若该 Url 不是来自此 Minio, 则 ok 为 false
func (ms *Service) GetObjectKeyFromUrl(fullUrl string) (objectKey string, ok bool) {
	objectKey = strings.TrimPrefix(fullUrl, (*ms).Domain+(*ms).Bucket+"/")
	if objectKey == fullUrl {
		return "", false
	}
	return objectKey, true
}

// DeleteObject 用于删除相应对象
func (ms *Service) DeleteObject(objectKey string) error {
	err := (*ms).Client.RemoveObject(
		context.Background(),
		(*ms).Bucket,
		objectKey,
		minio.RemoveObjectOptions{ForceDelete: true},
	)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

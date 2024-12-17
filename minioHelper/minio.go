package minioHelper

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Service 服务结构体
type Service struct {
	Client  *minio.Client
	Bucket  string
	Domain  string
	TempDir string
}

// InfoConfig 对象服务 MinIO 配置
type InfoConfig struct {
	EndPoint  string
	AccessKey string
	SecretKey string
	Secure    bool
	Bucket    string
	Domain    string
	TempDir   string
}

// Init 创建并返回 MinIO 服务实例
func Init(config *InfoConfig) (*Service, error) {
	// 初始化 MinIO 客户端对象
	client, err := minio.New(config.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("minio initialization failed: %w", err)
	}

	return &Service{
		Client:  client,
		Bucket:  config.Bucket,
		Domain:  config.Domain,
		TempDir: config.TempDir,
	}, nil
}

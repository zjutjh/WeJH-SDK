package cube

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

// 定义错误变量
var (
	ErrBucketNotSpecified  = errors.New("bucket not specified")
	ErrRequestBizCodeNotOK = errors.New("request business code not OK")
	ErrHttpStatusCodeNotOK = errors.New("HTTP status code not OK")
)

// Config 定义 Cube 配置结构体
type Config struct {
	BaseURL    string `mapstructure:"base_url"`    // 基础 URL
	APIKey     string `mapstructure:"api_key"`     // 应用密钥
	BucketName string `mapstructure:"bucket_name"` // 存储桶名称
}

// Client 定义 Cube 客户端
type Client struct {
	conf   Config
	client *resty.Client
}

// UploadResponse 定义 Cube 文件上传响应体
type UploadResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ObjectKey string `json:"object_key"`
	} `json:"data"`
}

// DeleteResponse 定义 Cube 文件删除响应体
type DeleteResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// New 创建一个 CubeClient
func New(conf Config) *Client {
	client := resty.New().
		SetBaseURL(strings.TrimRight(conf.BaseURL, "/")).
		SetHeader("Key", conf.APIKey)

	return &Client{
		conf:   conf,
		client: client,
	}
}

// UploadFile 上传文件到存储立方
func (c *Client) UploadFile(filename string, reader io.Reader, location string, convertWebp, useUUID bool) (*UploadResponse, error) {
	bucketName := c.conf.BucketName
	if bucketName == "" {
		return nil, ErrBucketNotSpecified
	}

	form := map[string]string{
		"bucket":       bucketName,
		"location":     location,
		"convert_webp": strconv.FormatBool(convertWebp),
		"use_uuid":     strconv.FormatBool(useUUID),
	}

	var result UploadResponse
	resp, err := c.client.R().
		SetFileReader("file", filename, reader).
		SetFormData(form).
		SetResult(&result).
		Post("/api/upload")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", ErrHttpStatusCodeNotOK, resp.StatusCode())
	}
	if result.Code != 200 {
		return &result, fmt.Errorf("%w: %v", ErrRequestBizCodeNotOK, result.Code)
	}
	return &result, nil
}

// DeleteFile 删除文件
func (c *Client) DeleteFile(objectKey string) (*DeleteResponse, error) {
	bucketName := c.conf.BucketName
	if bucketName == "" {
		return nil, ErrBucketNotSpecified
	}

	var result DeleteResponse
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"bucket":     bucketName,
			"object_key": objectKey,
		}).
		SetResult(&result).
		Delete("/api/delete")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", ErrHttpStatusCodeNotOK, resp.StatusCode())
	}
	if result.Code != 200 {
		return &result, fmt.Errorf("%w: %v", ErrRequestBizCodeNotOK, result.Code)
	}
	return &result, nil
}

// GetFileURL 返回文件访问 URL
func (c *Client) GetFileURL(objectKey string, thumbnail bool) string {
	baseURL := strings.TrimRight(c.conf.BaseURL, "/")
	params := url.Values{}
	params.Add("bucket", c.conf.BucketName)
	params.Add("object_key", objectKey)
	params.Add("thumbnail", strconv.FormatBool(thumbnail))
	return fmt.Sprintf("%s/api/file?%s", baseURL, params.Encode())
}

// GetObjectKeyFromUrl 从 URL 解析出 objectKey
func (c *Client) GetObjectKeyFromUrl(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	queryParams := u.Query()
	return queryParams.Get("object_key")
}

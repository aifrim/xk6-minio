package xk6_minio

import (
	"context"
	"encoding/json"
	"log"

	"go.k6.io/k6/js/modules"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Transforms input from goja runtime into go structs
// It uses json marshaling
func bridge(from interface{}, out any) error {
	bytes, err := json.Marshal(&from)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &out)

	if err != nil {
		return err
	}

	return nil
}

// Register the extensions on module initialization.
func init() {
	modules.Register("k6/x/minio", New())
}

type Minio struct{}

type Client struct {
	minioClient *minio.Client
}

func New() *Minio {
	return &Minio{}
}

func (*Minio) NewClient(endpoint string, keyId string, accessKey string, secure bool, region string) interface{} {
	log.Println("use ssl", secure)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(keyId, accessKey, ""),
		Secure: secure,
		Region: region,
	})

	if err != nil {
		return err
	}

	return &Client{minioClient: client}
}

func (client *Client) BucketExists(bucketName string) interface{} {
	ctx := context.Background()

	exists, err := client.minioClient.BucketExists(ctx, bucketName)

	if err != nil {
		return err
	}

	return exists
}

// type MinioGetObjectOptions struct {
// 	headers   map[string]string   `json:"headers"`
// 	reqParams map[string][]string `json:"reqParams"`
// }

// func (client *Client) StatObject(bucketName string, objectName string, options interface{}) interface{} {
// 	ctx := context.Background()

// 	exists, err := client.minioClient.StatObject(ctx, bucketName, objectName, opts)

// 	if err != nil {
// 		return err
// 	}

// 	return exists
// }

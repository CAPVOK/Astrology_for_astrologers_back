package minioclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient представляет клиент MinIO.
type MinioClient struct {
	Client *minio.Client
}

// NewMinioClient создает новый экземпляр клиента MinIO.
func NewMinioClient() (*MinioClient, error) {
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	endpoint := os.Getenv("MINIO_ENDPOINT")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinioClient{
		Client: client,
	}, nil
}

// UploadServiceImage загружает изображение в MinIO и возвращает URL изображения.
func (mc *MinioClient) UploadServiceImage(planetId int, imageBytes []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("planets/%d/image", planetId)

	// Используйте io.NopCloser вместо ioutil.NopCloser
	reader := io.NopCloser(bytes.NewReader(imageBytes))

	_, err := mc.Client.PutObject(context.TODO(), "space-images", objectName, reader, int64(len(imageBytes)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Формирование URL изображения
	imageURL := fmt.Sprintf("http://localhost:9000/space-images/%s", objectName)
	return imageURL, nil
}

// RemoveServiceImage удаляет изображение услуги из MinIO.
func (mc *MinioClient) RemoveServiceImage(planetId int) error {
	objectName := fmt.Sprintf("planets/%d/image", planetId)
	err := mc.Client.RemoveObject(context.TODO(), "space-images", objectName, minio.RemoveObjectOptions{})
	if err != nil {
		fmt.Println("Failed to remove object from MinIO:", err)
		// Обработка ошибки удаления изображения из MinIO
		return err
	}
	fmt.Println("Object removed from MinIO successfully:", objectName)
	return nil
}

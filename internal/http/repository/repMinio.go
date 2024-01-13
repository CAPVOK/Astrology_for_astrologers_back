package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioRepository interface {
	UploadServiceImage(userID, planetID uint64, imageBytes []byte, contentType string) (string, error)
	RemoveServiceImage(userID, planetID uint64) error
}

func (r *Repository) UploadServiceImage(planetID, userID uint, imageBytes []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("planets/%d/image", planetID)
	reader := io.NopCloser(bytes.NewReader(imageBytes))
	_, err := r.mc.PutObject(context.TODO(), "space-images", objectName, reader, int64(len(imageBytes)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", errors.New("ошибка при добавлении изображения в минио бакет")
	}
	// Формирование URL изображения
	imageURL := fmt.Sprintf("http://localhost:9000/space-images/%s", objectName)
	return imageURL, nil
}

func (r *Repository) RemoveServiceImage(planetID, userID uint) error {
	objectName := fmt.Sprintf("planets/%d/image", planetID)
	err := r.mc.RemoveObject(context.TODO(), "space-images", objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.New("не удалось удалить изображение из бакет")
	}

	if err := r.db.Table("planets").
		Where("planet_id = ?", planetID).
		Update("image_name", nil).Error; err != nil {
		return errors.New("ошибка при обновлении URL изображения в базе данных")
	}

	return nil
}

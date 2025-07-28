package images

import "github.com/google/uuid"

type ImagesRepo interface {
	// CreateImage создает изображение
	CreateImage() ([]byte, error)
	// GetImage получает изображение
	GetImage(imageName string) ([]byte, error)
	// LoadImage загружает изображение
	LoadImage(mimeType string, image []byte, userID uuid.UUID) (*string, error)
}

package images

import (
	"fmt"

	"github.com/google/uuid"
)

// LoadImage загружает изображение
func (repo *ImagesDBRepository) LoadImage(mimeType string, image []byte, userID uuid.UUID) (*string, error) {
	name := fmt.Sprintf("image%s", uuid.New().String())
	query := `INSERT INTO images (name, mimetype, image_data, user_id)
	         VALUES ($1, $2, $3, $4) RETURNING name;`

	var imageName string
	err := repo.dtb.QueryRow(query, name, mimeType, image, userID.String()).Scan(&imageName)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: загрузка изображения: %v", err)
	}

	var loadedImage []byte
	err = repo.dtb.QueryRow("SELECT image_data FROM images WHERE name = $1;", imageName).Scan(&loadedImage)
	if err != nil {
		return nil, fmt.Errorf("error while selecting the image: %v", err)
	}
	return &imageName, nil
}

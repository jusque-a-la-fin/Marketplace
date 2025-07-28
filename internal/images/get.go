package images

import "fmt"

// GetImage получает изображение
func (repo *ImagesDBRepository) GetImage(imageName string) ([]byte, error) {
	var image []byte
	err := repo.dtb.QueryRow("SELECT image_data FROM images WHERE name = $1;", imageName).Scan(&image)
	if err != nil {
		return nil, fmt.Errorf("error while selecting the image: %v", err)
	}
	return image, nil
}

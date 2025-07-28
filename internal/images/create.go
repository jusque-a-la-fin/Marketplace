package images

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
)

const (
	// maxImageBytes — максимальный размер файла изображения в байтах (2 Мбайта)
	MaxImageBytes int = 2000000
	// minImageDim — минимальное количество пикселей по каждой стороне изображения
	MinImageDim int = 500
	// maxImageDim — максимальное количество пикселей по каждой стороне изображения
	MaxImageDim int = 2000
	// minAspectRatio — минимально допустимое соотношение ширины к высоте изображения
	MinAspectRatio float64 = 0.8
	// maxAspectRatio — максимально допустимое соотношение ширины к высоте изображения
	MaxAspectRatio float64 = 1.2
	// maxAttempts — максимальное количество попыток создать изображение
	maxAttempts int = 10
)

// CreateImage создает изображение
func (repo *ImagesDBRepository) CreateImage() ([]byte, error) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		var width, height int
		for {
			width = rand.Intn(MaxImageDim-MinImageDim+1) + MinImageDim
			height = rand.Intn(MaxImageDim-MinImageDim+1) + MinImageDim
			ratio := float64(width) / float64(height)
			if ratio >= MinAspectRatio && ratio <= MaxAspectRatio {
				break
			}
		}

		randColor := color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 0xFF,
		}

		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, randColor)
			}
		}

		var buf bytes.Buffer
		opts := &jpeg.Options{Quality: 95}
		if err := jpeg.Encode(&buf, img, opts); err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %v", err)
		}

		size := buf.Len()
		if size <= MaxImageBytes {
			return buf.Bytes(), nil
		}

		fmt.Printf(
			"Попытка %d: размер (%d байтов) созданного изображения превышает максимально допустимый размер, повтор создания изображения...\n",
			attempt, size,
		)
	}

	return nil, fmt.Errorf(
		"не получилось создать изображение в течение %d попыток",
		maxAttempts,
	)
}

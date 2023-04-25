package services

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/imaging"
)

type ImageService struct {
	encodingService ImageEncodingService
}

func NewImageService(encodingService ImageEncodingService) *ImageService {
	return &ImageService{
		encodingService: encodingService,
	}
}

type ImageOptions struct {
	Width   int
	Quality int
	Format  string
}

func (s *ImageService) ProcessImage(img image.Image, options ImageOptions) ([]byte, error) {
	if options.Width > 0 {
		img = s.resizeImage(img, options.Width)
	}

	encoder, err := s.encodingService.NewEncoder(options.Format, EncoderOptions{Quality: options.Quality})
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}

	imgBytes, err := encoder.Encode(img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return imgBytes, nil
}

func (s *ImageService) resizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	maxWidth := bounds.Max.X

	if width > maxWidth {
		width = maxWidth
	}

	img = imaging.Resize(img, width, 0, imaging.Lanczos)

	return img
}

func (s *ImageService) LoadImage(imagePath string) (image.Image, error) {
	inputFile, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, err
	}

	return img, nil
}

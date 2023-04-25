package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
	"github.com/kagami/go-avif"
)

type ImageEncodingService struct {
}

func NewImageEncodingService() *ImageEncodingService {
	return &ImageEncodingService{}
}

type Encoder interface {
	Encode(img image.Image) ([]byte, error)
}

type EncoderOptions struct {
	Quality int
}

func (s *ImageEncodingService) NewEncoder(format string, options EncoderOptions) (Encoder, error) {
	switch format {
	case "jpeg":
		return &JPEGEncoder{options}, nil
	case "png":
		// todo: handle png options
		return &PNGEncoder{}, nil
	case "webp":
		return &WebPEncoder{options}, nil
	case "avif":
		return &AVIFEncoder{options}, nil
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
}

type JPEGEncoder struct {
	options EncoderOptions
}

func (e *JPEGEncoder) Encode(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: e.options.Quality})
	return buf.Bytes(), err
}

type PNGEncoder struct{}

func (e *PNGEncoder) Encode(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	return buf.Bytes(), err
}

type WebPEncoder struct {
	options EncoderOptions
}

func (e *WebPEncoder) Encode(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := webp.Encode(buf, img, &webp.Options{Lossless: false, Quality: float32(e.options.Quality)})
	return buf.Bytes(), err
}

type AVIFEncoder struct {
	options EncoderOptions
}

func (e *AVIFEncoder) Encode(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	quality := int((1 - float64(e.options.Quality)/100) * 63)
	err := avif.Encode(buf, img, &avif.Options{Quality: quality})
	return buf.Bytes(), err
}
